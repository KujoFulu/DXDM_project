package main

import (
	"fmt"
)

func main() {
	fmt.Println("Simulation of LV model starts!")

	// initialize number of species, interaction matrix, and deathGrowth matrix
	numSpecies := 3
	interaction := InitializeInteractionMatrix(numSpecies)
	deathGrowth := IniRateMatrix(numSpecies)

	// initialize an Ecosystem object
	initialEcosystem := InitializeEcosystem(numSpecies, interaction, deathGrowth)

	// initialize number of generations and time interval
	numGens := 100
	time := 0.1

	timePoints := SimulateEcosystem(initialEcosystem, numGens, time)

	// print out the population of each species in the last ecosystem
	fmt.Println("The population of the first specie in the last ecosystem is:", timePoints[numGens].species[0].population)

}
