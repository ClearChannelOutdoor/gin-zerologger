package ginzerologger

import "github.com/rs/zerolog"

type HTTPStatus int

var HTTPStatusCodes = struct {
	EqualToOrGreaterThan200 HTTPStatus
	EqualToOrGreaterThan300 HTTPStatus
	EqualToOrGreaterThan400 HTTPStatus
	EqualToOrGreaterThan500 HTTPStatus
}{
	EqualToOrGreaterThan200: 2,
	EqualToOrGreaterThan300: 3,
	EqualToOrGreaterThan400: 4,
	EqualToOrGreaterThan500: 5,
}

type logLevel interface {
	*zerolog.Event | string
}

type loggingOption struct {
	Key   string
	Value any
}

func newLoggingOption(key string, value any) *loggingOption {
	return &loggingOption{
		Key:   key,
		Value: value,
	}
}

func IncludeRequestBody(sts HTTPStatus) *loggingOption {
	return newLoggingOption("includeRequestBody", sts)
}

func LogLevel200[T logLevel](val T) *loggingOption {
	return &loggingOption{
		Key:   "default200",
		Value: val,
	}
}

func LogLevel300[T logLevel](val T) *loggingOption {
	return &loggingOption{
		Key:   "default300",
		Value: val,
	}
}

func LogLevel400[T logLevel](val T) *loggingOption {
	return &loggingOption{
		Key:   "default400",
		Value: val,
	}
}

func LogLevel500[T logLevel](val T) *loggingOption {
	return &loggingOption{
		Key:   "default500",
		Value: val,
	}
}

func PathExclusion(paths ...string) *loggingOption {
	if len(paths) > 0 {
		return newLoggingOption("excludes", paths)
	}

	return nil
}

type optionsSearch struct {
	opts []*loggingOption
}

func newOptionsSearch(opts ...*loggingOption) *optionsSearch {
	return &optionsSearch{
		opts: opts,
	}
}

func (o *optionsSearch) Find(key string) (*loggingOption, bool) {
	for _, opt := range o.opts {
		if opt.Key == key {
			return opt, true
		}
	}

	return nil, false
}
