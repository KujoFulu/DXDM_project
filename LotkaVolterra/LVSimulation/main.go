package main

import (
	"fmt"
	"gifhelper"
	"os"
	"strconv"
)

func main() {
	fmt.Println("Simulation of LV model starts!")

	fmt.Println("Reading input parameters...")

	// Main function's own input for original testing:

	// // initialize number of species, interaction matrix, and deathGrowth matrix
	// numSpecies := 3

	// // population matrix
	// // 1. orginal paper
	// pop := []float64{50.0, 10.0, 5.0}
	// // // 2. stable equilibrium & limit cycle & extinction
	// // pop := []float64{0.1, 0.8, 0.3}
	// // // 3. chaotic dynamics - 4 species
	// // pop := []float64{0.1, 0.8, 0.3, 0.5}

	// // interaction matrix
	// // // method 1: randomly generate
	// // interaction := InitializeInteractionMatrix(numSpecies)
	// // method 2: set a matrix by hand
	// // 1. original paper parameters
	// interactionSlice := []float64{0, -0.04, -0.04, 0.04, 0, -0.02, 0.02, 0.04, 0}
	// // // 2. stable equilibrium
	// // interactionSlice := []float64{-2, -1, 0, 0, -1, -2, -2.6, -1.6, -3}
	// // // 3. limit cycle
	// // interactionSlice := []float64{-0.5, -1, 0, 0, -1, -2, -2.6, -1.6, -3}
	// // // 4. extinction of one species
	// // interactionSlice := []float64{-2, -1, -1, -1, -1, -2, -2.6, -1.6, -3}
	// // // 5. extinction of two species
	// // interactionSlice := []float64{-0.1, -1, -0.1, -1, -0.1, -2, -2.6, -0.6, -3}
	// // // 6. chaotic dynamics
	// // interactionSlice := []float64{-1, -1.09, -1.52, 0, 0, -0.72, -0.3168, -0.9792, -3.5649, 0, -1.53, -0.7191, -1.5367, -0.6477, -0.4445, -1.27}

	// // deathGrowth matrix
	// // // method 1: randomly generate
	// // deathGrowth := IniRateMatrix(numSpecies)
	// // method 2: set a matrix by hand
	// // 1. original paper parameters
	// rateSlice := []float64{0.25, -0.5, -0.5}
	// // // 2. stable equilibrium & limit cycle & extinction
	// // rateSlice := []float64{3, 4, 7.2}
	// // // 3. chaotic dynamics - 4 species
	// // rateSlice := []float64{1, 0.72, 1.53, 1.27}

	// ************************* CLAs *************************
	// Read in CLAs from user: numSpecies, pop - slice, interaction matrix, deathGrowth matrix

	// os.Args[0] is the name of the program (./sandpile)
	fmt.Println(os.Args[0])

	// take in numSpecies CLA
	numSpecies, err1 := strconv.Atoi(os.Args[1])
	if err1 != nil {
		panic(err1)
	}
	if numSpecies <= 0 {
		panic("Error: nonpositive number given as size of the board.")
	}

	// take in pop CLA
	pop := make([]float64, numSpecies)
	for i := 0; i < numSpecies; i++ {
		pop[i], err1 = strconv.ParseFloat(os.Args[2+i], 64)
		if err1 != nil {
			panic(err1)
		}
		if pop[i] <= 0 {
			panic("Error: nonpositive number given as population.")
		}
	}

	// take in interaction matrix CLA
	interactionSlice := make([]float64, numSpecies*numSpecies)
	for i := 0; i < numSpecies*numSpecies; i++ {
		interactionSlice[i], err1 = strconv.ParseFloat(os.Args[2+numSpecies+i], 64)
		if err1 != nil {
			panic(err1)
		}
	}

	// take in deathGrowth matrix CLA
	rateSlice := make([]float64, numSpecies)
	for i := 0; i < numSpecies; i++ {
		rateSlice[i], err1 = strconv.ParseFloat(os.Args[2+numSpecies*numSpecies+i+numSpecies], 64)
		if err1 != nil {
			panic(err1)
		}
	}

	// print out all CLAs in one line
	fmt.Println("numSpecies:", numSpecies, "pop:", pop, "interactionSlice:", interactionSlice, "rateSlice:", rateSlice)

	// set interaction and deathGrowth matrix: for simulation
	transposedSlice := transposeSquareMatrix(numSpecies, interactionSlice)
	interaction := SetInteractionMatrix(transposedSlice, numSpecies)
	deathGrowth := SetRateMatrix(rateSlice)

	// initialize number of generations and time interval: for simulation
	numGens := 50000
	time := 0.002

	// initialize canvas width and frequency: for drawing
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

	gifhelper.ImagesToGIF(images, "LVSimulation/output/test")

	fmt.Println("GIF drawn!")

	// writing data to csv file
	fmt.Println("Writing data to csv file...")
	WriteToCSV(timePoints, "LVSimulation/output/test.csv")
	fmt.Println("Data written to csv file!")

}
