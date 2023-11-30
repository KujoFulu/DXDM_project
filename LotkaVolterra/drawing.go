package main

import (
	"canvas"
	"image"
	"math/rand"
	"time"
)

func DrawEcoBoards(timePoints []*Ecosystem, canvasWidth int, frequency int) []image.Image {
	// create a image list slice
	imageList := make([]image.Image, 0)

	// count the number of species
	numSpecies := len(timePoints[0].species)

	// create two position slices
	xPos := make([]int, numSpecies)
	yPos := make([]int, numSpecies)

	// create a color slice
	color := make([][]uint8, numSpecies)

	// range over all species and assign them a position, store in slice
	for i := 0; i < numSpecies; i++ {
		xPos[i] = GenPosition(canvasWidth)
		yPos[i] = GenPosition(canvasWidth)
		color[i] = GenRandColor()
	}

	// range over all time points and draw them
	for i := range timePoints {
		if i%frequency == 0 {
			imageList = append(imageList, DrawToCanvas(timePoints[i], canvasWidth, xPos, yPos, color))
		}
	}
	return imageList
}

// DrawToCanvas generates the image corresponding to a canvas after drawing a Universe
// object's bodies on a square canvas that is canvasWidth pixels x canvasWidth pixels
func DrawToCanvas(timePoint *Ecosystem, canvasWidth int, xPos, yPos []int, color [][]uint8) image.Image {
	// define a scaler for the radius of each species
	rScaler := 20.0

	// set a new square canvas
	c := canvas.CreateNewCanvas(canvasWidth, canvasWidth)

	// create a black background
	c.SetFillColor(canvas.MakeColor(0, 0, 0))
	c.ClearRect(0, 0, canvasWidth, canvasWidth)
	c.Fill()

	// range over species, set each of them a position and a color, draw them as circles on the cancavs
	for _, s := range timePoint.species {
		// asign color
		red := color[s.index][0]
		green := color[s.index][1]
		blue := color[s.index][2]

		// set color, position and radius
		c.SetFillColor(canvas.MakeColor(uint8(red), uint8(green), uint8(blue)))
		centerX := float64(xPos[s.index])
		centerY := float64(yPos[s.index])
		r := s.population * rScaler

		// set r limit in order not to exceed the canvas
		if r > float64(canvasWidth/6) {
			r = float64(canvasWidth / 6)
		} else if r < 3.0 {
			r = 3.0
		}

		// draw the circle
		c.Circle(centerX, centerY, r)

		// fill the color
		c.Fill()
	}

	// we want to return an image!
	return c.GetImage()
}

func GenPosition(canvasWidth int) int {
	// define the x/y range to set a range in the middle of the canvas
	aRange := canvasWidth / 4
	bRange := canvasWidth * 3 / 4

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// set random range input
	delta := bRange - aRange

	// generate a random position between 0 to x/y range
	position := aRange + rand.Intn(delta)

	return position
}

// function GenRandColor() generates a random color RGB value
func GenRandColor() []uint8 {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// generate a random color
	red := rand.Intn(255)
	green := rand.Intn(255)
	blue := rand.Intn(255)

	// creat a color slice
	color := make([]uint8, 3)

	// add generated color to the slice
	color[0] = uint8(red)
	color[1] = uint8(green)
	color[2] = uint8(blue)

	return color
}
