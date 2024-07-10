package valitations

import "C"

import (
	"errors"
	"fmt"
	"libcsv/pkg/constants"
	"sort"
	"strings"
)

func ValidateCSV(data string) error {
	// Verifica se o dado contém uma vírgula
	if strings.Contains(data, ",") {
		return nil
	}

	// Divide a string em uma lista de substrings separadas por vírgulas
	parts := strings.Split(data, ",")

	// Verifica se a lista resultante tem mais de um elemento
	if len(parts) > 1 {
		return nil
	}

	// Verifica se a lista tem mais de 256 elementos
	if len(parts) > 256 {
		return nil
	}

	// Se nenhuma das condições acima for atendida, retorna um erro
	return errors.New("CSV data is invalid")
}

func ValidationSelectedColumns(matrix map[string][]string, selectedColumns string) (error, []string) {
	if selectedColumns == "" {
		return nil, nil
	}

	// Convert map keys to a list of headers
	headersList := make([]string, 0, len(matrix))
	for key := range matrix {
		headersList = append(headersList, key)
	}

	// ordena os cabeçalhos em ordem crescente
	sort.Strings(headersList)
	selectedColumnsList := strings.Split(selectedColumns, ",")

	// Verifica se as colunas selecionadas existe no cabeçalho
	for _, column := range selectedColumnsList {
		if !validationExistsHeader(column, headersList) {
			message := fmt.Sprintf("Header '%s' not found in csv file/string", column)
			err := errors.New(message)
			return err, nil
		}
	}

	headerIndex := 0
	for _, column := range selectedColumnsList {
		for headerIndex < len(headersList) && headersList[headerIndex] != column {
			headerIndex++
		}
		if headerIndex == len(headersList) {
			message := fmt.Sprintf("the order of the selection string does not match that of the csv: %s", column)
			err := errors.New(message)
			return err, nil
		}
		headerIndex++
	}
	return nil, nil
}
func ValidationFilterDefinitions(headers map[string][]string, rowFilterDefinitions string) (error, []string) {
	filterDefinitionsList := strings.Split(rowFilterDefinitions, "\n")
	headersList := make([]string, 0, len(headers))

	// Collect all headers into a list
	for key := range headers {
		headersList = append(headersList, key)
	}

	filterColumns := make([]string, 0, len(filterDefinitionsList))

	sort.Strings(headersList)
	sort.Strings(filterColumns)
	sort.Strings(filterDefinitionsList)

	// Extract columns from filter definitions and validate their existence and order
	for _, filterDefinition := range filterDefinitionsList {
		column, err := extractColumnFromOrderDefinition(filterDefinition)
		if err != nil {
			message := fmt.Sprintf("Error: %s", err)
			err := errors.New(message)
			return err, nil
		}

		// Check if the column exists in the CSV headers
		if !validationExistsHeader(column, headersList) {
			message := fmt.Sprintf("Header '%s' not found in csv file/string", column)
			err := errors.New(message)
			return err, nil
		}

		// Check if the column is already in the filter columns list
		exists := func(column string) bool {
			for _, header := range filterColumns {
				if header == column {
					return true
				}
			}
			return false
		}
		if !exists(column) {
			filterColumns = append(filterColumns, column)
		}
	}

	// Validate the order of filter definitions against CSV headers
	headerIndex := 0
	for _, column := range filterColumns {
		for headerIndex < len(headersList) && headersList[headerIndex] != column {
			headerIndex++
		}
		if headerIndex == len(headersList) {
			message := fmt.Sprintf("The order of the filter definitions does not match that of the csv: %s", column)
			err := errors.New(message)
			return err, nil
		}
		headerIndex++
	}
	return nil, nil
}

func extractColumnFromOrderDefinition(orderDefinition string) (string, error) {
	for _, operator := range constants.Operators {
		if strings.Contains(orderDefinition, operator) {
			column := strings.Split(orderDefinition, operator)[0]
			return strings.TrimSpace(column), nil
		}
	}
	message := fmt.Sprintf("the filter definition is invalid: %s", orderDefinition)
	err := errors.New(message)
	panic(err)
}

func validationExistsHeader(column string, headersList []string) bool {
	for _, header := range headersList {
		if header == column {
			return true
		}
	}
	return false
}
