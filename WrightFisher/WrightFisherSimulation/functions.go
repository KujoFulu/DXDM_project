package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"math/rand"
	"time"

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

	//Copy the previous generation
	newPop := CopyGeneration(currentPop)

	//Calculate the probability of pulling the allele from the previous generation
	//This incorperates the selection coefficient
	//It then normalizes it into a frequency
	prob := (currentPop.freqNum * (1 + currentPop.selCo)) / (currentPop.freqNum*(1+currentPop.selCo) + float64(currentPop.popSize) - currentPop.freqNum)
	
	//We used the distuv packages Binomial function
	//If this does not work then we have provided another function that you can incorperate here
	var b distuv.Binomial
	b.N = float64(newPop.popSize)
	b.P = prob
	
	//This simulates taking n random pulls from the binomial distribution created by the probability
	//It returns a number of alleles
	newPop.freqNum = distuv.Binomial.Rand(b)
	
	//We also record the frequency
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
	
	//Create an array of simulations
	runs := make([][]*Population, numRuns)

	//Loops through the number of runs
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
	//folderName := "WrightFisher"

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Construct the full path to the folder
	folderPath := currentDir

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
func WriteParameters(popSize int, selCo, freqStart float64, numGen, numRuns int, filename string) {
	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Construct the full path to the folder
	folderPath := currentDir

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



//The following is an alternative to the random binomial package that needs to be installed for the main WF simulation
//Please use this if you are having difficulties in installing the package
//This should also replicate the function np.random.binomial in the numpy package in python
func Binomial(n int, p float64) int {

	b := 0

	//This loops through every n that could be our allele of interest 
	for i := 0; i < n; i++ {

		//If a random number happens to be below the probability it adds one to the output
		if rand.Float64() < p {
			b++
		}
	}

	return b
}





//Below this line are functions relating to the Two Loci Simulation

type TwoLPop struct {
	//n is the number of individuals/haplotypes in the population
	n int
	
	//haplotypeQuant contains the quantity of the four haplotypes
	//They are in the following order: 0:AB, 1:aB, 2:Ab, 3:ab
	haplotypeQuant []int

	//selCo is the selection coefficient that selects for the singleton mutation
	selCo float64

	//selProbs is the asjusted probabilities of pulling a haplotype from the population after adjusting w/ the selCo
	selProbs []float64

	//gen is the generation number
	gen int

	//singleton is the int corresponding to the mutated haplotype (same order as the )
	singleton int

	//recomb is the recombination coefficient that determines the rate of recombination
	recomb float64
}




//SimulateTwoLoci is a function that runs the entire recombination simulation
//It takes in the size of the population, the selection coefficient and the recombination coefficient set by the users
//It outputs an array of pointers to two loci populations 
//This array is the generations of the population
//The simulation runs until it is fixed or until it hits 100 generations
func SimulateTwoLoci(n int, selCo, recomb float64) []*TwoLPop {

	//Makes a new seed for the random number generator
	rand.Seed(time.Now().UnixNano())

	//This creates the initial population. It has a single copy of a mutated genotype
	initial := InitializeTwoLoci(n, selCo, recomb)

	generations := make([]*TwoLPop, 1)
	//This sets the inital population as the first in the array
	generations[0] = initial
		
	var fixed bool
	fixed = false

	i := 0

		//While the population is not fixed or less than 100
		for fixed == false || i < 100 {

			//increase the generation
			i++

			//The oldGen is the previous generation and the newGen is a copy with all haplotypes at 0
			oldGen := generations[i-1]
			newGen := CreateEmptyLociPop(oldGen)
			generations = append(generations, newGen)
			
			var freqA float64
			var freqB float64

			//freqA and freqB are the frequencies of the dominant allele at leach locus
			//ie. A/A+a and B/B+b
			freqA = float64(oldGen.haplotypeQuant[0]+oldGen.haplotypeQuant[2])/float64(oldGen.n)
			freqB = float64(oldGen.haplotypeQuant[0]+oldGen.haplotypeQuant[1])/float64(oldGen.n)

			//For each individual in the population it assigns it a hapltype
			for j := 0; j < n; j++ {

				//It selects a haplotype based on the frequencies in the previous generation
				//as well as the selection coefficient selecting for the mutation
				var newHaplo int
				newHaplo = RandomSelection(oldGen)

				//After the Haplotype has been chosen then, based on the recombination coefficient, be recombined
				recombBool := ChooseRecombination(oldGen)

				if recombBool == false {
					//no recombination means that haplotype just increases in quantity
					newGen.haplotypeQuant[newHaplo]++

				} else {

					//if it is recombined then it will be randomly recombined based on the allele quantities in the previous generation
					var recombHaplo int
					recombHaplo = RecombineGenotype(newHaplo, freqA, freqB)

					//Adding to the quantity of the newly recombined haplotype
					newGen.haplotypeQuant[recombHaplo]++
				}

			}
			
			//Checking to see if it is fixed
			fixed = IsTLFixed(newGen)
		}

	return generations
}



//InitializeTwoLoci is a function that takes in a population size, a selection coefficient and a recombintion coefficient
//It returns a pointer to a two loci population
//It has a single  mutation at the second locus (either b or B)
func InitializeTwoLoci(n int, selCo, recomb float64) *TwoLPop {

	//The frequeny of the dominant allele (A) at the A locus is determined by running the WF single site simulation
	//It is run 100 times for 100 gens, and then the quantity is pulled at random from the the distribution of frequencies 
	quantA := AlleleDistribution(n)

	var pop TwoLPop

	//copy over the selection coefficient
	pop.selCo = selCo
	
	//This is the first generation
	pop.gen = 0

	//copy over the recombination coefficient
	pop.recomb = recomb

	//Randomly choose a singleton
	s := rand.Intn(4)
	pop.singleton = rand.Intn(4)

	if s == 0 {

		//This means the singleton is AB and there is only one copy of B
		pop.haplotypeQuant[0] = 1 //AB
		pop.haplotypeQuant[1] = 0 //aB
		pop.haplotypeQuant[2] = quantA-1 //Ab (the 1 is for the singleton)
		pop.haplotypeQuant[3] = n - quantA //ab

	} else if s == 1 {

		//This means the singleton is aB and there is only one copy of B
		pop.haplotypeQuant[0] = 0
		pop.haplotypeQuant[1] = 1
		pop.haplotypeQuant[2] = quantA
		pop.haplotypeQuant[3] = n - quantA -1

	} else if s == 2 {

		//This means the singleton is Ab and there is only one copy of b
		pop.haplotypeQuant[0] = quantA-1
		pop.haplotypeQuant[1] = n - quantA
		pop.haplotypeQuant[2] = 1
		pop.haplotypeQuant[3] = 0
		
	} else {

		//This means the singleton is ab and there is only one copy of b
		pop.haplotypeQuant[0] = quantA
		pop.haplotypeQuant[1] = n - quantA-1
		pop.haplotypeQuant[2] = 1
		pop.haplotypeQuant[3] = 0

	}

	//The selection probabilities are calculated from the quanties of the genotypes and if it contains the mutant at the b locus
	//if it contains the mutant the 
	selectionProbs := SelectionProbabilities(&pop)
	
	//Copy over the selection probabilities into the population
	pop.selProbs[0] = selectionProbs[0]
	pop.selProbs[1] = selectionProbs[1]
	pop.selProbs[2] = selectionProbs[2]
	pop.selProbs[3] = selectionProbs[3]


	//return the pointer to the population
	return &pop
}


//AlleleDistribution is a function that take in a population number
//It returns a quantity of an allele taken from the distribution of frequencies after running the WF 100 times
func AlleleDistribution(n int) int {

	//Run the Wright Fisher simulation for a single site 100 times
	//This will have the same population size as the two loci, run for 100 generations, have a selection coefficient of 0 and a starting frequency of 0.5
	simulation := SimulateMultipleRuns(100, n, 100, 0.0, 0.5)

	//Randonly choose a run to pull from
	randRun := rand.Intn(1000)

	//Randomly choose a generation to pull from
	randGen := rand.Intn(100)

	//Pull quantity of that allele
	freqNum := int(simulation[randRun][randGen].freqNum)
	return freqNum
}


//CreateEmptyLociPop is a function that takes in a pointer to a two loci population
//It then copies that population but all of the haplotype quantities are set to 0
//Additionally it calculates the selection probabilites based on the frequencies of the previous generation
func CreateEmptyLociPop(oldGen *TwoLPop) *TwoLPop {
	
	var newGen TwoLPop

	newGen.n = oldGen.n

	hapQuant := make([]int, 4)
	newGen.haplotypeQuant = hapQuant

	newGen.selCo = oldGen.selCo 

	s := SelectionProbabilities(oldGen)
	newGen.selProbs[0] = s[0]
	newGen.selProbs[1] = s[1]
	newGen.selProbs[2] = s[2]
	newGen.selProbs[3] = s[3]

	//This is the next generation
	newGen.gen = oldGen.gen + 1

	newGen.singleton = oldGen.singleton

	newGen.recomb = oldGen.recomb

	return &newGen
}


//SelectionProbabilities is a function that takes in a pointer to a two loci population
//It returns an array of floats representing the different selection probabilties of the four haplotypes
//The selection probabilities are based on the frequencies in the population that is fed into it
func SelectionProbabilities(oldGen *TwoLPop) []float64 {

	//First calculate the frequencies of each haplotype based on their quantities
	freqAB := float64(oldGen.haplotypeQuant[0])/float64(oldGen.n)
	freqaB := float64(oldGen.haplotypeQuant[1])/float64(oldGen.n)
	freqAb := float64(oldGen.haplotypeQuant[2])/float64(oldGen.n)
	freqab := float64(oldGen.haplotypeQuant[3])/float64(oldGen.n)
	
	//If the singleton mutation is B
	if oldGen.singleton == 0 || oldGen.singleton == 1 {
		
		//All of the haplotypes with B have the selection coefficient added to them
		freqAB += oldGen.selCo
		freqaB += oldGen.selCo
		//All of the haplotypes with b have the selection coefficient reduced from them
		freqAb -= oldGen.selCo
		freqab -= oldGen.selCo

	} else {
		//if the singleton muation is b

		//All of the haplotypes with B have the selection coefficient reduced from them
		freqAB -= oldGen.selCo
		freqaB -= oldGen.selCo
		//All of the haplotypes with b have the selection coefficient added to them
		freqAb += oldGen.selCo
		freqab += oldGen.selCo
	}

	//normalize the frequencies after adjusting with the selection coefficient
	sum := freqAB + freqaB + freqAb + freqab

	freqAB = freqAB/sum
	freqaB = freqaB/sum
	freqAb = freqAb/sum
	freqab = freqab/sum

	//Create array and return the new frequencies after they have been adjusted
	selectProbs := make([]float64, 4)

	selectProbs[0] = freqAB
	selectProbs[1] = freqaB
	selectProbs[2] = freqAb
	selectProbs[3] = freqab

	return selectProbs

}


//RandomSelection takes in a pointer to a generation
//It returns a haplotype based on the selection probabilities in that population
func RandomSelection(oldGen *TwoLPop) int {

	pull := rand.Float64()

	if pull <= oldGen.selProbs[0] {

		return 0
	
	} else if pull <= oldGen.selProbs[0] + oldGen.selProbs[1] {

		return 1

	} else if pull <= oldGen.selProbs[0] + oldGen.selProbs[1] + oldGen.selProbs[2] {

		return 2

	} else {

		return 3

	}

}


//ChooseRecombination is a function that takes in a pointer to a two loci population
//It returns a boolean of if a gene should be recombined or not
//This is based off of the recombination coefficient
func ChooseRecombination(oldGen *TwoLPop) bool {

	pull := rand.Float64()

	if pull >= oldGen.recomb {
		return true
	} else {
		return false
	}
	
}


//RecombineGenotype takes in a Haplotype integer and the frequencies of the dominant alleles of the two loci (A and B)
//It returns an integer corresponding to the recombined haplotype
//It randomly chooses the loci to recombine
//It then chooses the new allele at that loci based on the allele frequencies
func RecombineGenotype(newHaplo int, freqA, freqB float64) int {

	//Essentially a coin flip to determine which locus recombines
	pull := rand.Intn(2)

	if pull == 0 {

		//This keeps the A locus and is recombining the B locus

		pullB := rand.Float64()

		if pullB <= freqB {
			//The new allele at the B locus is B

			if newHaplo == 0 || newHaplo == 2 {
				//If the A locus is A
				return 0
				//AB

			} else {
				//If the A locus is a
				return 1
				//aB
			}

		} else {
			//The new allele at the B locus is b

			if newHaplo == 0 || newHaplo == 2 {
				//If the A locus is A
				return 2
				//Ab

			} else {
				//If the A locus is a
				return 3
				//ab
			}

		}

	} else {
		//The locus getting recombined is the A locus

		pullA := rand.Float64()

		if pullA <= freqA {
			//The new allele at the A locus is A

			if newHaplo == 0 || newHaplo == 1 {
				//If the B locus is B
				return 0
				//AB

			} else {
				//If the B locus is b
				return 2
				//Ab
			}

		} else {
			//The new allele at the A locus is a

			if newHaplo == 0 || newHaplo == 1 {
				//If the B locus is B
				return 1
				//aB

			} else {
				//If the B locus is b
				return 3
				//ab
			}

		}
	}

}


//IsTLFixed is a function that takes in a pointer to a two loci population
//It returns a boolean of if that population is fixed or not
func IsTLFixed(gen *TwoLPop) bool {
	
	var fixed bool
	fixed = false

	var freqA float64
	var freqB float64

	//freqA and freqB are the frequencies of the dominant allele at leach locus
	//ie. A/A+a and B/B+b
	freqA = float64(gen.haplotypeQuant[0]+gen.haplotypeQuant[2])/float64(gen.n)
	freqB = float64(gen.haplotypeQuant[0]+gen.haplotypeQuant[1])/float64(gen.n)

	//Checking to see if any of the alleles have been fixed or lost
	if freqA == 0.0 || freqA == 1.0 {
		fixed = true
	} else if freqB == 0.0 || freqB == 1.0 {
		fixed = true
	}
	
	return fixed
}


// Function to write two loci simulation data to CSV file
// Input the simulation results which is timePoints, and the filename(path)
// Output a csv file
func WritetwoLToCSV(timePoints []*TwoLPop, filename string) {
	// Specify the folder name
	//folderName := "WrightFisher"

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current working directory:", err)
		return
	}

	// Construct the full path to the folder
	folderPath := currentDir

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
	header := []string{"Generations", "PopulationSize", "SelectionCoefficient", "RecombinationCoefficient", "NumSingleton", "SingletonFrequency"}
	writer.Write(header)

	// Write data for each generation
	for _, pop := range timePoints {
		row := []string{
			fmt.Sprint(pop.gen),
			fmt.Sprint(pop.n),
			fmt.Sprint(pop.selCo),
			fmt.Sprint(pop.recomb),
			fmt.Sprint(pop.haplotypeQuant[pop.singleton]),
			fmt.Sprint(pop.haplotypeQuant[pop.singleton]/pop.n),
		}
		writer.Write(row)
	}

	fmt.Println("CSV file created:", fullPath)
}