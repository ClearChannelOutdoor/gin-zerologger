package ginzerologger

import (
	"bytes"
	"io"
	"strings"
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
		5: "default500",
		4: "default400",
		3: "default300",
		2: "default200",
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

func pathIsExcluded(path string, opt *loggingOption) bool {
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

func GinZeroLogger(opts ...*loggingOption) gin.HandlerFunc {
	// create a search for the options
	search := newOptionsSearch(opts...)

	return func(ctx *gin.Context) {
		// capture request duration
		t := time.Now()

		var bdy []byte

		// check to see if we should collect the request body
		if _, ok := search.Find("includeRequestBody"); ok {
			// read the request body
			bdy, _ = io.ReadAll(ctx.Request.Body)

			// restore the io.ReadCloser to its original state for downstream
			// processing...
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bdy))
		}

		// process request
		ctx.Next()

		// do not log for any excluded paths (i.e. /v1/status)
		if excludes, ok := search.Find("excludes"); ok {
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

		// add X-Correlation-ID and X-Request-ID if the exist
		for _, hdr := range []string{"X-Correlation-ID", "X-Request-ID"} {
			if rid := ctx.Request.Header.Get(hdr); rid != "" {
				le = le.Str(hdr, rid)
			}
		}

		// check to see if request body should be included in the log
		if opt, ok := search.Find("includeRequestBody"); ok && len(bdy) > 0 {
			if logSts, ok := opt.Value.(HTTPStatus); ok {
				if ctx.Writer.Status()/100 == int(logSts) {
					if ct := ctx.Request.Header.Get("content-type"); strings.Contains(ct, "application/json") {
						le.RawJSON("body", bdy)
					} else {
						le.Str("body", string(bdy))
					}
				}
			}
		}

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
