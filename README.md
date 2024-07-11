# Processador de CSV

# Como Criar e Utilizar uma Biblioteca Compartilhada em Go e C

Este guia explica como criar uma biblioteca compartilhada (`.so`) usando Go, e como utilizá-la em um programa C. Vamos demonstrar isso com um exemplo que processa arquivos CSV.

## Objetivo

A implementação da biblioteca _shared object_ (.so) que processe arquivos CSV, aplicando filtros e selecionando colunas conforme especificado, é capaz de integrar com a interface definida em C abaixo.

```c
/**
 * Process the CSV data by applying filters and selecting columns.
 *
 * @param csv The CSV data to be processed.
 * @param selectedColumns The columns to be selected from the CSV data.
 * @param rowFilterDefinitions The filters to be applied to the CSV data.
 *
 * @return void
 */
void processCsv(const char[], const char[], const char[]);

/**
 * Process the CSV file by applying filters and selecting columns.
 *
 * @param csvFilePath The path to the CSV file to be processed.
 * @param selectedColumns The columns to be selected from the CSV data.
 * @param rowFilterDefinitions The filters to be applied to the CSV data.
 *
 * @return void
 */
void processCsvFile(const char[], const char[], const char[]);
```

## Compilação do código Go para uma _shared object_ - Caos de uso.

### Pass 1: 

Use o comando `go build` para compilar o código fonte em uma biblioteca compartilhada (`.so`).

```sh
go build -buildmode=c-shared -o libcsv.so main.go
```

### Passo 2: 

Escreva um programa C que utilize a biblioteca compartilhada Go. Vamos assumir que o código está em um arquivo chamado `main.c`.

```c
#include <stdio.h>
#include "libcsv.h"

int main() {
    char csv[] = "col1,col2,col3\nl1c1,l1c2,l1c3\nl2c1,l2c2,l2c3\nl3c1,l3c2,l3c3";
    char selectedColumns[] = "col1,col3";
    char rowFilterDefinitions[] = "col2=l2c2";
    
    printf("processCsv output:\n");
    processCsv(csv, selectedColumns, rowFilterDefinitions);

    char csvFilePath[] = "data.csv";
    printf("\nprocessCsvFile output:\n");
    processCsvFile(csvFilePath, selectedColumns, rowFilterDefinitions);

    return 0;
}
```
### Passo 3: 

Use o comando `gcc` para compilar o programa C e linkar com a biblioteca compartilhada.

```sh
gcc -o test main.c ./libcsv.so
```
### Passo 4: 

Execute o programa compilado.

```sh
./test
```

### Saída Esperada

A execução do programa deve gerar a seguinte saída:

```
processCsv output:
col1,col3,col4,col7
l2c1,l2c3,l2c4,l2c7,
l3c1,l3c3,l3c4,l3c7,

processCsvFile output:
col1,col3,col4,col7
l2c1,l2c3,l2c4,l2c7,
l3c1,l3c3,l3c4,l3c7,
```
---
