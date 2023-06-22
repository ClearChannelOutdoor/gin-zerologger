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