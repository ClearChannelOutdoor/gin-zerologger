package ginzerologger

type StatusLevel int

const (
	Above200 StatusLevel = 200
	Above300 StatusLevel = 300
	Above400 StatusLevel = 400
	Above500 StatusLevel = 500
)

var LogBodyErrorStatusLevel = struct {
	Above200 StatusLevel
	Above300 StatusLevel
	Above400 StatusLevel
	Above500 StatusLevel
}{
	Above200: Above200,
	Above300: Above300,
	Above400: Above400,
	Above500: Above500,
}

type LoggingOption struct {
	Key   string
	Value any
}

func newLoggingOption(key string, value any) LoggingOption {
	return LoggingOption{
		Key:   key,
		Value: value,
	}
}

func LogRequestBodyOption(lvl StatusLevel) LoggingOption {
	return newLoggingOption("logRequestBody", lvl)
}

func LogPathExclusion(path string) LoggingOption {
	return newLoggingOption("exclude", path)
}

type optionsSearch struct {
	opts []LoggingOption
}

func newOptionsSearch(opts ...LoggingOption) *optionsSearch {
	return &optionsSearch{
		opts: opts,
	}
}

func (o *optionsSearch) Find(key string) (LoggingOption, bool) {
	for _, opt := range o.opts {
		if opt.Key == key {
			return opt, true
		}
	}

	return LoggingOption{}, false
}
