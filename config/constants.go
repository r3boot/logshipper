package config

/*
 * Various constants used to configure the input system
 */
const T_SYSLOG string = "syslog" // Syslog logging with RFC3339 timestamps
const T_CLF string = "clf"       // HTTP Common Log Format
const T_JSON string = "json"     // JSON log format

/*
 * Various constants used to configure the output system
 */
const T_STDOUT int = 0 // Write to stdout
const T_REDIS int = 1  // Write to redis

/*
 * Various constants used to control the input/output threads
 */
const CMD_CLEANUP int = 0 // Stop whatever you're doing and cleanup
