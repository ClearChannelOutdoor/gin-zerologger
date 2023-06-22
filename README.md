# Gin Logger

## Installation

```bash
go get github.com/clearchanneloutdoor/gin-logger
```

## Usage

```go
package main

import (
  "github.com/gin-gonic/gin"
  gzl "github.com/clearchanneloutdoor/gin-logger"
)

func main() {
  r := gin.New()

  // attach logger and recovery middleware
  r.Use(gzl.GinZeroLogger())

  // attach routes
  r.GET("/", func(c *gin.Context) {
    c.String(200, "Hello, World!")
  })

  // start the server
  r.Run(":8080")
}
```

Voila! Requests issued to the gin server will be logged via rs/zerolog!


### Options

But wait, there's more! The `GinZeroLogger()` function accepts a variadic list of options to further tailor how request logging is handled. The following options are available:

#### NewLoggingOption

The `NewLoggingOption` method can be used to create new logging options to provide to the middleware.

##### exclude

To exclude certain paths from the logging middleware, use the `exclude` option. The `exclude` option accepts a slice of strings representing the paths to exclude from logging. For example, to exclude the `/health` endpoint from logging:

```go
r.Use(gzl.GinZeroLogger(
  NewLoggingOption("exclude", []string{"/health"})
))
```

##### default200LogLevel

To override the log level for requests that return a 200 status code (the default is `log.Info()`), use the `default200LogLevel` option. The `default200LogLevel` option accepts a string representing the log level to use for requests that return a 200 status code. For example, to set the log level to `log.Debug()` for requests that return a 200 status code:

```go
r.Use(gzl.GinZeroLogger(
  NewLoggingOption("default200LogLevel", log.Debug())
))
```

##### default400LogLevel

To override the log level for requests that return a 400 status code (the default is `log.Warn()`), use the `default400LogLevel` option. The `default400LogLevel` option accepts a string representing the log level to use for requests that return a 400 status code. For example, to set the log level to `log.Info()` for requests that return a 400 status code:

```go
r.Use(gzl.GinZeroLogger(
  NewLoggingOption("default400LogLevel", log.Info())
))
```

##### default500LogLevel

To override the log level for requests that return a 500 status code (the default is `log.Error()`), use the `default500LogLevel` option. The `default500LogLevel` option accepts a string representing the log level to use for requests that return a 500 status code. For example, to set the log level to `log.Warn()` for requests that return a 500 status code:

```go
r.Use(gzl.GinZeroLogger(
  NewLoggingOption("default500LogLevel", log.Warn())
))
```

##### Providing multiple options

Multiple options can be provided to the `GinZeroLogger()` function. For example, to exclude the `/health` endpoint from logging and set the log level to `log.Debug()` for requests that return a 200 status code:

```go
r.Use(gzl.GinZeroLogger(
  NewLoggingOption("exclude", []string{"/health"}),
  NewLoggingOption("default200LogLevel", log.Debug())
))
```

