package handlers

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/deside01/is_available/internal/utils"
	gofpdf "github.com/jung-kurt/gofpdf"
)

func GetData(w http.ResponseWriter, r *http.Request) {
	file, err := os.Open("data.json")
	if err != nil {
		log.Printf("open file: %v", err)
		return
	}
	defer file.Close()

	data := make(map[string]map[string]string)
	json.NewDecoder(file).Decode(&data)

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	for _, dataMap := range data {
		for link, status := range dataMap {
			if link == "links_num" {
				continue
			}

			if pdf.GetY() > 255 {
				pdf.AddPage()
			}

			pdf.SetFont("Times", "B", 16)
			pdf.Cell(45, 10, fmt.Sprintf("%v:", link))
			pdf.SetFont("Times", "", 16)
			pdf.SetX(pdf.GetX() + 40)
			pdf.Cell(45, 10, status)
			// log.Printf("%v: %v. %v", link, status, listNum)

			pdf.SetY(pdf.GetY() + 10)

			// footer
			pdf.SetFooterFunc(func() {
				pdf.SetY(-15)
				pdf.SetFont("Times", "", 10)
				pdf.CellFormat(0, 10, strconv.Itoa(pdf.PageNo()), "", 0, "C", false, 0, "")
			})
		}
	}

	filename := fmt.Sprintf("%v.pdf", generateName())
	err = pdf.OutputFileAndClose(filename)
	if err != nil {
		panic(err)
	}

	utils.ResPDF(w, 200, filename)

	err = os.Remove(filename)
	if err != nil {
		log.Printf("os.Remove: %v", err)
	}
}

func generateName() string {
	n := 5
	b := make([]byte, n)
	rand.Read(b)

	s := fmt.Sprintf("%X", b)
	return s
}
