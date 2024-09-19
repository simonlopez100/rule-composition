package reader

import (
	"encoding/csv"
	"fmt"
	"strconv"
)

type RowData struct {
	AssignmentResult  string `json:"Assigment_Result"`
	LineIndex         int
	Order             int    `json:"Order"`
	ProcessName       string `json:"Process Name"`
	RuleCategoryName  string `json:"Rule Category Name"`
	RuleConditionName string `json:"Rule Condition Name"`
	RuleName          string `json:"Rule Name"`
	RuleType          string `json:"Rule Type"`
	SentenceResult    string `json:"Sentence Result"`
}

func ReadCSV(reader *csv.Reader, headerIndex map[string]int, expectedHeaders []string) ([]RowData, error) {
	// Verificar que todas las columnas esperadas están presentes
	for _, expected := range expectedHeaders {
		if _, ok := headerIndex[expected]; !ok {
			return nil, fmt.Errorf("falta la columna esperada: %s", expected)
		}
	}

	var result []RowData

	// Leer todas las filas
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("error al leer el archivo CSV: %w", err)
	}

	// Iterar sobre las filas y almacenar solo las columnas de interés
	for i, row := range rows {
		order, err := strconv.Atoi(row[headerIndex["Order"]])
		if err != nil {
			return nil, fmt.Errorf("error al convertir Order a entero: %w", err)
		}

		rowData := RowData{
			LineIndex:         i + 1,
			ProcessName:       row[headerIndex["Process Name"]],
			RuleCategoryName:  row[headerIndex["Rule Category Name"]],
			Order:             order,
			RuleName:          row[headerIndex["Rule Name"]],
			RuleConditionName: row[headerIndex["Rule Condition Name"]],
			RuleType:          row[headerIndex["Rule Type"]],
			SentenceResult:    row[headerIndex["Sentence Result"]],
			AssignmentResult:  row[headerIndex["Assignment Result"]],
		}
		result = append(result, rowData)
	}

	return result, nil
}
