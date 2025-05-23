# Change Log
All notable changes in gin-zerologger will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [1.3.4] - 2025-04-01

- Fixed an issue where 500-level responses were logged at `warn` level instead of `error` by default

## [1.3.3] - 2025-02-28

- Fixed an issue when global logging level conflicts with level-specific logging, log messages were produced unexpectedly:
  - for example, when global logging level is set to `info` while a 200-level logging is set to `debug`, the logger will suppress any 200-level log messages (thus honoring the `info` global logging level)
- Enhancement added to the example to support an optionally provided global logging level via an environment variable `LOGGING_LEVEL`

## [1.3.2] - 2023-12-06

- Fixed issue where log level override for 400s and 500s was not being applied correctly

## [1.3.1] - 2023-07-20

- When present in the gin context, adding additional configurable values to the log output

## [1.3.0] - 2023-07-18

- When present in the context, including `X-Correlation-ID` and `X-Request-ID` in the log output
- Fixed issue where multiple (and incomplete) log messages were being sent per single request

## [1.2.0] - 2023-06-27

- Updated documentation to include logging options for request body
- Modified the log level options to shorten the option name

## [1.1.0] - 2023-06-22

- Added logging options to `GinZeroLogger()` function to further tailor how request logging is handled
- Added documentation to the README.md
- Added `example/main.go` to demonstrate usage of the module

## [1.0.0] - 2023-06-21

- Initial implementation of gin-zerologger