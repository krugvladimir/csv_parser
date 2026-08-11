// csv_parser.go — Go версия

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func detectDelimiter(filePath string) rune {
	file, _ := os.Open(filePath)
	defer file.Close()
	reader := csv.NewReader(file)
	firstLine, _ := reader.Read()
	if len(firstLine) == 0 {
		return ','
	}
	// Пробуем определить по количеству полей
	for _, delim := range []rune{',', ';', '\t', '|'} {
		r := csv.NewReader(strings.NewReader(strings.Join(firstLine, string(delim))))
		fields, _ := r.Read()
		if len(fields) > 1 {
			return delim
		}
	}
	return ','
}

func parseValue(value string) interface{} {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil
	}
	switch strings.ToLower(value) {
	case "true":
		return true
	case "false":
		return false
	case "null", "none":
		return nil
	}
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}
	return value
}

func parseNestedKeys(headers []string) map[string]interface{} {
	result := make(map[string]interface{})
	for _, key := range headers {
		parts := strings.Split(key, ".")
		current := result
		for _, part := range parts[:len(parts)-1] {
			if _, ok := current[part]; !ok {
				current[part] = make(map[string]interface{})
			}
			current = current[part].(map[string]interface{})
		}
		current[parts[len(parts)-1]] = nil
	}
	return result
}

func buildNestedRecord(row []string, headers []string) map[string]interface{} {
	record := make(map[string]interface{})
	for i, header := range headers {
		parts := strings.Split(header, ".")
		current := record
		for _, part := range parts[:len(parts)-1] {
			if _, ok := current[part]; !ok {
				current[part] = make(map[string]interface{})
			}
			current = current[part].(map[string]interface{})
		}
		val := ""
		if i < len(row) {
			val = row[i]
		}
		current[parts[len(parts)-1]] = parseValue(val)
	}
	return record
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run csv_parser.go <input.csv> [output.json]")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	outputFile := "output.json"
	if len(os.Args) > 2 {
		outputFile = os.Args[2]
	}

	fmt.Println("📊 CSV → JSON Парсер (Go)")

	file, err := os.Open(inputFile)
	if err != nil {
		fmt.Printf("❌ Ошибка: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	delim := detectDelimiter(inputFile)
	fmt.Printf("🔍 Определён разделитель: %q\n", delim)

	reader := csv.NewReader(file)
	reader.Comma = delim
	reader.FieldsPerRecord = -1

	headers, err := reader.Read()
	if err != nil {
		fmt.Printf("❌ Ошибка чтения заголовков: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("📋 Заголовки: %v\n", headers)

	var data []map[string]interface{}
	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Printf("⚠️ Ошибка строки: %v\n", err)
			continue
		}
		record := make(map[string]interface{})
		for i, header := range headers {
			val := ""
			if i < len(row) {
				val = row[i]
			}
			record[header] = parseValue(val)
		}
		data = append(data, record)
	}

	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		fmt.Printf("❌ Ошибка JSON: %v\n", err)
		os.Exit(1)
	}

	err = os.WriteFile(outputFile, jsonData, 0644)
	if err != nil {
		fmt.Printf("❌ Ошибка сохранения: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✅ Экспортировано %d записей в %s\n", len(data), outputFile)
}
