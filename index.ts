import winston from 'winston';

// Load environmental variables
// URL to send heartbeat to
const heartbeatUrl = process.env.HEARTBEAT_URL;
// Optional test url to test for OK status before sending heartbeat
const testUrl = process.env.TEST_URL;
// Default interval is every 30 seconds
const interval = process.env.INTERVAL
    ? parseInt(process.env.INTERVAL)
    : 30;
// Uses winston log levels
const logLevel = process.env.VERBOSE
    ? 'verbose'
    : 'info';

// Create logger
const logger = winston.createLogger({
    level: logLevel,
    format: winston.format.cli(),
    transports: [new winston.transports.Console()],
});

logger.verbose(`HEARTBEAT_URL Provided: ${heartbeatUrl}`);
logger.verbose(`Interval Set: ${interval} Seconds`);
logger.verbose(testUrl
    ? `TEST_URL Provided: ${testUrl}`
    : 'No TEST_URL Provided'
);

// Fail if required environmental variable is missing
if (!heartbeatUrl) {
    logger.error('Missing HEARTBEAT_URL.  Please provide a HEARTBEAT_URL as an environmental variable.');
    process.exit(1);
}

// Perform heartbeat cycle
const cycleHeartbeat = async () => {
    logger.verbose('Starting Heartbeat Cycle...');
    if (testUrl) {
        logger.verbose(`Checking Status of Test URL ${testUrl}...`);
        try {
            const req = await fetch(testUrl);
            logger.verbose(`GET to Test URL ${testUrl} returned status code of ${req.status}`);
            if (req.ok) {
                logger.verbose('Test URL returned OK status; proceeding with this cycle...');
                sendHeartbeat();
            } else {
                logger.verbose('Test URL did not return OK status; skipping this cycle.');
            }
        } catch(e) {
            logger.error('Error occured when sending GET to Test URL; skipping this cycle.');
            logger.error(e);
        }
    } else {
        logger.verbose('No Test URL Provided, proceeding with sendinng heartbeat')
        sendHeartbeat();
    }
};

// Send GET to heartbeatUrl
const sendHeartbeat = async () => {
    logger.verbose(`Sending GET request to Heartbeat URL ${heartbeatUrl}...`);
    try {
        const req = await fetch(heartbeatUrl);
        logger.verbose(`GET to Heartbeat URL ${heartbeatUrl} returned status code of ${req.status}`);
        if (req.ok) {
            logger.info('Heartbeat Sent!');
        } else {
            logger.error(`Heartbear URL did not return OK status; status was ${req.status}`);
        }
    } catch(e) {
        logger.error('Error occured when sending GET to Heartbeat URL; skipping this cycle.');
            logger.error(e);
    }
}

setInterval(cycleHeartbeat, interval * 1000);

logger.info('Heartbeat Process Started!');