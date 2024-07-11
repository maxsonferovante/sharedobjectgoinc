# Processador de CSV

## Execução de teste

```bash
docker build -t image_csv . 
```

```bash
docker run --rm -it image_csv
```
## Objetivo

Implementa uma biblioteca _shared object_ (.so) que processe arquivos CSV, aplicando filtros e selecionando colunas conforme especificado. A solução deve ser capaz de integrar com a interface definida em C abaixo.

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


## Exemplo:

```c
const char csv[] = "header1,header2,header3\n1,2,3\n4,5,6\n7,8,9";
processCsv(csv, "header1,header3", "header1>1\nheader3<8");

// output
// header1,header3
// 4,6

const char csv_file[] = "path/to/csv_file.csv";
processCsvFile(csv_file, "header1,header3", "header1>1\nheader3<8");

// output
// header1,header3
// 4,6
```

---