package main

import (
	"fmt"
	"log"
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

	// Set parameters for simulation
	popSize := 150   // Set population size
	selCo := 0.0015  // Set select coefficient
	freqStart := 0.6 // Set the start allele frequency
	numGen := 200    // Set the number of generations
	numRuns := 200   // Set the number of simulation runs

	fmt.Println("Population size =", popSize)
	fmt.Println("Select coefficient =", selCo)
	fmt.Println("Start allele frequency =", freqStart)
	fmt.Println("Number of simulating generations =", numGen)
	fmt.Println("Number of simulation runs =", numRuns)
	fmt.Println("All parameters loaded!")

	// Specify the folder name(for R plotting)
	//folderName := "WrightFisher"
	filename := "SimulationParameters.csv"

	// Write parameters to csv file
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

}
