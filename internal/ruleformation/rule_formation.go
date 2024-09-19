package ruleformation

import (
	"fmt"
	"strings"

	reader "github.com/mercadolibre/rule-formation/internal/readcsv"
)

const (
	updateQuery   = "UPDATE `TABLA_ORIGEN` SET  %s WHERE %s"
	selectQuery = "SELECT * FROM `TBL.MANTIS` WHERE %s"
	groupRuleAudience = "AUDIENCE"
	groupRuleCrossFilters = "CROSS_FILTERS"
	groupRuleCategoryFilters = "CATEGORY_FILTERS"
	groupRuleAccess = "ACCESS"
	groupRulePricing = "PRICING"
	groupRuleProfile = "PROFILE"
	processGroupAcquisition = "ACQUISITION"
	acquisitionFilterFlag = "FILTER_FLAG=1"
	subProcessessFilter = "RuleCategoryName"
)

type Process struct {
	ProcessName       string
	ProcessRows       []reader.RowData
	ProccessFunctions map[string]func([]reader.RowData) ([]BigQueryTable, error)
}

//To-do update function generic

//to-do select function generic

func NewProcess(processName string, processRows []reader.RowData, processFunctions map[string]func([]reader.RowData) ([]BigQueryTable, error)) *Process {
	return &Process{
		ProcessName:       processName,
		ProcessRows:       processRows,
		ProccessFunctions: processFunctions,
	}
}

func (p *Process) ExecuteProcess() (result []BigQueryTable, err error) {
	//execute processes into processFunctions and append into result
	for _, processFunc := range p.ProccessFunctions {
		processResult, err := processFunc(p.ProcessRows)
		if err != nil {
			return nil, err
		}
		result = append(result, processResult...)
	}
	return result, nil
}

func AudienceSubProcess(ad []reader.RowData) (bg []BigQueryTable, err error) {
	subProccess, err := PartitionByField(ad, subProcessessFilter)
	if err != nil {
		return nil, err
	}
	audienceProcess, err := getProcessArray(subProccess, groupRuleAudience)
	if err != nil {
		return nil, err
	}

	for _, row := range audienceProcess {
		queryResult := createDynamicQuery(row)
		BigQueryObject, err := formatBigQueryTable(row, queryResult)
		if err != nil {
			return nil, err
		}
		bg = append(bg, BigQueryObject)
	}

	return
}

func CrossFiltersSubProcess(ad []reader.RowData) (bg []BigQueryTable, err error) {
		subProccess, err := PartitionByField(ad, subProcessessFilter)
	if err != nil {
		return nil, err
	}
	crossFiltersProcces, err := getProcessArray(subProccess, groupRuleCrossFilters)
	if err != nil {
		return nil, err
	}

	for _, row := range crossFiltersProcces {
		if strings.ToUpper(row.ProcessName) == processGroupAcquisition {
			row.AssignmentResult = acquisitionFilterFlag
		}
		queryResult := createDynamicQuery(row)
		BigQueryObject, err := formatBigQueryTable(row, queryResult)
		if err != nil {
			return nil, err
		}
		bg = append(bg, BigQueryObject)
	}

	return
}

func CategoryFiltersSubProcess(ad []reader.RowData) (bg []BigQueryTable, err error) {
		subProccess, err := PartitionByField(ad, subProcessessFilter)
	if err != nil {
		return nil, err
	}
	categoryFiltersProcces, err := getProcessArray(subProccess, groupRuleCategoryFilters)
	if err != nil {
		return nil, err
	}



	for _, row := range categoryFiltersProcces {
		if strings.ToUpper(row.ProcessName) == processGroupAcquisition {
			row.AssignmentResult = acquisitionFilterFlag
		}
		queryResult := createDynamicQuery(row)
		BigQueryObject, err := formatBigQueryTable(row, queryResult)
		if err != nil {
			return nil, err
		}
		bg = append(bg, BigQueryObject)
	}

	return
}

func AccessSubProcess(ad []reader.RowData) (bg []BigQueryTable, err error) {
		subProccess, err := PartitionByField(ad, subProcessessFilter)
	if err != nil {
		return nil, err
	}
	accessProcess, err := getProcessArray(subProccess, "ACCESS")
	if err != nil {
		return nil, err
	}

	for _, row := range accessProcess {
		queryResult := createDynamicQuery(row)
		BigQueryObject, err := formatBigQueryTable(row, queryResult)
		if err != nil {
			return nil, err
		}
		bg = append(bg, BigQueryObject)
	}

	return
}

func PricingSubProcess(ad []reader.RowData) (bg []BigQueryTable, err error) {
	subProccess, err := PartitionByField(ad, subProcessessFilter)
	if err != nil {
		return nil, err
	}
	pricingProcess, err := getProcessArray(subProccess, groupRulePricing)
	if err != nil {
		return nil, err
	}

	for _, row := range pricingProcess {
		queryResult := createDynamicQuery(row)
		BigQueryObject, err := formatBigQueryTable(row, queryResult)
		if err != nil {
			return nil, err
		}
		bg = append(bg, BigQueryObject)
	}

	return
}

func ProfileSubProcess(ad []reader.RowData) (bg []BigQueryTable, err error) {
		subProccess, err := PartitionByField(ad, subProcessessFilter)
	if err != nil {
		return nil, err
	}
	audienceProcess, err := getProcessArray(subProccess, groupRuleProfile)
	if err != nil {
		return nil, err
	}

	for _, row := range audienceProcess {
		queryResult := createDynamicQuery(row)
		BigQueryObject, err := formatBigQueryTable(row, queryResult)
		if err != nil {
			return nil, err
		}
		bg = append(bg, BigQueryObject)
	}

	return
}

func getProcessArray(partitions interface{}, subProcess string)([]reader.RowData, error) {
		for key, process := range partitions.(map[string]interface{}) {
		if len(process.([]reader.RowData)) == 0 {
			continue
		}
		if key == subProcess {
			return process.([]reader.RowData), nil
		}
	}
	return nil, nil
}



func createDynamicQuery(row reader.RowData) string {
	if strings.ToUpper(row.RuleCategoryName) == groupRuleAudience {
		return fmt.Sprintf(selectQuery, row.SentenceResult)
	}
	return fmt.Sprintf(updateQuery, row.AssignmentResult, row.SentenceResult)
}