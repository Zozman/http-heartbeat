package main

import (
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Setup logger
	prodConfig := zap.NewProductionConfig()
	prodConfig.Encoding = "console"
	prodConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	prodConfig.EncoderConfig.EncodeDuration = zapcore.StringDurationEncoder
	// Enable debug logs if VERBOSE is enabled
	_, isVerboseDefined := os.LookupEnv("VERBOSE")
	if isVerboseDefined {
		prodConfig.Level.SetLevel(zap.DebugLevel)
	}
	logger, _ := prodConfig.Build()
	defer logger.Sync()
	sugar := logger.Sugar()

	// Check HEARTBEAT_URL
	heartbeatUrl, heartbeatUrlDefined := os.LookupEnv("HEARTBEAT_URL")
	sugar.Debugf("HEARTBEAT_URL Provided: %s", heartbeatUrl)

	// Check INTERVAL
	interval, intervalDefined := os.LookupEnv("INTERVAL")
	var intervalValue int
	var err error
	if intervalDefined {
		intervalValue, err = strconv.Atoi(interval)
		if err != nil {
			sugar.Warn("Error reading INTERVAL value")
			sugar.Warn(err)
			os.Exit(1)
		}
	} else {
		intervalValue = 30
	}
	sugar.Debugf("INTERVAL Set: %d Seconds", intervalValue)

	// Check TEST_URL
	testUrl, testUrlDefined := os.LookupEnv("TEST_URL")
	if testUrlDefined {
		sugar.Debugf("TEST_URL Provided: %s", testUrl)
	} else {
		sugar.Debug("No TEST_URL Provided")
	}

	// Fail if required environmental variable is missing
	if !heartbeatUrlDefined {
		sugar.Warn("Missing HEARTBEAT_URL.  Please provide a HEARTBEAT_URL as an environmental variable.")
		os.Exit(1)
	}

	// Perform check along interval
	for range time.Tick(time.Duration(intervalValue) * time.Second) {
		sugar.Debug("Starting heartbeat cycle...")
		// If there is a TEST_URL, check it first before performing the heartbeat
		if testUrlDefined {
			sugar.Debugf("Checking status of TEST_URL %s...", testUrl)
			resp, err := http.Get(testUrl)
			if err != nil {
				sugar.Warn("Error occured when sending GET to TEST_URL; skipping this cycle.")
				sugar.Warnf("%v", err)
			} else {
				// Read and close the body or memory will leak
				// https://stackoverflow.com/questions/69623172/go-http-package-why-do-i-need-to-close-res-body-after-reading-from-it
				io.ReadAll(resp.Body)
				resp.Body.Close()
				sugar.Debugf("GET to TEST_URL %s returned status code of %d", testUrl, resp.StatusCode)

				if resp.StatusCode == http.StatusOK {
					sugar.Debug("TEST_URL returned OK status; proceeding with this cycle...")
					sendHeartbeat(heartbeatUrl, sugar)
				} else {
					sugar.Debug("TEST_URL did not return OK status; skipping this cycle.")
				}
			}
			// Else just perform the heartbeat
		} else {
			sugar.Debug("No TEST_URL provided, proceeding with sending heartbeat")
			sendHeartbeat(heartbeatUrl, sugar)
		}
	}
}

func sendHeartbeat(targetUrl string, sugar *zap.SugaredLogger) {
	sugar.Debugf("Sending GET request to HEARTBEAT_URL %s...", targetUrl)
	resp, err := http.Get(targetUrl)
	if err != nil {
		sugar.Warn("Error occured when sending GET to HEARTBEAT_URL; skipping this cycle.")
		sugar.Warn(err)
	} else {
		// Read and close the body or memory will leak
		// https://stackoverflow.com/questions/69623172/go-http-package-why-do-i-need-to-close-res-body-after-reading-from-it
		io.ReadAll(resp.Body)
		resp.Body.Close()
		sugar.Debugf("GET to HEARTBEAT_URL %s returned status code of %d", targetUrl, resp.StatusCode)

		if resp.StatusCode == http.StatusOK {
			sugar.Info("Heartbeat Sent!")
		} else {
			sugar.Warnf("Heartbear URL did not return OK status; status was %d", resp.StatusCode)
		}
	}
}
