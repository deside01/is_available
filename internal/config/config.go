package config

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/joho/godotenv"
)

type ConfigData struct {
	Mu         sync.Mutex
	Headers    map[string][]string
	QueueLimit int
	Address    string
}

var Data = ConfigData{
	QueueLimit: 3,
	Headers: map[string][]string{
		"Accept":          {"*/*"},
		"Accept-Encoding": {"gzip, deflate, br zstd"},
		"Connection":      {"Close"},
		"User-Agent":      {"Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:144.0) Gecko/20100101 Firefox/144.0"},
	},
	Address: getAddress(),
}

func getAddress() string {
	godotenv.Load()

	host := strings.TrimSpace(os.Getenv("HOST"))
	if host == "" {
		host = "127.0.0.1"
	}

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		log.Fatal("empty PORT env")
	}

	return fmt.Sprintf("%v:%v", host, port)
}
