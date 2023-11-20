package main

import (

	"gonum.org/v1/gonum/stat/distuv"
)

//InitializePopulation takes in a number of parameters such as popSize, selection Coefficient and the starting allele frequency
//It returns a population object with all of these parameters
func InitializePopulation(popSize int, selCo, freqStart float64) *Population {
	var initialPop Population
	initialPop.popSize = popSize
	initialPop.selCo = selCo
	initialPop.freqStart = freqStart
	initialPop.freqNum = freqStart*float64(popSize)
	initialPop.freq = freqStart
	initialPop.gen = 0
	
	return &initialPop
}

//SimulatePopulationTimePoints takes in a population object, number of generations
//it returns a slice of numGen number of pointers to populations
func SimulatePopulationTimePoints(initialPop *Population, numGen int) []*Population {
	timePoints := make([]*Population, numGen)
	timePoints[0] = initialPop
	for i := 1; i < numGen; i++ {
		timePoints[i] = SimulateOneGeneration(timePoints[i-1])
	}

	return timePoints
}

//SimulateOneGeneration takes in a population object
//It returns another population object with the frequency of the allele updated by the WF Equation
func SimulateOneGeneration(currentPop *Population) *Population {
	newPop := CopyGeneration(currentPop)
	prob := (currentPop.freqNum * (1+currentPop.selCo))/(currentPop.freqNum * (1+currentPop.selCo) + float64(currentPop.popSize) - currentPop.freqNum)
	var b distuv.Binomial
	b.N = float64(newPop.popSize)
	b.P = prob
	newPop.freqNum = distuv.Binomial.Rand(b)
	newPop.freq = newPop.freqNum/(float64(newPop.popSize))

	return newPop
}

//CopyGeneration takes in one Population.
//It returns another population with the same parameters but one greater gen
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

//SimulateMultipleRuns is a function that takes in an int and the different parameters
//It returns a slice of population generation slices
//func SimulateMultipleRuns(numRuns, popSize int, selCo, freqStart float64) *[][]Population {

//}