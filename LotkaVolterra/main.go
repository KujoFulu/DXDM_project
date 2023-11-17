package main

import (
	"fmt"
)

func main() {
	fmt.Println("Simulation of LV model starts!")

	fmt.Println("Reading input parameters...")

	// initialize number of species, interaction matrix, and deathGrowth matrix
	numSpecies := 3

	// population matrix
	// // 1. orginal paper
	// pop := []float64{50.0, 10.0, 5.0}

	// 2. stable equilibrium & limit cycle
	pop := []float64{0.1, 0.8, 0.3}

	// // method 1: randomly generate
	// interaction := InitializeInteractionMatrix(numSpecies)
	// method 2: set a matrix by hand

	// // 1. original paper parameters
	// interactionSlice := []float64{0, -0.04, -0.04, 0.04, 0, -0.02, 0.02, 0.04, 0}

	// // 2. stable equilibrium
	// interactionSlice := []float64{-2, -1, 0, 0, -1, -2, -2.6, -1.6, -3}

	// 3. limit cycle
	interactionSlice := []float64{-0.5, -1, 0, 0, -1, -2, -2.6, -1.6, -3}

	interaction := SetInteractionMatrix(interactionSlice, numSpecies)

	// // method 1: randomly generate
	// deathGrowth := IniRateMatrix(numSpecies)

	// method 2: set a matrix by hand

	// // 1. original paper parameters
	// rateSlice := []float64{0.25, -0.5, -0.5}

	// 2. stable equilibrium & limit cycle
	rateSlice := []float64{3, 4, 7.2}

	deathGrowth := SetRateMatrix(rateSlice)

	// initialize number of generations and time interval
	numGens := 50000
	time := 0.002

	fmt.Println("parameters read! Initilizing ecosystem...")

	// initialize an Ecosystem object
	initialEcosystem := InitializeEcosystem(numSpecies, pop, interaction, deathGrowth)

	fmt.Println("Ecosystem initialized! Simulating ecosystem...")

	timePoints := SimulateEcosystem(initialEcosystem, numGens, time)

	fmt.Println("Ecosystem simulated! Printing results...")

	// print out the population of each species in the last ecosystem
	fmt.Println("The population of the first specie in the last ecosystem is:", timePoints[numGens].species[0].population)

	// writing data to csv file
	fmt.Println("Writing data to csv file...")
	WriteToCSV(timePoints, "output/data_limit_cycle.csv")
	fmt.Println("Data written to csv file!")

}
