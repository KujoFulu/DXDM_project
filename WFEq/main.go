package main

import (
	"encoding/csv"
	"fmt"
	"os"
	
)

func main() {

startPop:= InitializePopulation(100, -0.1, 0.3)
var populations []*Population
// populations := the final function
populations = SimulatePopulationTimePoints(startPop, 10)

// Set the output file name
filePath := "WrightFishe_output.csv"


// creat file path
file, err := os.Create(filePath)
if err != nil {
fmt.Println("Error creating file:", err)
return
}
defer file.Close()


//writing the file
writer := csv.NewWriter(file)
defer writer.Flush()


// set hearder of the output table
header := []string{"Population_Size", "Generations", "Selection_coefficient", "Starting_allele_frequency", "Allele_frequency"}
if err := writer.Write(header); err != nil {
fmt.Println("Failed to write header:", err)
return
}


// each row contain one parameter
for _, p := range populations {
row := []string{
fmt.Sprintf("%d", p.popSize),
fmt.Sprintf("%d", p.gen),
fmt.Sprintf("%.2f", p.selCo),
fmt.Sprintf("%.2f", p.freqStart),
fmt.Sprintf("%.2f", p.freq),
}
if err := writer.Write(row); err != nil {
fmt.Println("Failed to write row:", err)
return
}
}


fmt.Println("CSV file generation successful.")
}


