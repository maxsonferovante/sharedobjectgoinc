#include "libcsv.h"

/// @brief 
/// @return 
int main(){

    /* csv := "header1,header2,header3\n1,2,3\n4,5,6\n7,8,9"
    processCsv(csv, "header1,header3", "header1>1\nheader3<8")
    // output
    // header1,header3
    // 4,6
    */
    
    char csv[] = "header1,header2,header3\n1,2,3\n4,5,6\n7,8,9";
    char selectedColumns[] = "header1,header3";
    char rowFilterDefinitions[] = "header1>1\nheader3<8";

    processCsv(csv, selectedColumns, rowFilterDefinitions);
    
    return 0;
}