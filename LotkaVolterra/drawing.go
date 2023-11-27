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

	// range over all species and assign them a position, store in slice
	for i := 0; i < numSpecies; i++ {
		xPos[i] = GenPosition(canvasWidth)
		yPos[i] = GenPosition(canvasWidth)
	}

	// range over all time points and draw them
	for i := range timePoints {
		if i%frequency == 0 {
			imageList = append(imageList, DrawToCanvas(timePoints[i], canvasWidth, xPos, yPos))
		}
	}
	return imageList
}

// DrawToCanvas generates the image corresponding to a canvas after drawing a Universe
// object's bodies on a square canvas that is canvasWidth pixels x canvasWidth pixels
func DrawToCanvas(timePoint *Ecosystem, canvasWidth int, xPos, yPos []int) image.Image {
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
		// set color, position and radius
		c.SetFillColor(canvas.MakeColor(255, 192, 203))
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
