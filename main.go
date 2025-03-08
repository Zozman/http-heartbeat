package main

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	// Setup logger
	logger := buildLogger()

	// Check HEARTBEAT_URL
	heartbeatUrl, heartbeatUrlDefined := os.LookupEnv("HEARTBEAT_URL")
	logger.Debug("HEARTBEAT_URL Provided", zap.String("url", heartbeatUrl))

	// Check INTERVAL
	interval, intervalDefined := os.LookupEnv("INTERVAL")
	var intervalValue int
	var err error
	if intervalDefined {
		intervalValue, err = strconv.Atoi(interval)
		if err != nil {
			logger.Warn("Error reading INTERVAL value", zap.Error(err))
			os.Exit(1)
		}
	} else {
		intervalValue = 30
	}
	logger.Debug("INTERVAL Set", zap.Int("interval", intervalValue))

	// Check TEST_URL
	testUrl, testUrlDefined := os.LookupEnv("TEST_URL")
	if testUrlDefined {
		logger.Debug("TEST_URL Provided", zap.String("url", testUrl))
	} else {
		logger.Debug("No TEST_URL Provided")
	}

	// Fail if required environmental variable is missing
	if !heartbeatUrlDefined {
		logger.Warn("Missing HEARTBEAT_URL. Please provide a HEARTBEAT_URL as an environmental variable.")
		os.Exit(1)
	}

	// Perform check along interval
	ticker := time.NewTicker(time.Duration(intervalValue) * time.Second)
	defer ticker.Stop()

	logger.Info("Heartbeat Service Started!")

	for range ticker.C {
		logger.Debug("Starting heartbeat cycle...")
		// If there is a TEST_URL, check it first before performing the heartbeat
		if testUrlDefined {
			logger.Debug("Checking status of TEST_URL", zap.String("url", testUrl))
			resp, err := http.Get(testUrl)
			if err != nil {
				logger.Warn("Error occurred when sending GET to TEST_URL; skipping this cycle.", zap.Error(err))
				continue
			} else {
				resp.Body.Close()
				logger.Debug("GET to TEST_URL returned status code", zap.String("url", testUrl), zap.Int("status", resp.StatusCode))

				if resp.StatusCode == http.StatusOK {
					logger.Debug("TEST_URL returned OK status; proceeding with this cycle...")
					sendHeartbeat(heartbeatUrl, logger)
				} else {
					logger.Debug("TEST_URL did not return OK status; skipping this cycle.")
				}
			}
			// Else just perform the heartbeat
		} else {
			logger.Debug("No TEST_URL provided, proceeding with sending heartbeat")
			sendHeartbeat(heartbeatUrl, logger)
		}
	}
}

func sendHeartbeat(targetUrl string, logger *zap.Logger) {
	logger.Debug("Sending GET request to HEARTBEAT_URL", zap.String("url", targetUrl))
	resp, err := http.Get(targetUrl)
	if err != nil {
		logger.Warn("Error occurred when sending GET to HEARTBEAT_URL; skipping this cycle.", zap.Error(err))
		return
	} else {
		resp.Body.Close()
		logger.Debug("GET to HEARTBEAT_URL returned status code", zap.String("url", targetUrl), zap.Int("status", resp.StatusCode))

		if resp.StatusCode == http.StatusOK {
			logger.Info("Heartbeat Sent!")
		} else {
			logger.Warn("Heartbeat URL did not return OK status", zap.Int("status", resp.StatusCode))
		}
	}
}

func buildLogger() *zap.Logger {
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
	return logger
}
