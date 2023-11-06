package main

import (
	"math"

	"gonum.org/v1/gonum/mat"
)

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
	for i := 1; i < numGens+1; i++ {
		timePoints[i] = UpdateEcosystem(timePoints[i-1], time)
	}

	return timePoints
}

// UpdateEcosystem() takes a pointer of Ecosystem object, a float64 object time,
// It returns a new Ecosystem object with updated population of each species.
func UpdateEcosystem(currEcosystem *Ecosystem, time float64) *Ecosystem {
	if currEcosystem == nil {
		// Handle the nil case appropriately, possibly returning nil or an error
		return nil //if that behavior is desired
	}

	// initialize a new Ecosystem object
	newEcosystem := Copy(currEcosystem)

	// get new updated popuation matrix
	newPop := UpdatePopulation(newEcosystem, time)

	// range over the species slice in the new Ecosystem object, and assign updated population to each specie
	for _, specie := range newEcosystem.species {
		// extract out the value of the updated population matrix at the index of the specie
		// row: specie.index; column: 0
		specie.population = newPop.At(specie.index, 0)
	}

	return newEcosystem
}

// Copy takes a pointer of Ecosystem object, and returns a new Ecosystem object with the same attributes.
func Copy(ecosystem *Ecosystem) *Ecosystem {
	// initialize a new Ecosystem object
	newEcosystem := &Ecosystem{}

	// copy the species slice
	newEcosystem.species = CopySpecies(ecosystem.species)

	// copy the interaction matrix
	newEcosystem.interaction = DeepCopyMatrix(ecosystem.interaction) // ecosystem.interaction

	// copy the deathGrowth matrix
	newEcosystem.deathGrowth = DeepCopyMatrix(ecosystem.deathGrowth) // ecosystem.deathGrowth

	return newEcosystem
}

func DeepCopyMatrix(m mat.Matrix) mat.Matrix {
	// Type assert to check if it's a *mat.Dense
	if md, ok := m.(*mat.Dense); ok {
		// Use the RawMatrix method to get the underlying data slice
		data := md.RawMatrix()

		// Make a new slice with the same data
		copiedData := make([]float64, len(data.Data))
		copy(copiedData, data.Data)

		// Create a new Dense matrix with the copied data
		return mat.NewDense(data.Rows, data.Cols, copiedData)
	}

	panic("Unsupported matrix type!")
}

// CopySpecies takes a slice of Specie pointers, and returns a new slice of Specie pointers with the same attributes
func CopySpecies(species []*Specie) []*Specie {
	// // Method 1 - some nil pointers problems
	// // initialize a new slice of Specie pointers
	// newSpecies := make([]*Specie, len(species))

	// // range over the species slice, and copy each specie
	// for i, specie := range newSpecies {
	// 	specie.index = species[i].index
	// 	specie.population = species[i].population
	// }

	// Method 2 - chatGPT corrected
	// Initialize a new slice of Specie pointers with the same length as the input
	newSpecies := make([]*Specie, len(species))

	// Range over the original species slice, and copy each Specie
	for i, specie := range species {
		// Make sure to initialize a new Specie to avoid nil pointer dereference
		newSpecies[i] = &Specie{
			index:      specie.index,
			population: specie.population,
		}
	}

	return newSpecies

}

// UpdatePopulation(specie, time) takes a pointer of object Species, and a float64 object time
// It returns a float64 object which is the updated population of this specie.
func UpdatePopulation(ecosystem *Ecosystem, time float64) mat.Matrix {
	// initialize a new population variable
	var newP mat.Matrix

	// calculate F = ∆t · G + 1, return a Matrix
	f := CalculateF(ecosystem.deathGrowth, time)

	// calculate H = ∆t · D, return a Matirx
	h := CalculateH(ecosystem.interaction, time)

	p := InitializePop(ecosystem.species)

	// calculate updated population
	// newPop = (h*p + f) * p
	newP = CalculatePop(f, h, p)

	return newP
}

// // CalculateF calculates the matrix F based on the deathGrowth matrix and time scalar.
// // It performs the operation F = ∆t * G + I, where I is an identity matrix of the same size as G.
func CalculateF(G mat.Matrix, deltaTime float64) mat.Matrix {
	// Get the dimensions of the deathGrowth matrix
	r, _ := G.Dims()

	// Create a new dense matrix to hold the result
	F := mat.NewDense(r, 1, nil)

	// Add the identity matrix to F
	for i := 0; i < r; i++ {
		F.Set(i, 0, G.At(i, 0)*deltaTime+1)
	}

	return F
}

// CalculateH calculates the matrix H based on the interaction matrix and time scalar.
// It performs the operation H = ∆t * D, where D is the interaction matrix.
func CalculateH(D mat.Matrix, deltaTime float64) mat.Matrix {
	// Clone the interaction matrix to avoid altering the original data
	H := mat.DenseCopyOf(D)

	// Scale the interaction matrix by deltaTime to get H
	H.Scale(deltaTime, H)

	return H
}

// CalculatePop calculates the updated population based on the given matrices.
// newPop = (h x p + f) * p
func CalculatePop(f, h, p mat.Matrix) mat.Matrix {
	// get the dimensions of the population matrix
	r, c := p.Dims()

	// initialize a new mat.matrix hp to store the result of h * p
	hp := mat.NewDense(r, c, nil)
	hp.Mul(h, p)

	// add f to the result of h * p
	hp.Add(hp, f)

	// element-wise multiplication
	newP := mat.NewDense(r, c, nil)
	newP.MulElem(hp, p)

	// apply max function to ensure no population goes below 0
	for i := 0; i < r; i++ {
		for j := 0; j < c; j++ {
			newP.Set(i, j, math.Max(0, newP.At(i, j)))
		}
	}

	return newP
}
