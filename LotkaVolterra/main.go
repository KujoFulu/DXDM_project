package main

import (
	"fmt"
	"gifhelper"
)

func main() {
	fmt.Println("Simulation of LV model starts!")

	fmt.Println("Reading input parameters...")

	// initialize number of species, interaction matrix, and deathGrowth matrix
	numSpecies := 3

	// population matrix
	// // 1. orginal paper
	// pop := []float64{50.0, 10.0, 5.0}

	// // 2. stable equilibrium & limit cycle & extinction
	// pop := []float64{0.1, 0.8, 0.3}

	// 3. chaotic dynamics - 4 species
	pop := []float64{0.1, 0.8, 0.3, 0.5}

	// // method 1: randomly generate
	// interaction := InitializeInteractionMatrix(numSpecies)
	// method 2: set a matrix by hand

	// 1. original paper parameters
	interactionSlice := []float64{0, -0.04, -0.04, 0.04, 0, -0.02, 0.02, 0.04, 0}

	// // 2. stable equilibrium
	// interactionSlice := []float64{-2, -1, 0, 0, -1, -2, -2.6, -1.6, -3}

	// // 3. limit cycle
	// interactionSlice := []float64{-0.5, -1, 0, 0, -1, -2, -2.6, -1.6, -3}

	// // 4. extinction of one species
	// interactionSlice := []float64{-2, -1, -1, -1, -1, -2, -2.6, -1.6, -3}

	// // 5. extinction of two species
	// interactionSlice := []float64{-0.1, -1, -0.1, -1, -0.1, -2, -2.6, -0.6, -3}

	// // 6. chaotic dynamics
	// interactionSlice := []float64{-1, -1.09, -1.52, 0, 0, -0.72, -0.3168, -0.9792, -3.5649, 0, -1.53, -0.7191, -1.5367, -0.6477, -0.4445, -1.27}

	interaction := SetInteractionMatrix(interactionSlice, numSpecies)

	// // method 1: randomly generate
	// deathGrowth := IniRateMatrix(numSpecies)

	// method 2: set a matrix by hand

	// 1. original paper parameters
	rateSlice := []float64{0.25, -0.5, -0.5}

	// // 2. stable equilibrium & limit cycle & extinction
	// rateSlice := []float64{3, 4, 7.2}

	// // 3. chaotic dynamics - 4 species
	// rateSlice := []float64{1, 0.72, 1.53, 1.27}

	deathGrowth := SetRateMatrix(rateSlice)

	// initialize number of generations and time interval
	numGens := 50000
	time := 0.002
	canvasWidth := 500
	frequency := 200

	fmt.Println("parameters read! Initilizing ecosystem...")

	// initialize an Ecosystem object
	initialEcosystem := InitializeEcosystem(numSpecies, pop, interaction, deathGrowth)

	fmt.Println("Ecosystem initialized! Simulating ecosystem...")

	timePoints := SimulateEcosystem(initialEcosystem, numGens, time)

	// drawing ecosystem gifs
	fmt.Println("Simulation done! Drawing the ecosystem...")

	images := DrawEcoBoards(timePoints, canvasWidth, frequency)

	fmt.Println("Images drawn!")

	fmt.Println("Generating an animated GIF.")

	gifhelper.ImagesToGIF(images, "output/original_paper")

	fmt.Println("GIF drawn!")

	// // writing data to csv file
	// fmt.Println("Writing data to csv file...")
	// WriteToCSV(timePoints, "output/original_paper.csv")
	// fmt.Println("Data written to csv file!")

}
