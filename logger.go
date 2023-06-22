package ginzerologger

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggingDetails interface {
	Details() map[string]any
}

func augmentLogEvent(err LoggingDetails, le *zerolog.Event) {
	for k, v := range err.Details() {
		le = le.Interface(k, v)
	}
}

func GinZeroLogger(opts ...LoggingOption) gin.HandlerFunc {
	// create a search for the options
	search := NewOptionsSearch(opts...)

	return func(ctx *gin.Context) {
		// capture request duration
		t := time.Now()

		// process request
		ctx.Next()

		// do not log for any excluded paths (i.e. /v1/status)
		if excludes, ok := search.Find("exclude"); ok {
			for _, path := range excludes.Value.([]string) {
				if ctx.Request.URL.Path == path {
					return
				}
			}
		}

		// create a logging event and augment it
		var le *zerolog.Event

		switch sts := ctx.Writer.Status(); {
		case sts >= 500:
			if dle, ok := search.Find("default500LogLevel"); ok {
				le = dle.Value.(*zerolog.Event)
				break
			}

			le = log.Error()
		case sts >= 400:
			if dle, ok := search.Find("default400LogLevel"); ok {
				le = dle.Value.(*zerolog.Event)
				break
			}

			le = log.Warn()
		default:
			if dle, ok := search.Find("default200LogLevel"); ok {
				le = dle.Value.(*zerolog.Event)
				break
			}

			le = log.Info()
		}

		// add request detail to the error
		le = le.
			Dur("duration", time.Since(t)).
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Int("status", ctx.Writer.Status())

		// add query if there is one
		if ctx.Request.URL.RawQuery != "" {
			le = le.Str("query", ctx.Request.URL.RawQuery)
		}

		// request has a single error
		if len(ctx.Errors) == 1 {
			err := ctx.Errors[0].Err
			le = le.Err(err)

			// check to see if the error has any additional details
			if dErr, ok := err.(LoggingDetails); ok {
				augmentLogEvent(dErr, le)
			}

			le.Send()
			return
		}

		// more than 1 error
		if len(ctx.Errors) > 1 {
			err := ctx.Errors.Last().Err
			le = le.Err(err)

			// check to see if the error has any additional details
			if dErr, ok := err.(LoggingDetails); ok {
				augmentLogEvent(dErr, le)
			}

			le.Msg(ctx.Errors.String())
			return
		}

		// send the details
		le.Send()
	}
}
