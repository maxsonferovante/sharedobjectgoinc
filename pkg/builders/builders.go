package builders

import "C"

import (
	"errors"
	"fmt"
	"libcsv/pkg/constants"
	"os"
	"sort"
	"strings"
)

func BuilderMatrix(csv string) map[string][]string {
	data := strings.Split(csv, "\n")
	matrix := make(map[string][]string)
	for i, line := range data {

		if i == 0 {
			headers := strings.Split(line, ",")
			for _, header := range headers {

				if _, exists := matrix[header]; exists {
					message := fmt.Sprintf("Error: Header %s already exists", header)
					err := errors.New(message)
					panic(err)
				}

				if len(matrix) >= 256 {
					message := fmt.Sprintf("Error: The number of headers exceeds the limit of 256: %s", header)
					err := errors.New(message)
					panic(err)
				}
				matrix[header] = make([]string, 0)
			}
		} else {
			values := strings.Split(line, ",")
			if len(values) == 1 && values[0] == "" {
				continue
			}
			headers := strings.Split(data[0], ",")
			for j := range headers {
				if j < len(values) {
					matrix[headers[j]] = append(matrix[headers[j]], values[j])
				}
			}
		}
	}
	return matrix
}

func BuilderSelectedMatrix(matrix map[string][]string, selectedColumns string) map[string][]string {
	if selectedColumns == "" {
		return matrix
	}
	columns := strings.Split(selectedColumns, ",")

	selectedMatrix := make(map[string][]string)
	for _, column := range columns {
		if values, exists := matrix[column]; exists {
			selectedMatrix[column] = values
		}
	}
	return selectedMatrix
}

func BuilderFilteredMatrix(matrix map[string][]string, rowFilterDefinitions string) (map[string][]string, error) {

	expressions, err := BuilderFilters(rowFilterDefinitions)

	if err != nil {
		message := fmt.Sprintf("Error: %s", err)
		err := errors.New(message)
		return nil, err
	}

	maxtrixFiltered := make(map[string][]string)
	for i := range matrix[expressions[0][0]] {
		validRow := true
		for _, expression := range expressions {
			value := matrix[expression[0]][i]

			if !applyOperation(expression[2], value, expression[1]) {
				validRow = false
				break
			}
		}
		if validRow {
			for column, values := range matrix {
				if _, exists := maxtrixFiltered[column]; !exists {
					maxtrixFiltered[column] = []string{}
				}
				maxtrixFiltered[column] = append(maxtrixFiltered[column], values[i])
			}
		}
	}
	return maxtrixFiltered, nil
}

func BuilderFilters(rowFilterDefinitions string) ([][3]string, error) {
	var expressions [][3]string

	filtersDefinitions := strings.Split(rowFilterDefinitions, "\n")

	for _, filterDefinition := range filtersDefinitions {
		var validExpression bool
		for _, operator := range constants.Operators {
			if strings.Contains(filterDefinition, operator) {
				parts := strings.Split(filterDefinition, operator)
				if len(parts) < 2 {
					return nil, errors.New("Invalid filter: " + filterDefinition)
				}
				column := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])

				// Ajusta o operador '=' para '=='
				if operator == "=" {
					operator = "=="
				}

				expression := [3]string{column, fmt.Sprintf("%s%s%s", column, operator, value), operator}
				expressions = append(expressions, expression)
				validExpression = true
				break
			}
		}
		if !validExpression {
			return nil, errors.New("Invalid filter: " + filterDefinition)
		}
	}
	return expressions, nil
}
func BuilderApreciation(matrix map[string][]string) {

	// Ordena as chaves do mapa
	headers := make([]string, 0, len(matrix))
	for header := range matrix {
		headers = append(headers, header)
	}
	sort.Strings(headers)
	os.Stdout.WriteString("\n")
	os.Stdout.WriteString(strings.Join(headers, ",") + "\n")
	for i := 0; i < len(matrix[headers[0]]); i++ {
		for _, header := range headers {
			os.Stdout.WriteString(matrix[header][i] + ",")
		}
		os.Stdout.WriteString("\n")
	}
	os.Stdout.WriteString("\n")
}

func contains(values []string, value string) bool {
	for _, v := range values {
		if v == value {
			return true
		}
	}
	return false
}

func applyOperation(operator, value1, value2 string) bool {
	switch operator {
	case ">":
		value2 = strings.Split(value2, ">")[1]
		return value1 > value2
	case "<":
		value2 = strings.Split(value2, "<")[1]
		return value1 < value2
	case ">=":
		value2 = strings.Split(value2, ">=")[1]
		return value1 >= value2
	case "<=":
		value2 = strings.Split(value2, "<=")[1]
		return value1 <= value2
	case "==":
		value2 = strings.Split(value2, "==")[1]
		return value1 == value2
	case "!=":
		value2 = strings.Split(value2, "!=")[1]
		return value1 != value2
	}
	return false
}
