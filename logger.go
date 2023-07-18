package ginzerologger

import (
	"bytes"
	"io"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LoggingDetails interface {
	Details() map[string]any
}

func augmentLogEvent(err LoggingDetails, lctx *zerolog.Context) {
	lc := *lctx
	for k, v := range err.Details() {
		lc = lc.Interface(k, v)
	}
	*lctx = lc
}

func logEventWithContext(sts int, search *optionsSearch, lctx zerolog.Context, msg ...string) {
	var lgr *zerolog.Event
	l := lctx.Logger()

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
				if val != nil {
					t := reflect.ValueOf(val).Elem().FieldByName("level")
					switch t.Int() {
					case -1:
						lgr = l.Trace()
					case 0:
						lgr = l.Debug()
					case 1:
						lgr = l.Info()
					case 2:
						lgr = l.Warn()
					case 3:
						lgr = l.Error()
					case 4:
						lgr = l.Fatal()
					case 5:
						lgr = l.Panic()
					default:
						lgr = l.Info()
					}
				}
			case string:
				switch val {
				case "trace":
					lgr = l.Trace()
				case "debug":
					lgr = l.Debug()
				case "info":
					lgr = l.Info()
				case "warn":
					lgr = l.Warn()
				case "error":
					lgr = l.Error()
				case "fatal":
					lgr = l.Fatal()
				case "panic":
					lgr = l.Panic()
				default:
					lgr = l.Info()
				}
			}
		}
	}

	// default 400s to warn
	if sts >= 400 {
		lgr = l.Warn()
	}

	// default 500s to warn
	if sts >= 500 {
		lgr = l.Error()
	}

	// default to info
	if lgr == nil {
		lgr = l.Info()
	}

	// send the log event with a message if one is provided
	if len(msg) > 0 && msg[0] != "" {
		lgr.Msg(msg[0])
		return
	}

	// send the log event
	lgr.Send()
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

		var lctx zerolog.Context

		// add request detail to the error
		lctx = log.With().
			Dur("duration", time.Since(t)).
			Str("method", ctx.Request.Method).
			Str("path", ctx.Request.URL.Path).
			Int("status", ctx.Writer.Status())

		// add X-Correlation-ID and X-Request-ID if the exist
		for _, hdr := range []string{"X-Correlation-ID", "X-Request-ID"} {
			if rid := ctx.Request.Header.Get(hdr); rid != "" {
				lctx = lctx.Str(hdr, rid)
			}
		}

		// check to see if request body should be included in the log
		if opt, ok := search.Find("includeRequestBody"); ok && len(bdy) > 0 {
			if logSts, ok := opt.Value.(HTTPStatus); ok {
				if ctx.Writer.Status()/100 == int(logSts) {
					if ct := ctx.Request.Header.Get("content-type"); strings.Contains(ct, "application/json") {
						lctx = lctx.RawJSON("body", bdy)
					} else {
						lctx = lctx.Str("body", string(bdy))
					}
				}
			}
		}

		// add query if there is one
		if ctx.Request.URL.RawQuery != "" {
			lctx = lctx.Str("query", ctx.Request.URL.RawQuery)
		}

		// more than 1 error
		if len(ctx.Errors) > 1 {
			err := ctx.Errors.Last().Err
			lctx = lctx.Err(err)

			// check to see if the error has any additional details
			if dErr, ok := err.(LoggingDetails); ok {
				augmentLogEvent(dErr, &lctx)
			}

			logEventWithContext(ctx.Writer.Status(), search, lctx, ctx.Errors.String())
			return
		}

		// request has a single error
		if len(ctx.Errors) == 1 {
			err := ctx.Errors[0].Err
			lctx = lctx.Err(err)

			// check to see if the error has any additional details
			if dErr, ok := err.(LoggingDetails); ok {
				augmentLogEvent(dErr, &lctx)
			}
		}

		// send the details
		logEventWithContext(ctx.Writer.Status(), search, lctx)
	}
}
