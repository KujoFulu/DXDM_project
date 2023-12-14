package main

import (
	"testing"
)

// Test function SimulateMultipleRuns
func TestSimulateMultipleRuns(t *testing.T) {
	// Test parameters
	popSize := 100
	selCo := 0.02
	freqStart := 0.5
	numGen := 10
	numRuns := 5

	// Run simulations
	simulationRuns := SimulateMultipleRuns(numRuns, popSize, numGen, selCo, freqStart)
	count := 0
	// Check the consistency of results
	for runIndex, timePoints := range simulationRuns {
		count += 1
		for genIndex, pop := range timePoints {
			// Common assertions for all runs and generations
			if pop.popSize != popSize {
				t.Errorf("Run %d, Generation %d: Incorrect population size. Got %d, want %d", runIndex, genIndex, pop.popSize, popSize)
			}

			// Example: Check if allele frequency is within a reasonable range
			if pop.freq < 0 || pop.freq > 1 {
				t.Errorf("Run %d, Generation %d: Allele frequency out of valid range. Got %f, want [0, 1]", runIndex, genIndex, pop.freq)
			}

		}
	}
	if count != numGen {
		t.Errorf("Incorrect number of Gnerations. Got %d, want %d", count, numGen)
	}

}

// Test function CopyGeneration
func TestCopyGeneration(t *testing.T) {
	// Test parameters
	popSize := 100
	selCo := 0.02
	freqStart := 0.5

	// Create an initial population
	initialPop := InitializePopulation(popSize, selCo, freqStart)

	// Copy the initial population
	copiedPop := CopyGeneration(initialPop)

	// Check if the copied population has the same parameters
	if copiedPop.popSize != initialPop.popSize {
		t.Errorf("Incorrect population size. Got %d, want %d", copiedPop.popSize, initialPop.popSize)
	}

	if copiedPop.gen != initialPop.gen+1 {
		t.Errorf("Incorrect generation. Got %d, want %d", copiedPop.gen, initialPop.gen+1)
	}

	if copiedPop.selCo != initialPop.selCo {
		t.Errorf("Incorrect selection coefficient. Got %f, want %f", copiedPop.selCo, initialPop.selCo)
	}

	if copiedPop.freqStart != initialPop.freqStart {
		t.Errorf("Incorrect starting allele frequency. Got %f, want %f", copiedPop.freqStart, initialPop.freqStart)
	}

	if copiedPop.freqNum != initialPop.freqNum {
		t.Errorf("Incorrect allele frequency count. Got %f, want %f", copiedPop.freqNum, initialPop.freqNum)
	}

	if copiedPop.freq != initialPop.freq {
		t.Errorf("Incorrect allele frequency. Got %f, want %f", copiedPop.freq, initialPop.freq)
	}
}

// Test function SimulatePopulationTimePoints
func TestSimulatePopulationTimePoints(t *testing.T) {
	// Test parameters
	popSize := 100
	selCo := 0.02
	freqStart := 0.5
	numGen := 10

	// Create an initial population
	initialPop := InitializePopulation(popSize, selCo, freqStart)

	// Simulate population time points
	timePoints := SimulatePopulationTimePoints(initialPop, numGen)

	// Check the length of the timePoints slice
	if len(timePoints) != numGen {
		t.Errorf("Incorrect number of generations. Got %d, want %d", len(timePoints), numGen)
	}

	// Check the consistency of results for each generation
	for genIndex, pop := range timePoints {
		// Example assertions (modify based on your simulation logic)
		if pop.popSize != popSize {
			t.Errorf("Generation %d: Incorrect population size. Got %d, want %d", genIndex, pop.popSize, popSize)
		}

		// Example: Check if allele frequency is within a reasonable range
		if pop.freq < 0 || pop.freq > 1 {
			t.Errorf("Generation %d: Allele frequency out of valid range. Got %f, want [0, 1]", genIndex, pop.freq)
		}
	}

}

// Test function InitializePopulation
func TestInitializePopulation(t *testing.T) {
	// Test parameters
	popSize := 100
	selCo := 0.02
	freqStart := 0.5

	// Initialize population
	initialPop := InitializePopulation(popSize, selCo, freqStart)

	// Check the initialized population parameters
	if initialPop.popSize != popSize {
		t.Errorf("Incorrect population size. Got %d, want %d", initialPop.popSize, popSize)
	}

	if initialPop.selCo != selCo {
		t.Errorf("Incorrect selection coefficient. Got %f, want %f", initialPop.selCo, selCo)
	}

	if initialPop.freqStart != freqStart {
		t.Errorf("Incorrect starting allele frequency. Got %f, want %f", initialPop.freqStart, freqStart)
	}

	if initialPop.freqNum != freqStart*float64(popSize) {
		t.Errorf("Incorrect allele frequency count. Got %f, want %f", initialPop.freqNum, freqStart*float64(popSize))
	}

	if initialPop.freq != freqStart {
		t.Errorf("Incorrect allele frequency. Got %f, want %f", initialPop.freq, freqStart)
	}

	if initialPop.gen != 0 {
		t.Errorf("Incorrect generation. Got %d, want %d", initialPop.gen, 0)
	}
}
