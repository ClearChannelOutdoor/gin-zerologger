package ginzerologger

type LoggingOption struct {
	Key   string
	Value any
}

func NewLoggingOption(key string, value any) LoggingOption {
	return LoggingOption{
		Key:   key,
		Value: value,
	}
}

type optionsSearch struct {
	opts []LoggingOption
}

func NewOptionsSearch(opts ...LoggingOption) *optionsSearch {
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
