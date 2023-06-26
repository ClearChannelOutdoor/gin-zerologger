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


### Options

But wait, there's more! The `GinZeroLogger()` function accepts a variadic list of options to further tailor how request logging is handled. The following options are available:

#### Logging Options

Logging options can be provided to the GinZeroLogger at construction...

##### PathExclusion

To exclude certain paths from the logging middleware, use the `PathExclusion` option. The `PathExclusion` option accepts a variadic number of strings representing the paths to exclude from logging. For example, to exclude the `/routeToExclude` and `/anotherRouteToExclude` endpoints from logging:

```go
r.Use(gzl.GinZeroLogger(
  gzl.PathExclusion("/routeToExclude", "/anotherRouteToExclude")
))
```

##### LogLevelEqualToOrGreaterThan200

To override the log level for requests that return a 200 status code (the default is `log.Info()`), use the `LogLevelEqualToOrGreaterThan200` option. The `LogLevelEqualToOrGreaterThan200` option accepts either a string representing the log level to use, or a `*zerolog.Event` for requests that return a 200 status code. For example, to set the log level to `log.Debug()` for requests that return a 200 status code:

```go
r.Use(gzl.GinZeroLogger(
  gzl.LogLevelEqualToOrGreaterThan200(log.Debug()),
))
```

##### LogLevelEqualToOrGreaterThan400

To override the log level for requests that return a 400 status code (the default is `log.Warn()`), use the `LogLevelEqualToOrGreaterThan400` option. The `LogLevelEqualToOrGreaterThan400` option accepts a either string representing the log level to use, or a `*zerolog.Event` for requests that return a 400 status code. For example, to set the log level to `log.Info()` for requests that return a 400 status code:

```go
r.Use(gzl.GinZeroLogger(
  gzl.LogLevelEqualToOrGreaterThan400(log.Info()),
))
```

##### LogLevelEqualToOrGreaterThan500

To override the log level for requests that return a 500 status code (the default is `log.Error()`), use the `LogLevelEqualToOrGreaterThan500` option. The `LogLevelEqualToOrGreaterThan500` option accepts either a string representing the log level to use, or a `*zerolog.Event` for requests that return a 500 status code. For example, to set the log level to `log.Warn()` for requests that return a 500 status code:

```go
r.Use(gzl.GinZeroLogger(
  gzl.LogLevelEqualToOrGreaterThan500(log.Warn()),
))
```

##### Providing multiple options

Multiple options can be provided to the `GinZeroLogger()` function. For example, to exclude the `/health` endpoint from logging and set the log level to `log.Debug()` for requests that return a 200 status code:

```go
r.Use(gzl.GinZeroLogger(
  gzl.PathExclusion("/health", "/status"),
  gzl.LogLevelEqualToOrGreaterThan200(log.Debug()),
))
```

