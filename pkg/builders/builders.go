package builders

import "C"

import (
	"errors"
	"fmt"
	"libcsv/pkg/constants"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/Knetic/govaluate"
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

func BuilderFilteredMatrix(matrix map[string][]string, rowFilterDefinitions string) map[string][]string {

	expressions, err := BuilderFilters(rowFilterDefinitions)

	if err != nil {
		message := fmt.Sprintf("Error: %s", err)
		err := errors.New(message)
		panic(err)
	}

	maxtrixFiltered := make(map[string][]string)
	for i := range matrix[expressions[0][0]] {
		validRow := true
		for _, expression := range expressions {
			value := matrix[expression[0]][i]

			valueInt, err := strconv.Atoi(value)

			if err != nil {
				message := fmt.Sprintf("Error: %s", err)
				err := errors.New(message)
				panic(err)
			}
			paramenters := map[string]interface{}{
				expression[0]: valueInt,
			}

			evaluableExpression, _ := govaluate.NewEvaluableExpression(expression[1])

			result, err := evaluableExpression.Evaluate(paramenters)

			if err != nil {
				panic(err)
			}
			if result == false {
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
	return maxtrixFiltered
}

func BuilderFilters(rowFilterDefinitions string) ([][2]string, error) {
	var expressions [][2]string

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

				expression := [2]string{column, fmt.Sprintf("%s%s%s", column, operator, value)}
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

	// Prepara e imprime as linhas de dados
	rows := make([]string, len(headers))
	for i, header := range headers {
		rows[i] = strings.Join(matrix[header], ",")
	}

	processed := strings.Join(headers, ",") + "\n" + strings.Join(rows, ",") + "\n"
	os.Stdout.WriteString(processed)
	os.Stdout.Sync()
}
