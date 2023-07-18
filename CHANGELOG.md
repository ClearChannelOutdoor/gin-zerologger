# Change Log
All notable changes in gin-zerologger will be documented in this file.
 
The format is based on [Keep a Changelog](http://keepachangelog.com/)
and this project adheres to [Semantic Versioning](http://semver.org/).

## [1.3.0] - 2023-07-18

- When present in the context, including `X-Correlation-ID` and `X-Request-ID` in the log output

## [1.2.0] - 2023-06-27

- Updated documentation to include logging options for request body
- Modified the log level options to shorten the option name

## [1.1.0] - 2023-06-22

- Added logging options to `GinZeroLogger()` function to further tailor how request logging is handled
- Added documentation to the README.md
- Added `example/main.go` to demonstrate usage of the module

## [1.0.0] - 2023-06-21

- Initial implementation of gin-zerologger