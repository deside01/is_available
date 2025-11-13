package config

import "sync"

type ConfigData struct {
	Mu         sync.Mutex
	Headers    map[string][]string
	QueueLimit int
}

var Data = ConfigData{
	QueueLimit: 3,
	Headers: map[string][]string{
		"Accept":          {"*/*"},
		"Accept-Encoding": {"gzip, deflate, br zstd"},
		"Connection":      {"Close"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:144.0) Gecko/20100101 Firefox/144.0"},
	},
}
