package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	"gonum.org/v1/gonum/stat/distuv"
)

// InitializePopulation takes in a number of parameters such as popSize, selection Coefficient and the starting allele frequency
// It returns a population object with all of these parameters
func InitializePopulation(popSize int, selCo, freqStart float64) *Population {
	var initialPop Population
	initialPop.popSize = popSize
	initialPop.selCo = selCo
	initialPop.freqStart = freqStart
	initialPop.freqNum = freqStart * float64(popSize)
	initialPop.freq = freqStart
	initialPop.gen = 0

	return &initialPop
}

// SimulatePopulationTimePoints takes in a population object, number of generations
// it returns a slice of numGen number of pointers to populations
func SimulatePopulationTimePoints(initialPop *Population, numGen int) []*Population {
	timePoints := make([]*Population, numGen)
	timePoints[0] = initialPop
	for i := 1; i < numGen; i++ {
		timePoints[i] = SimulateOneGeneration(timePoints[i-1])
	}

	return timePoints
}

// SimulateOneGeneration takes in a population object
// It returns another population object with the frequency of the allele updated by the WF Equation
func SimulateOneGeneration(currentPop *Population) *Population {
	newPop := CopyGeneration(currentPop)
	prob := (currentPop.freqNum * (1 + currentPop.selCo)) / (currentPop.freqNum*(1+currentPop.selCo) + float64(currentPop.popSize) - currentPop.freqNum)
	var b distuv.Binomial
	b.N = float64(newPop.popSize)
	b.P = prob
	newPop.freqNum = distuv.Binomial.Rand(b)
	newPop.freq = newPop.freqNum / (float64(newPop.popSize))

	return newPop
}

// CopyGeneration takes in one Population.
// It returns another population with the same parameters but one greater gen
func CopyGeneration(currentPop *Population) *Population {
	var newPop Population
	newPop.popSize = currentPop.popSize
	newPop.gen = currentPop.gen + 1
	newPop.selCo = currentPop.selCo
	newPop.freqStart = currentPop.freqStart
	newPop.freqNum = currentPop.freqNum
	newPop.freq = currentPop.freq
	return &newPop
}

// SimulateMultipleRuns is a function that takes in an int for number of runs and the different parameters
// It returns a slice of population generation slices
func SimulateMultipleRuns(numRuns, popSize, numGen int, selCo, freqStart float64) [][]*Population {
	runs := make([][]*Population, numRuns)

	for i := 0; i < numRuns; i++ {
		iniPop := InitializePopulation(popSize, selCo, freqStart)
		timePoints := SimulatePopulationTimePoints(iniPop, numGen)
		runs = append(runs, timePoints)
	}

	return runs
}

// Function to write simulation data to CSV file
// Input the simulation results which is timePoints, and the filename(path)
// Output a csv file
func WriteToCSV(timePoints []*Population, filename string) {
	// Specify the folder name
	folderName := "WrightFisher"

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Construct the full path to the folder
	folderPath := filepath.Join(currentDir, folderName)

	// Check if the folder exists, create it if not
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}
	}

	// Construct the full path to the CSV file
	fullPath := filepath.Join(folderPath, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"Generations", "PopulationSize", "SelectionCoefficient", "StartAlleleFrequency", "NumAlleles", "AlleleFrequency"}
	writer.Write(header)

	// Write data for each generation
	for _, pop := range timePoints {
		row := []string{
			fmt.Sprint(pop.gen),
			fmt.Sprint(pop.popSize),
			fmt.Sprint(pop.selCo),
			fmt.Sprint(pop.freqStart),
			fmt.Sprint(pop.freqNum),
			fmt.Sprint(pop.freq),
		}
		writer.Write(row)
	}

	fmt.Println("CSV file created:", fullPath)
}

// Function to write parameters to CSV file
// Input all parameters, the output folder name, and filename(path)
// Output a csv file in the folder
func WriteParameters(popSize int, selCo, freqStart float64, numGen, numRuns int, folderName, filename string) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Construct the full path to the folder
	folderPath := filepath.Join(currentDir, folderName)

	// Check if the folder exists, create it if not
	if _, err := os.Stat(folderPath); os.IsNotExist(err) {
		err := os.Mkdir(folderPath, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating folder:", err)
			return
		}
	}

	// Construct the full path to the CSV file
	fullPath := filepath.Join(folderPath, filename)

	file, err := os.Create(fullPath)
	if err != nil {
		fmt.Println("Error creating CSV file:", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write CSV header
	header := []string{"PopulationSize", "SelectionCoefficient", "StartAlleleFrequency", "NumGenerations", "NumberRuns"}
	writer.Write(header)

	row := []string{
		fmt.Sprint(popSize),
		fmt.Sprint(selCo),
		fmt.Sprint(freqStart),
		fmt.Sprint(numGen),
		fmt.Sprint(numRuns),
	}
	writer.Write(row)
}
