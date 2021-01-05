package main

import (
	// "fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	// "math"
	"os"
)

const (
	//Width and height values
	WIDTH  = 5
	HEIGHT = 5
)

func initializeMap(noiseMap *[HEIGHT][WIDTH]uint8) {
	for _, y := range noiseMap {
		for i, _ := range y {
			y[i] = 0
		}
	}
}

func writeImage(noiseMap [HEIGHT][WIDTH]uint8, image *image.RGBA) {
	for yi, yv := range noiseMap {
		for xi, xv := range yv {
			image.Set(xi, yi, color.RGBA{xv, xv, xv, 255})
		}
	}
}

func main() {

	//Creates the map
	var noiseMap [HEIGHT][WIDTH]uint8

	//Fills the map with black
	initializeMap(&noiseMap)

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	writeImage(noiseMap, img)

	//Creates an image
	f, err := os.Create("result.png")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// Encodes and saves to png
	err = png.Encode(f, img)
	if err != nil {
		log.Fatal(err)
	}
}
