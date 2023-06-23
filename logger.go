package ginzerologger

import (
	"bytes"
	"io"
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

func defaultLogLevelEvent(sts int, search *optionsSearch) *zerolog.Event {
	for lvl, key := range map[int]string{
		5: "default500LogLevel",
		4: "default400LogLevel",
		3: "default300LogLevel",
		2: "default200LogLevel",
	} {
		// convert a XXX error to X for comparison purposes
		if sts/100 != lvl {
			continue
		}

		// check to see if there is a specific logging level to use
		if dle, ok := search.Find(key); ok {
			switch val := dle.Value.(type) {
			case *zerolog.Event:
				return val
			case string:
				return getLogEventForString(val)
			}
		}
	}

	// default 500s to warn
	if sts >= 500 {
		return log.Error()
	}

	// default 400s to warn
	if sts >= 400 {
		return log.Warn()
	}

	// default to info
	return log.Info()
}

func getLogEventForString(level string) *zerolog.Event {
	switch level {
	case "debug":
		return log.Debug()
	case "info":
		return log.Info()
	case "warn":
		return log.Warn()
	case "error":
		return log.Error()
	case "fatal":
		return log.Fatal()
	case "panic":
		return log.Panic()
	default:
		return log.Info()
	}
}

func pathIsExcluded(path string, opt LoggingOption) bool {
	switch val := opt.Value.(type) {
	case []string:
		for _, p := range val {
			if path == p {
				return true
			}
		}
	case string:
		if path == val {
			return true
		}
	}

	return false
}

func GinZeroLogger(opts ...LoggingOption) gin.HandlerFunc {
	// create a search for the options
	search := newOptionsSearch(opts...)

	return func(ctx *gin.Context) {
		// capture request duration
		t := time.Now()

		var buf bytes.Buffer

		// check to see if we should log the request body
		if _, ok := search.Find("logRequestBody"); ok {
			// read the request body
			io.Copy(&buf, ctx.Request.Body)

			// restore the io.ReadCloser to its original state
			ctx.Request.Body = io.NopCloser(&buf)
		}

		// process request
		ctx.Next()

		// do not log for any excluded paths (i.e. /v1/status)
		if excludes, ok := search.Find("exclude"); ok {
			if pathIsExcluded(ctx.Request.URL.Path, excludes) {
				return
			}
		}

		// create a logging event and augment it
		le := defaultLogLevelEvent(ctx.Writer.Status(), search)

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
