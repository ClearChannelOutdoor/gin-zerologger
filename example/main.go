package main

import (
	"io"
	"net/http"
	"os"

	gzl "github.com/clearchanneloutdoor/gin-zerologger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
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

	// if a global logging level was passed in via an environment variable
	// set the global logging level to that level
	if lvl, ok := os.LookupEnv("LOGGING_LEVEL"); ok {
		l, err := zerolog.ParseLevel(lvl)
		if err != nil {
			log.Fatal().Err(err).Msg("unable to parse log level")
		}
		zerolog.SetGlobalLevel(l)
		log.Info().Str("level", l.String()).Msg("global logging level set")
	}

	// attach logger and recovery middleware
	r.Use(gzl.GinZeroLogger(
		gzl.IncludeRequestBody(gzl.HTTPStatusCodes.EqualToOrGreaterThan200),
		gzl.PathExclusion("/notlogged"),
		gzl.LogLevel200(log.Debug()),
		gzl.LogLevel300(log.Info()),
		gzl.IncludeContextValues("clientID", "random"),
	))

	// attach routes
	r.GET("/", func(ctx *gin.Context) {
		log.Trace().Str("level", "trace").Msg("running / route...")
		log.Debug().Str("level", "debug").Msg("running / route...")
		log.Info().Str("level", "info").Msg("running / route...")
		ctx.JSON(http.StatusOK, map[string]string{
			"message": "hello world!",
		})
	})

	r.GET("/300", func(ctx *gin.Context) {
		ctx.JSON(http.StatusPermanentRedirect, map[string]string{
			"message": "redirecting",
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

	r.GET("/logcontext", func(ctx *gin.Context) {
		ctx.Set("clientID", "12345")
		ctx.Set("random", "random value")

		ctx.JSON(http.StatusOK, map[string]string{
			"message": "context included",
		})
	})

	// start the server
	r.Run(":8080")
}
