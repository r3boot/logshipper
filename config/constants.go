package config

/*
 * Various constants used to configure the input system
 */
const T_SYSLOG int = 0 // Syslog logging with RFC3339 timestamps
const T_CLF int = 1    // HTTP Common Log Format

/*
 * Various constants used to configure the output system
 */
const T_STDOUT int = 0 // Write to stdout
const T_REDIS int = 1  // Write to redis

/*
 * Various constants used to control the input/output threads
 */
const CMD_CLEANUP int = 0 // Stop whatever you're doing and cleanup
