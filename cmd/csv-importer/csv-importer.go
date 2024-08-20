package main

import (
	"encoding/csv"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/wolveix/gravityforms-go"
)

func main() {
	apiKey := flag.String("API_KEY", "", "Your Gravity Forms key")
	apiSecret := flag.String("API_SECRET", "", "Your Gravity Forms secret")
	apiURL := flag.String("API_URL", "https://your.domain.here/wp-json/gf/v2", "Your WordPress installation's URL")
	csvFilePath := flag.String("CSV_FILE_PATH", "", "Your input CSV (headers should be field IDs")
	debug := flag.Bool("DEBUG", true, "Enable debug logging for the API")
	formID := flag.Int("FORM_ID", 0, "Your Gravity Form ID")
	timeout := flag.Duration("TIMEOUT", 300*time.Second, "Specify the API timeout")
	flag.Parse()

	if *apiKey == "" {
		log.Fatal("Missing API_KEY")
	}

	if *apiSecret == "" {
		log.Fatal("Missing API_SECRET")
	}

	if *apiURL == "" {
		log.Fatal("Missing API_URL")
	}

	file, err := os.Open(*csvFilePath)
	if err != nil {
		log.Fatal("Failed to open " + *csvFilePath)
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		log.Fatal("Error reading CSV: ", err)
	}

	var columns []string

	entries := make([]*gravityforms.Entry, len(records)-1)

	for rowIndex, row := range records {
		if rowIndex == 0 {
			columns = row
			continue
		}

		rawEntry := make(map[string]string)

		for columnIndex, value := range row {
			if columnIndex > len(columns)-1 {
				continue
			}

			rawEntry[columns[columnIndex]] = strings.TrimPrefix(value, "\ufeff")
		}

		entryJSON, err := json.Marshal(rawEntry)
		if err != nil {
			fmt.Println("Failed to marshal entry for row:", rowIndex)
			continue
		}

		var entry gravityforms.Entry

		if err = json.Unmarshal(entryJSON, &entry); err != nil {
			fmt.Println("Failed to unmarshal entry for row:", rowIndex)
			continue
		}

		entries[rowIndex-1] = &entry
	}

	service := gravityforms.New(*apiURL, *apiKey, *apiSecret, *timeout, *debug)

	for index, entry := range entries {
		if err = service.CreateEntry(*formID, entry); err != nil {
			fmt.Println("Failed to create entry for row:", index-1, "error:", err.Error())
			continue
		}
	}
}
