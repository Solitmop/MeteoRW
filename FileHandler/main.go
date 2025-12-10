package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

// Meteostation структура метеостанции
type Meteostation struct {
	Index     uint    `json:"index"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`  // широта
	Longitude float64 `json:"longitude"` // долгота
	Altitude  int     `json:"altitude"`  // высота над уровнем моря
	Geohash   string  `json:"geohash"`
}

// CSVHandler обработчик CSV файлов
type CSVHandler struct {
	APIURL      string
	fieldNumber int
	client      *http.Client
}

// NewCSVHandler создает новый обработчик CSV
func NewCSVHandler(apiURL string) *CSVHandler {
	return &CSVHandler{
		APIURL: apiURL,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProcessFile обрабатывает CSV файл построчно и отправляет данные
func (h *CSVHandler) ProcessFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	//reader.FieldsPerRecord = -1 // Разрешаем разное количество полей

	successCount := 0
	failCount := 0
	lineNumber := 1
	_, err = reader.Read() // Пропускаем заголовок
	if err != nil {
		return fmt.Errorf("ошибка чтения заголовка: %v", err)
	}

	for {
		lineNumber++
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Ошибка чтения строки %d: %v", lineNumber, err)
			failCount++
			continue
		}

		station, err := h.parseRecord(record)
		if err != nil {
			log.Printf("Ошибка парсинга строки %d: %v", lineNumber, err)
			failCount++
			continue
		}

		if err := h.sendToAPI(station); err != nil {
			log.Printf("Ошибка отправки строки %d: %v", lineNumber, err)
			failCount++
		} else {
			successCount++
		}
	}

	log.Printf("Обработка завершена. Успешно: %d. Ошибок: %d. Всего файлов: %d", successCount, failCount, lineNumber)
	return nil
}

// parseRecord преобразует CSV запись в структуру Meteostation
func (h *CSVHandler) parseRecord(record []string) (Meteostation, error) {
	var station Meteostation

	// Парсим Index
	index, err := strconv.ParseUint(record[0], 10, 32)
	if err != nil {
		return station, fmt.Errorf("ошибка парсинга Index: %v", err)
	}
	station.Index = uint(index)

	// Name
	station.Name = record[1]

	// Latitude
	station.Latitude, err = strconv.ParseFloat(record[2], 64)
	if err != nil {
		return station, fmt.Errorf("ошибка парсинга Latitude: %v", err)
	}

	// Longitude
	station.Longitude, err = strconv.ParseFloat(record[3], 64)
	if err != nil {
		return station, fmt.Errorf("ошибка парсинга Longitude: %v", err)
	}

	// Altitude
	if record[4] != "" {
		station.Altitude, err = strconv.Atoi(record[4])
		if err != nil {
			return station, fmt.Errorf("ошибка парсинга Altitude: %v", err)
		}
	}

	// Geohash
	station.Geohash = ""

	return station, nil
}

// sendToAPI отправляет одну метеостанцию на API
func (h *CSVHandler) sendToAPI(station Meteostation) error {
	jsonData, err := json.Marshal(station)
	if err != nil {
		return fmt.Errorf("ошибка обработки JSON: %v", err)
	}

	resp, err := h.client.Post(h.APIURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("ошибка отправки запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API вернул ошибку %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

func main() {
	/*if len(os.Args) < 3 {
		log.Fatal("Использование: ./program <путь_к_csv_файлу> <api_url>")
	}

	filePath := os.Args[1]
	apiURL := os.Args[2]
	
	log.Printf("Начинаем обработку файла: %s", filePath)
	log.Printf("Отправка данных на API: %s", apiURL)
	*/

	filePath := "meteost.csv"
	apiURL := "http://localhost:8081/api/meteostations"
	handler := NewCSVHandler(apiURL)

	if err := handler.ProcessFile(filePath); err != nil {
		log.Fatalf("Ошибка обработки файла: %v", err)
	}

	log.Println("Обработка завершена")
}
