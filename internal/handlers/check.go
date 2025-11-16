package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"maps"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/deside01/is_available/internal/config"
	"github.com/deside01/is_available/internal/utils"
)

func Check(w http.ResponseWriter, r *http.Request) {
	links, err := parseBody(r)
	if err != nil {
		utils.ResERR(w, 400, err.Error())
		return
	}

	wg := &sync.WaitGroup{}
	queue := make(chan struct{}, config.Data.QueueLimit)

	dataMap := make(map[string]map[string]any)
	newData := make(map[string]any)

	for _, link := range links {
		wg.Add(1)
		queue <- struct{}{}

		go func(link string) {
			defer func() {
				<-queue
				wg.Done()
			}()

			tempData := make(map[string]any)

			err := getStatus(link, tempData)
			if err != nil {
				log.Println(err)
				return
			}

			config.Data.Mu.Lock()
			maps.Copy(newData, tempData)
			config.Data.Mu.Unlock()

			log.Printf("%v: OK", link)
		}(link)
	}
	wg.Wait()

	config.Data.Mu.Lock()

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

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	encoder.Encode(dataMap)

	config.Data.Mu.Unlock()

	newData["links_num"] = oldDataLength

	utils.ResJSON(w, 201, newData)
}

func getStatus(link string, newData map[string]any) error {
	safeLink := link
	if !hasProtocol(link) {
		safeLink = fmt.Sprintf("https://%v", link)
	}

	req, err := http.NewRequest(http.MethodGet, safeLink, nil)
	if err != nil {
		return fmt.Errorf("не удалось создать запрос: %w", err)
	}
	req.Header = config.Data.Headers

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		newData[link] = "not available"
		return fmt.Errorf("не удалось сделать запрос: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		newData[link] = "not available"
		return nil
	}

	newData[link] = "available"

	return nil
}

func hasProtocol(link string) bool {
	return strings.HasPrefix(link, "http://") || strings.HasPrefix(link, "https://")
}
