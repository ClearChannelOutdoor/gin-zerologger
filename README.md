# Gin Logger

## Installation

```bash
go get github.com/clearchanneloutdoor/gin-zerologger
```

## Usage

```go
package main

import (
  "github.com/gin-gonic/gin"
  gzl "github.com/clearchanneloutdoor/gin-zerologger"
)

func main() {
  r := gin.New()

  // attach logger and recovery middleware
  r.Use(gzl.GinZeroLogger())

  // attach routes
  r.GET("/", func(ctx *gin.Context) {
    ctx.String(200, "Hello, World!")
  })

  // start the server
  r.Run(":8080")
}
```

Voila! Requests issued to the gin server will be logged via rs/zerolog!

## Example

To see an example with a few of the logging options in use, check out the [`example/main.go`](example/main.go) file in this repository.

### Logging Options

But wait, there's more! The `GinZeroLogger()` function accepts a variadic list of options to further tailor how request logging is handled. Logging options can be provided to the GinZeroLogger at construction...

##### IncludeRequestBody

To include the request body in the log output, use the `IncludeRequestBody` option. The `IncludeRequestBody` option accepts an `HTTPStatus` value. For example, to include the request body in the log output for all HTTP responses with a status of 500 or greater:

```go
r.Use(gzl.GinZeroLogger(
  gzl.IncludeRequestBody(gzl.HTTPStatusCodes.EqualToOrGreaterThan500)
))
```

The following status options are supported:

* `HTTPStatusCodes.EqualToOrGreaterThan500`
* `HTTPStatusCodes.EqualToOrGreaterThan400`
* `HTTPStatusCodes.EqualToOrGreaterThan300`
* `HTTPStatusCodes.EqualToOrGreaterThan200`

##### IncludeContextValues

To include values from the gin context in the log output, use the `IncludeContextValues` option. The `IncludeContextValues` option accepts a variadic number of strings representing the keys of the values to include in the log output. For example, to include the values of the `ClientID` and `SessionID` gin context key/value pairs in the log output:

```go
r.Use(gzl.GinZeroLogger(
  gzl.IncludeContextValues("ClientID", "SessionID")
))
```

##### PathExclusion

To exclude certain paths from the logging middleware, use the `PathExclusion` option. The `PathExclusion` option accepts a variadic number of strings representing the paths to exclude from logging. For example, to exclude the `/routeToExclude` and `/anotherRouteToExclude` endpoints from logging:

```go
r.Use(gzl.GinZeroLogger(
  gzl.PathExclusion("/routeToExclude", "/anotherRouteToExclude")
))
```

##### LogLevel200

To override the log level for requests that return a 200 status code (the default is `log.Info()`), use the `LogLevel200` option. The `LogLevel200` option accepts either a string representing the log level to use, or a `*zerolog.Event` for requests that return a 200 status code. For example, to set the log level to `log.Debug()` for requests that return a 200 status code:

```go
r.Use(gzl.GinZeroLogger(
  gzl.LogLevel200(log.Debug()),
))
```

##### LogLevel400

To override the log level for requests that return a 400 status code (the default is `log.Warn()`), use the `LogLevel400` option. The `LogLevel400` option accepts a either string representing the log level to use, or a `*zerolog.Event` for requests that return a 400 status code. For example, to set the log level to `log.Info()` for requests that return a 400 status code:

```go
r.Use(gzl.GinZeroLogger(
  gzl.LogLevel400(log.Info()),
))
```

##### LogLevel500

To override the log level for requests that return a 500 status code (the default is `log.Error()`), use the `LogLevel500` option. The `LogLevel500` option accepts either a string representing the log level to use, or a `*zerolog.Event` for requests that return a 500 status code. For example, to set the log level to `log.Warn()` for requests that return a 500 status code:

```go
r.Use(gzl.GinZeroLogger(
  gzl.LogLevel500(log.Warn()),
))
```

##### Providing multiple options

Multiple options can be provided to the `GinZeroLogger()` function. For example, to exclude the `/health` endpoint from logging and set the log level to `log.Debug()` for requests that result in a 200 status code:

```go
r.Use(gzl.GinZeroLogger(
  gzl.PathExclusion("/health", "/status"),
  gzl.LogLevel200(log.Debug()),
))
```

