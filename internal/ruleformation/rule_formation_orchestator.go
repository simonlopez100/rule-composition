package ruleformation

import (
	// "encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"slices"
	"strconv"
	"strings"

	readercsv "github.com/mercadolibre/rule-formation/internal/readcsv"
)

type BigQueryTable struct {
	GroupRule      string
	NoRule         int
	OrderExec      int
	OrderGroup     int
	Process        string
	ResultantQuery string
}

var activeProcesses = []string{
	"ACQUISITION",
	"SHUTDOWN",
	"REVALUATION",
	"REACTIVATION",
	"CHANGESHUTDOWN",
}

var activeSubProcesses = map[string]func(ad []readercsv.RowData) ([]BigQueryTable, error){
	"AUDIENCE":        AudienceSubProcess,
	"CROSS_FILTERS":   CrossFiltersSubProcess,
	"CATEGORY_FILTERS": CategoryFiltersSubProcess,
	"ACCESS":          AccessSubProcess,
	"PRICING":         PricingSubProcess,
	"PROFILE":         ProfileSubProcess,
}

var RuleCategoryNameTranslation = map[string]int{
	"AUDIENCE":         1,
	"CROSS_FILTERS":    2,
	"CATEGORY_FILTERS": 3,
	"ACCESS":           4,
	"PRICING":          5,
	"PROFILE":          6,
}

func RuleFormationOrchestator(data []readercsv.RowData) (bqArray []BigQueryTable, err error) {

	partitions, err := PartitionByField(data, "ProcessName")
	if err != nil {
		log.Fatalf("Error al particionar los datos: %v", err)
	}

	for key, process := range partitions.(map[string]interface{}) {
		if len(process.([]readercsv.RowData)) == 0 {
			continue
		}

		if !slices.Contains(activeProcesses, strings.ToUpper(key)){
			continue
		}

		result, err := ProcessHandler(process.([]readercsv.RowData))
		if err != nil {
			log.Fatalf("Error al procesar el proceso %s: %v", key, err)
			continue
		}

		bqArray = append(bqArray, result...)
	}

	return bqArray, nil

}

func ProcessHandler(data []readercsv.RowData) (bqArray []BigQueryTable, err error) {
	principalProcces := data[0].ProcessName

	// particionar los subprocessos por RuleCategoryName
	partitions, err := PartitionByField(data, "RuleCategoryName")
	if err != nil {
		log.Fatalf("Error al particionar los datos: %v", err)
	}

	subProcessFunctiones := make(map[string]func(ad []readercsv.RowData) ([]BigQueryTable, error))
	for key := range partitions.(map[string]interface{}) {

		processFunc, ok := activeSubProcesses[strings.ToUpper(key)]
		if !ok {
			continue
		}

		//validate funcion isnt exists in map
		_, ok = subProcessFunctiones[key]
		if ok {
			continue
		}
		subProcessFunctiones[key] = processFunc
	}

	process := NewProcess(principalProcces, data, subProcessFunctiones)
	result, err := process.ExecuteProcess()
	if err != nil {
		log.Fatalf("Error al procesar el proceso %s: %v", principalProcces, err)
	}

	bqArray = append(bqArray, result...)

	return bqArray, nil
}

func PartitionByField(array interface{}, field string) (interface{}, error) {
	val := reflect.ValueOf(array)

	// Validar que el input sea un slice
	if val.Kind() != reflect.Slice {
		return nil, errors.New("input must be a slice")
	}

	// Crear un mapa de particiones
	partitions := make(map[string]interface{})

	for i := 0; i < val.Len(); i++ {
		row := val.Index(i)
		fieldVal := row.FieldByName(field)

		if !fieldVal.IsValid() {
			return nil, fmt.Errorf("field not defined: %s", field)
		}

		var key string
		switch fieldVal.Kind() {
		case reflect.String:
			key = fieldVal.String()
		case reflect.Int:
			key = strconv.Itoa(int(fieldVal.Int()))
		default:
			return nil, fmt.Errorf("unsupported field type: %s", fieldVal.Kind())
		}

		if _, exists := partitions[key]; !exists {
			partitions[key] = reflect.MakeSlice(reflect.SliceOf(row.Type()), 0, 0).Interface()
		}

		partitions[key] = reflect.Append(reflect.ValueOf(partitions[key]), row).Interface()
	}

	return partitions, nil
}

func formatBigQueryTable(row readercsv.RowData, queryResult string) (BigQueryTable, error) {
	return BigQueryTable{
			GroupRule: row.RuleCategoryName,
			NoRule: row.LineIndex,
			OrderExec: row.Order,
			Process: row.ProcessName,
			OrderGroup: RuleCategoryNameTranslation[strings.ToUpper(row.RuleCategoryName)],
			ResultantQuery: queryResult,
	}, nil
}

