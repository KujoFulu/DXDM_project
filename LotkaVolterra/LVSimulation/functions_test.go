package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strconv"
	"strings"
	"testing"

	"gonum.org/v1/gonum/mat"
)

// CalculatePopTest struct holds information for a test case
type CalculatePopTest struct {
	fInput, hInput, pInput mat.Matrix
	expectedOutput         *mat.Dense // Changed to *mat.Dense
}

// TestCalculatePop tests the CalculatePop() function
func TestCalculatePop(t *testing.T) {
	tests := ReadCalculatePopTests("Tests/CalculatePop/")

	for _, test := range tests {
		// Run the function and get the result matrix
		resultPtr := CalculatePop(test.fInput, test.hInput, test.pInput)
		result := resultPtr // Dereference the pointer to get the matrix

		// Check the result
		if !mat.EqualApprox(result, test.expectedOutput, 0.001) {
			t.Errorf("CalculatePop() = %v, want %v", result, test.expectedOutput)
		}
	}
}

// ReadCalculatePopTests reads test cases for CalculatePop
func ReadCalculatePopTests(directory string) []CalculatePopTest {
	inputFiles := ReadDirectory(directory + "/input")
	outputFiles := ReadDirectory(directory + "/output")

	if len(inputFiles) != len(outputFiles) {
		panic("Error: number of input and output files do not match!")
	}

	tests := make([]CalculatePopTest, len(inputFiles))
	for i := range inputFiles {
		inputMatrices := ReadMatrices(directory + "/input/" + inputFiles[i].Name())
		expectedOutput := ReadMatrix(directory + "/output/" + outputFiles[i].Name())

		tests[i] = CalculatePopTest{
			fInput:         inputMatrices[0],
			hInput:         inputMatrices[1],
			pInput:         inputMatrices[2],
			expectedOutput: expectedOutput.(*mat.Dense), // Type assertion to *mat.Dense
		}
	}

	return tests
}

// ReadMatrices reads multiple matrices from a single file
func ReadMatrices(file string) []mat.Matrix {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var matrices []mat.Matrix
	var data []float64
	var nrows, ncols int

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" { // Assuming blank line as delimiter
			matrices = append(matrices, mat.NewDense(nrows, ncols, data))
			data = []float64{}
			nrows = 0
			continue
		}

		nrows++
		rowData := strings.Split(line, " ")
		if ncols == 0 {
			ncols = len(rowData)
		}
		for _, strVal := range rowData {
			val, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				panic(fmt.Sprintf("Error parsing matrix value: %s", err))
			}
			data = append(data, val)
		}
	}
	matrices = append(matrices, mat.NewDense(nrows, ncols, data)) // Append last matrix

	return matrices
}

// ReadMatrix reads in a mat.Matrix from a file
func ReadMatrix(file string) mat.Matrix {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var data []float64
	var nrows int
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		nrows++
		rowStr := scanner.Text()
		rowData := strings.Split(rowStr, " ")
		for _, strVal := range rowData {
			val, err := strconv.ParseFloat(strVal, 64)
			if err != nil {
				panic(fmt.Sprintf("Error parsing matrix value: %s", err))
			}
			data = append(data, val)
		}
	}

	// Assuming the matrix is square for simplicity
	ncols := len(data) / nrows
	return mat.NewDense(nrows, ncols, data)
}

// ReadDirectory() reads in a directory and returns a slice of fs.DirEntry objects containing file info for the directory
func ReadDirectory(dir string) []fs.DirEntry {
	// read in all the files in the given directory
	files, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	return files
}

type CalculateHTest struct {
	dInput         mat.Matrix
	deltaTime      float64
	expectedOutput mat.Matrix
}

// TestCalculateH tests the CalculateH() function
func TestCalculateH(t *testing.T) {
	tests := ReadCalculateHTests("Tests/CalculateH/")

	for _, test := range tests {
		// Run the function
		result := CalculateH(test.dInput, test.deltaTime)

		// Check the result
		if !mat.EqualApprox(result, test.expectedOutput, 0.001) {
			t.Errorf("CalculateH() = %v, want %v", result, test.expectedOutput)
		}
	}
}

// ReadCalculateHTests reads test cases for CalculateH
func ReadCalculateHTests(directory string) []CalculateHTest {
	inputFiles := ReadDirectory(directory + "/input")
	outputFiles := ReadDirectory(directory + "/output")

	if len(inputFiles) != len(outputFiles) {
		panic("Mismatch in number of input and output files")
	}

	tests := make([]CalculateHTest, len(inputFiles))
	for i := range inputFiles {
		tests[i].dInput, tests[i].deltaTime = ReadMatrixAndScalar(directory + "/input/" + inputFiles[i].Name())
		tests[i].expectedOutput = ReadMatrix(directory + "/output/" + outputFiles[i].Name())
	}

	return tests
}

// ReadMatrixAndScalar reads a matrix and a scalar value from a file
func ReadMatrixAndScalar(file string) (mat.Matrix, float64) {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var data []float64
	var nrows, ncols int
	var deltaTime float64
	readingMatrix := true

	for scanner.Scan() {
		line := scanner.Text()

		// Check for blank line indicating end of matrix data
		if line == "" {
			readingMatrix = false
			continue
		}

		if readingMatrix {
			// Reading matrix data
			rowData := strings.Split(line, " ")
			if ncols == 0 {
				ncols = len(rowData)
			}
			for _, strVal := range rowData {
				val, err := strconv.ParseFloat(strVal, 64)
				if err != nil {
					panic(fmt.Sprintf("Error parsing matrix value: %s", err))
				}
				data = append(data, val)
			}
			nrows++
		} else {
			// Reading the scalar value
			deltaTime, err = strconv.ParseFloat(line, 64)
			if err != nil {
				panic(fmt.Sprintf("Error parsing deltaTime value: %s", err))
			}
			break // Assuming only one scalar value after the matrix
		}
	}

	matrix := mat.NewDense(nrows, ncols, data)
	return matrix, deltaTime
}

type CalculateFTest struct {
	gInput         mat.Matrix
	deltaTime      float64
	expectedOutput mat.Matrix
}

// TestCalculateF tests the CalculateF() function
func TestCalculateF(t *testing.T) {
	tests := ReadCalculateFTests("Tests/CalculateF/")

	for _, test := range tests {
		// Run the function
		result := CalculateF(test.gInput, test.deltaTime)

		// Check the result
		if !mat.EqualApprox(result, test.expectedOutput, 0.001) {
			t.Errorf("CalculateF() = %v, want %v", result, test.expectedOutput)
		}
	}
}

// ReadCalculateFTests reads test cases for CalculateF
func ReadCalculateFTests(directory string) []CalculateFTest {
	inputFiles := ReadDirectory(directory + "/input")
	outputFiles := ReadDirectory(directory + "/output")

	if len(inputFiles) != len(outputFiles) {
		panic("Mismatch in number of input and output files")
	}

	tests := make([]CalculateFTest, len(inputFiles))
	for i := range inputFiles {
		tests[i].gInput, tests[i].deltaTime = ReadMatrixAndScalar(directory + "/input/" + inputFiles[i].Name())
		tests[i].expectedOutput = ReadMatrix(directory + "/output/" + outputFiles[i].Name())
	}

	return tests
}

type SetInteractionMatrixTest struct {
	interactionSlice []float64
	numSpecies       int
	expectedOutput   mat.Matrix
}

// TestSetInteractionMatrix tests the SetInteractionMatrix() function
func TestSetInteractionMatrix(t *testing.T) {
	tests := ReadSetInteractionMatrixTests("Tests/SetInteractionMatrix/")

	for _, test := range tests {
		// Run the function
		result := SetInteractionMatrix(test.interactionSlice, test.numSpecies)

		// Check the result
		if !mat.Equal(result, test.expectedOutput) {
			t.Errorf("SetInteractionMatrix() = %v, want %v", result, test.expectedOutput)
		}
	}
}

// ReadSetInteractionMatrixTests reads test cases for SetInteractionMatrix
func ReadSetInteractionMatrixTests(directory string) []SetInteractionMatrixTest {
	inputFiles := ReadDirectory(directory + "/input")
	outputFiles := ReadDirectory(directory + "/output")

	if len(inputFiles) != len(outputFiles) {
		panic("Mismatch in number of input and output files")
	}

	tests := make([]SetInteractionMatrixTest, len(inputFiles))
	for i := range inputFiles {
		tests[i].interactionSlice, tests[i].numSpecies = ReadInteractionSliceAndNumSpecies(directory + "/input/" + inputFiles[i].Name())
		tests[i].expectedOutput = ReadMatrix(directory + "/output/" + outputFiles[i].Name())
	}

	return tests
}

// ReadInteractionSliceAndNumSpecies reads the interaction slice and number of species from a file
func ReadInteractionSliceAndNumSpecies(file string) ([]float64, int) {
	f, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %s", err))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var interactionSlice []float64
	var numSpecies int
	readingSlice := true

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line is a delimiter (e.g., blank line) indicating the end of the interaction slice
		if line == "" {
			readingSlice = false
			continue
		}

		if readingSlice {
			// Read interaction slice values
			parts := strings.Split(line, " ")
			for _, p := range parts {
				val, err := strconv.ParseFloat(p, 64)
				if err != nil {
					panic(fmt.Sprintf("Error parsing float64 value in interaction slice: %s", err))
				}
				interactionSlice = append(interactionSlice, val)
			}
		} else {
			// Read number of species
			numSpecies, err = strconv.Atoi(line)
			if err != nil {
				panic(fmt.Sprintf("Error parsing number of species: %s", err))
			}
			break // Assuming number of species is the first line after the interaction slice
		}
	}

	return interactionSlice, numSpecies
}

type SetRateMatrixTest struct {
	rateSlice      []float64
	expectedOutput mat.Matrix
}

// TestSetRateMatrix tests the SetRateMatrix() function
func TestSetRateMatrix(t *testing.T) {
	tests := ReadSetRateMatrixTests("Tests/SetRateMatrix/")

	for _, test := range tests {
		// Run the function
		result := SetRateMatrix(test.rateSlice)

		// Check the result
		if !mat.Equal(result, test.expectedOutput) {
			t.Errorf("SetRateMatrix() = %v, want %v", result, test.expectedOutput)
		}
	}
}

// ReadSetRateMatrixTests reads test cases for SetRateMatrix
func ReadSetRateMatrixTests(directory string) []SetRateMatrixTest {
	inputFiles := ReadDirectory(directory + "/input")
	outputFiles := ReadDirectory(directory + "/output")

	if len(inputFiles) != len(outputFiles) {
		panic("Mismatch in number of input and output files")
	}

	tests := make([]SetRateMatrixTest, len(inputFiles))
	for i := range inputFiles {
		tests[i].rateSlice = ReadRateSlice(directory + "/input/" + inputFiles[i].Name())
		tests[i].expectedOutput = ReadMatrix(directory + "/output/" + outputFiles[i].Name())
	}

	return tests
}

// ReadRateSlice reads the rate slice from a file
func ReadRateSlice(file string) []float64 {
	f, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %s", err))
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var rateSlice []float64

	for scanner.Scan() {
		val, err := strconv.ParseFloat(scanner.Text(), 64)
		if err != nil {
			panic(fmt.Sprintf("Error parsing float64 value in rate slice: %s", err))
		}
		rateSlice = append(rateSlice, val)
	}

	return rateSlice
}

type InitializePopTest struct {
	species        []*Specie
	expectedOutput mat.Matrix
}

// TestInitializePop tests the InitializePop() function
func TestInitializePop(t *testing.T) {
	tests := ReadInitializePopTests("Tests/InitializePop/")

	for _, test := range tests {
		// Run the function
		result := InitializePop(test.species)

		// Check the result
		if !mat.Equal(result, test.expectedOutput) {
			t.Errorf("InitializePop() = %v, want %v", result, test.expectedOutput)
		}
	}
}

// ReadInitializePopTests reads test cases for InitializePop
func ReadInitializePopTests(directory string) []InitializePopTest {
	inputFiles := ReadDirectory(directory + "/input")
	outputFiles := ReadDirectory(directory + "/output")

	if len(inputFiles) != len(outputFiles) {
		panic("Mismatch in number of input and output files")
	}

	tests := make([]InitializePopTest, len(inputFiles))
	for i := range inputFiles {
		tests[i].species = ReadSpeciesArray(directory + "/input/" + inputFiles[i].Name())
		tests[i].expectedOutput = ReadMatrix(directory + "/output/" + outputFiles[i].Name())
	}

	return tests
}

// ReadSpeciesArray reads an array of Specie objects from a file
func ReadSpeciesArray(file string) []*Specie {
	f, err := os.Open(file)
	if err != nil {
		panic(fmt.Sprintf("Error opening file: %s", err))
	}
	defer f.Close()

	var speciesArray []*Specie
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		var s Specie
		line := scanner.Text()
		parts := strings.Split(line, " ")

		if len(parts) != 2 {
			panic("Incorrect format in species data")
		}

		s.index, err = strconv.Atoi(parts[0])
		if err != nil {
			panic(fmt.Sprintf("Error parsing species index: %s", err))
		}

		s.population, err = strconv.ParseFloat(parts[1], 64)
		if err != nil {
			panic(fmt.Sprintf("Error parsing species population: %s", err))
		}

		speciesArray = append(speciesArray, &s)
	}

	return speciesArray
}
