package main

import (
	"math"
	"gonum.org/v1/gonum/stat/distuv"
)

//InitializePopulation takes in a number of parameters such as popSize, selection Coefficient and the starting allele frequency
//It returns a population object with all of these parameters
func InitializePopulation(popSize int, selCo, freqStart float64) *Population {
	var initialPop Population
	initialPop.popSize = popSize
	initialPop.selCo = selCo
	initialPop.freqStart = freqStart
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