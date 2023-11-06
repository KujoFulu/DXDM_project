package main

import (
	"fmt"
)

func main() {
	fmt.Println("Simulation of LV model starts!")

	fmt.Println("Reading input parameters...")

	// initialize number of species, interaction matrix, and deathGrowth matrix
	numSpecies := 3

	// // method 1: randomly generate
	// interaction := InitializeInteractionMatrix(numSpecies)
	// method 2: set a matrix by hand
	interactionSlice := []float64{0, -0.04, -0.04, 0.04, 0, -0.02, 0.02, 0.04, 0}
	interaction := SetInteractionMatrix(interactionSlice, numSpecies)

	// // method 1: randomly generate
	// deathGrowth := IniRateMatrix(numSpecies)
	// method 2: set a matrix by hand
	rateSlice := []float64{0.25, -0.5, -0.5}
	deathGrowth := SetRateMatrix(rateSlice)

	// initialize number of generations and time interval
	numGens := 50000
	time := 0.002

	fmt.Println("parameters read! Initilizing ecosystem...")

	// initialize an Ecosystem object
	initialEcosystem := InitializeEcosystem(numSpecies, interaction, deathGrowth)

	fmt.Println("Ecosystem initialized! Simulating ecosystem...")

	timePoints := SimulateEcosystem(initialEcosystem, numGens, time)

	fmt.Println("Ecosystem simulated! Printing results...")

	// print out the population of each species in the last ecosystem
	fmt.Println("The population of the first specie in the last ecosystem is:", timePoints[numGens].species[0].population)

}
