package main

import "gonum.org/v1/gonum/mat"

// SimulateEcosystem() takes as input the initial *Ecosystem object set by the user, a number of generations that the simulation will run,
// and a time interval at which the ecosystem will be updated.
// It returns an array of numGens + 1 *Ecosystem pointers timePoints, where timePoints[0] is the initial ecosystem,
// and timePoints[i] represents the ecosystem object in the i-th time step of the ecosystem simulation starting with initialEcosystem,
// assuming that in each step of the simulation we use a time value equal to time interval.
func SimulateEcosystem(initialEcosystem *Ecosystem, numGens int, time float64) []*Ecosystem {
	// initialize an array of numGens + 1 *Ecosystem pointers
	timePoints := make([]*Ecosystem, numGens+1)
	// assign the first element of the array to be the initial ecosystem
	timePoints[0] = initialEcosystem

	// range over the number of ecosystems and set the i-th ecosystem equal to updating the (i-1)th ecosystem
	for i := 1; i < numGens; i++ {
		timePoints[i] = UpdateEcosystem(timePoints[i-1], time)
	}

	return timePoints
}

// UpdateEcosystem() takes a pointer of Ecosystem object, a float64 object time,
// It returns a new Ecosystem object with updated population of each species.
func UpdateEcosystem(currEcosystem *Ecosystem, time float64) *Ecosystem {
	// initialize a new Ecosystem object
	newEcosystem := Copy(currEcosystem)

	// get new updated popuation matrix
	newPop := UpdatePopulation(newEcosystem, time)

	// range over the species slice in the new Ecosystem object, and assign updated population to each specie
	for _, specie := range newEcosystem.species {
		specie.population = newPop[specie.index] // indexing problem with mat.Vector datatype
	}

	return newEcosystem
}

// UpdatePopulation(specie, time) takes a pointer of object Species, and a float64 object time
// It returns a float64 object which is the updated population of this specie.
func UpdatePopulation(ecosystem *Ecosystem, time float64) mat.Vector {
	// initialize a new population variable
	var newP float64

	// calculate F = ∆t · G + 1, return a Matrix
	f := CalculateF(ecosystem.deathGrowth, time)

	// calculate H = ∆t · D, return a Matirx
	h := CalculateH(ecosystem.interaction, time)

	p := InitializePop(ecosystem.species)

	// calculate updated population
	// newPop = (h*ecosystem.allPopulation + f) * ecosystem.allPopulation
	newP = CalculatePop(f, h, p)

	return newP
}
