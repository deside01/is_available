package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/deside01/is_available/internal/utils"
)

type ReqBody struct {
	Links []string `json:"links"`
}

var mu = sync.Mutex{}

func Check(w http.ResponseWriter, r *http.Request) {
	var body ReqBody
	json.NewDecoder(r.Body).Decode(&body)

	wg := &sync.WaitGroup{}
	queue := make(chan struct{}, 3)
	dataMap := make(map[string]map[string]any)
	newData := make(map[string]any)

	for _, link := range body.Links {
		wg.Add(1)
		queue <- struct{}{}

		go func(link string) {
			defer func() {
				<-queue
				wg.Done()
			}()

			stat, _ := getStatus(link, newData)
			log.Println(stat + link)
		}(link)
	}
	wg.Wait()

	file, err := os.OpenFile("data.json", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("openFile: %v", err)
		return
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&dataMap)

	oldDataLength := len(dataMap) + 1
	nextKey := strconv.Itoa(oldDataLength)

	dataMap[nextKey] = newData

	if _, err := file.Seek(0, 0); err != nil {
		panic(err)
	}
	if err := file.Truncate(0); err != nil {
		panic(err)
	}

	newData["links_num"] = oldDataLength

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	encoder.Encode(dataMap)

	utils.ResJSON(w, 201, newData)
}

func getStatus(link string, newData map[string]any) (string, error) {
	safeLink := link
	if !hasProtocol(link) {
		safeLink = fmt.Sprintf("https://%v", link)
	}

	resp, err := http.Get(safeLink)
	if err != nil {
		log.Printf("Не удалось сделать запрос: %v", err)
		return "not available", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Println(resp.Status)
		return "not available", nil
	}

	mu.Lock()
	defer mu.Unlock()

	newData[link] = "available"

	return "available", nil
}

func hasProtocol(link string) bool {
	return strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")
}
