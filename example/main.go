package main

import (
	"io"
	"net/http"

	gzl "github.com/clearchanneloutdoor/gin-zerologger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
)

type customError struct {
	msg string
}

func (ce *customError) Error() string {
	return ce.msg
}

func (ce *customError) Details() map[string]interface{} {
	return map[string]interface{}{
		"custom": "error",
		"test":   true,
	}
}

func main() {
	r := gin.New()

	// attach logger and recovery middleware
	r.Use(gzl.GinZeroLogger(
		gzl.IncludeRequestBody(gzl.HTTPStatusCodes.EqualToOrGreaterThan200),
		gzl.PathExclusion("/notlogged"),
		gzl.LogLevel200(log.Debug()),
	))

	// attach routes
	r.GET("/", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "hello world!",
		})
	})

	r.GET("/400", func(ctx *gin.Context) {
		ctx.JSON(http.StatusNotFound, map[string]string{
			"error": "something not found",
		})
	})

	r.GET("/500", func(ctx *gin.Context) {
		// demonstrate LoggingDetails interface usage
		ce := &customError{
			msg: "an internal server error",
		}
		ctx.Error(ce)

		ctx.JSON(http.StatusInternalServerError, map[string]string{
			"error": "internal server error",
		})
	})

	r.POST("/logbody", func(ctx *gin.Context) {
		bdy, _ := io.ReadAll(ctx.Request.Body)
		ctx.JSON(http.StatusAccepted, map[string]string{
			"body": string(bdy),
		})
	})

	r.GET("/notlogged", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "not logged",
		})
	})

	// start the server
	r.Run(":8080")
}
