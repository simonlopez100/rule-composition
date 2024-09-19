package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"os"

	readercsv "github.com/mercadolibre/rule-formation/internal/readcsv"
	ruleformation "github.com/mercadolibre/rule-formation/internal/ruleformation"
)

func main() {
	filePath := "data/example.csv"

	// Abrir el archivo .csv
	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("Error al abrir el archivo: %v", err)
	}
	defer file.Close()

	// Crear un lector CSV
	reader := csv.NewReader(file)

	// Leer el encabezado
	headers, err := reader.Read()
	if err != nil {
		log.Fatalf("Error al leer el encabezado del archivo: %v", err)
	}

	// Crear un mapa de índices de encabezado
	headerIndex := make(map[string]int)
	for i, header := range headers {
		headerIndex[header] = i
		// fmt.Printf("Header: %s, Index: %d\n", header, i)
	}

	// Definir los encabezados esperados
	expectedHeaders := []string{
		"Process Name", "Rule Category Name", "Order", "Rule Name",
		"Rule Condition Name", "Rule Type", "Sentence Result", "Assignment Result",
	}

	// Leer el archivo .csv y almacenarlo en una estructura
	data, err := readercsv.ReadCSV(reader, headerIndex, expectedHeaders)
	if err != nil {
		log.Fatalf("Error al leer el archivo: %v", err)
	}

	// Procesar los datos y agregar un ID único a cada fila
	result, err := ruleformation.RuleFormationOrchestator(data)
	if err != nil {
		log.Fatalf("Error al procesar los datos: %v", err)
	}

	// fmt.Printf("%+v\n", result)

	// // Particionar los datos por el campo "Process"
	// partitions, err := ruleformation.PartitionByField(result, "Process")
	// if err != nil {
	// 	log.Fatalf("Error al particionar los datos: %v", err)
	// }

	// // Convertir el resultado a JSON para una impresión más legible
	jsonResult, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		log.Fatalf("Error al convertir el resultado a JSON: %v", err)
	}

	// Imprimir el resultado en formato JSON
	fmt.Println(string(jsonResult))
}
