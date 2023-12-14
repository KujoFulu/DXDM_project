package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)


type Population struct {
	popSize   int     //Population size
	gen       int     //Gen which generation this is (we might not need this)
	selCo     float64 //Selection coefficient
	freqStart float64 //Starting allele frequency
	freqNum   float64 //Number of allelles
	freq      float64 //K This pop’s allele frequency
}


func main() {

	
	// Initialize variables
	var popSize int
	var selCo float64
	var freqStart float64
	var numGen int
	var numRuns int


	//the first parameter is population size
	popSize, err1 := strconv.Atoi(os.Args[1])
	if err1 != nil {
		//problem in converting this parameter
		panic(err1)
	}

	if popSize < 0 {
		panic("Error: negative number given as number of popSize.")
	}

	//the second parameter is Selection coefficient
	selCo, err2 := strconv.ParseFloat(os.Args[2], 64)
	if err2 != nil {
		//problem in converting this parameter
		panic(err2)
	}

	//the third parameter is Starting allele frequency
	freqStart, err3 := strconv.ParseFloat(os.Args[3], 64)
	if err3 != nil {
		//problem in converting this parameter
		panic(err3)
	}

	//the fourth parameter is number of generations
	numGen, err4 := strconv.Atoi(os.Args[4])
	if err4 != nil {
		//problem in converting this parameter
		panic(err4)
	}

	if numGen < 0 {
		panic("Error: negative number given as number of numGen.")
	}

	//the fifth parameter is number of runs
	numRuns, err5 := strconv.Atoi(os.Args[5])
	if err1 != nil {
		//problem in converting this parameter
		panic(err5)
	}

	if numRuns < 0 {
		panic("Error: negative number given as number of numRuns.")
	}


	// Print loaded parameters
	fmt.Println("Population size =", popSize)
	fmt.Println("Select coefficient =", selCo)
	fmt.Println("Start allele frequency =", freqStart)
	fmt.Println("Number of simulating generations =", numGen)
	fmt.Println("Number of simulation runs =", numRuns)
	fmt.Println("All parameters loaded!")


	// Specify the folder name (for R plotting)
	// folderName := "WrightFisher"
	filename := "SimulationParameters.csv"


	// Write parameters to CSV file
	WriteParameters(popSize, selCo, freqStart, numGen, numRuns, filename)
	fmt.Println("SimulationParameters.csv file created")


	// Run simulations using SimulateMultipleRuns
	fmt.Println("Start simulation!")
	startTime := time.Now()
	runs := SimulateMultipleRuns(numRuns, popSize, numGen, selCo, freqStart)
	log.Println("Runtime:", time.Since(startTime))
	fmt.Println("Simulation done, start output data")


	// Combine data from all runs into a single slice
	var allData []*Population
	for _, timePoints := range runs {
		allData = append(allData, timePoints...)
	}

	// Write all data to a single CSV file
	WriteToCSV(allData, "all_simulation_data.csv")
	fmt.Println("Data output successfully!")
	

	fmt.Println("Start simulate two loci model ")
    recomb := 0.1
    twpLoci := SimulateTwoLoci(popSize, selCo, recomb)
    WritetwoLToCSV(twpLoci, "two_loci_simulation_data")


}
