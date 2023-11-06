package main

import (
	"gonum.org/v1/gonum/mat"
)

type Ecosystem struct {
	species     []*Specie
	interaction mat.Matrix
	deathGrowth mat.Matrix
}

type Specie struct {
	population float64
	index      int
}
