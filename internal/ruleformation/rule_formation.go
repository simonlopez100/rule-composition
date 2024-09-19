package ruleformation

import (
	"fmt"
	"reflect"
	"strconv"

	readercsv "github.com/mercadolibre/rule-formation/internal/readcsv"
)

type BigQueryTable struct {
	NoRule         int
	Process        string
	OrderGroup     int
	GroupRule      string
	OrderExec      int
	ResultantQuery string
}

var RuleCategoryNameTranslation = map[string]int{
	"AUDCIENCE":        1,
	"CROSS__FILTERS":   2,
	"CATEGORY_FILTERS": 3,
	"ACCESS":           4,
	"PRICING":          5,
	"PROFILE":          6,
}

func RuleFormationOrchestator(data []readercsv.RowData) (bqArray []BigQueryTable, err error) {

	// Inicializar la lista de BigQueryTable, creacion index por fila (NO_RULE)
	for _, row := range data {
		bqTable := BigQueryTable{
			NoRule:         row.LineIndex,
			Process:        row.ProcessName,
			OrderGroup:     row.Order,
			GroupRule:      row.RuleCategoryName,
			OrderExec:      row.Order,
			ResultantQuery: "Aqui va la logica de la query",
		}
		// result[i] = bqTable
		bqArray = append(bqArray, bqTable)
	}

	// Poblar process, ordergroup, grouprule, orderexec

	// Armar query seleccion de tabla  SELECT * FROM `TABLE`

	// Armar query de estructura WHERE

	// Armar query de estructura WHERE con SET

	return bqArray, nil
}

func PartitionByField(array []BigQueryTable, field string) ([][]BigQueryTable, error) {
	// Crear un mapa de particiones
	partitions := make(map[string][]BigQueryTable)

	for _, row := range array {
		var key string
		val := reflect.ValueOf(row)
		fieldVal := val.FieldByName(field)
		if !fieldVal.IsValid() {
			return nil, fmt.Errorf("field not defined: %s", field)
		}

		switch fieldVal.Kind() {
		case reflect.String:
			key = fieldVal.String()
		case reflect.Int:
			key = strconv.Itoa(int(fieldVal.Int()))
		default:
			return nil, fmt.Errorf("unsupported field type: %s", fieldVal.Kind())
		}

		partitions[key] = append(partitions[key], row)
	}

	// Convertir el mapa de particiones a un slice de slices
	var result [][]BigQueryTable
	for _, partition := range partitions {
		result = append(result, partition)
	}
	fmt.Println("Resultado partition", result)

	return result, nil
}
