package utils

import (
	"net/http"
	"os"
	"strconv"
)

func ResPDF(w http.ResponseWriter, statusCode int, filename string) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/pdf")
	w.Header().Set("Content-Disposition", "attachment; filename=\""+filename+"\"")
	w.Header().Set("Content-Length", strconv.Itoa(len(fileBytes)))

	w.WriteHeader(statusCode)

	w.Write(fileBytes)
}
