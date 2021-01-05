package main

import (
	// "fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"time"
)

const (
	//Width and height values
	WIDTH      = 500
	HEIGHT     = 500
	CELL_COUNT = 200
	BOARD_SIZE = 10.0
)

func initialize_map(noise_map *[HEIGHT][WIDTH]uint8) {
	for _, y := range noise_map {
		for i, _ := range y {
			y[i] = 1
		}
	}
}

func abs(num float32) float32 {
	var zero float32 = 0.0
	if num < zero {
		return -num
	} else {
		return num
	}
}

func distance(p1, p2 [2]float32) float32 {
	delta_x := abs(p1[0] - p2[0])
	delta_y := abs(p1[1] - p2[1])
	return (delta_x * delta_x) + (delta_y * delta_y)
}

func clamp(v, lo, hi float32) float32 {
	if v < lo {
		return lo
	} else if v > hi {
		return hi
	} else {
		return v
	}
}

func min(d1, d2 float32) float32 {
	if d1 < d2 {
		return d1
	} else {
		return d2
	}
}

func cellular_noise(noise_map *[HEIGHT][WIDTH]uint8) {
	cells := make([][2]float32, CELL_COUNT)
	for i := range cells {
		t := time.Now()
		nsec := t.UnixNano()
		s := rand.NewSource(int64(nsec))
		r := rand.New(s)
		cells[i] = [2]float32{r.Float32() * BOARD_SIZE, r.Float32() * BOARD_SIZE}
	}
	for yi, yv := range noise_map {
		for xi, _ := range yv {
			current_point := [2]float32{float32(xi) / float32(WIDTH) * BOARD_SIZE, float32(yi) / float32(HEIGHT) * BOARD_SIZE}
			var minimum_distance float32 = math.MaxFloat32
			for _, v := range cells {
				distance := distance(current_point, v)
				minimum_distance = min(distance, minimum_distance)
			}
			noise_map[yi][xi] = uint8(clamp(minimum_distance, 0.0, 1.0) * 255)
		}
	}
}

func write_image(noise_map [HEIGHT][WIDTH]uint8, image *image.RGBA) {
	for yi, yv := range noise_map {
		for xi, xv := range yv {
			image.Set(xi, yi, color.RGBA{xv, xv, xv, 255})
		}
	}
}

func main() {

	//Creates the map
	var noise_map [HEIGHT][WIDTH]uint8

	//Fills the map with white
	initialize_map(&noise_map)

	cellular_noise(&noise_map)

	img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	write_image(noise_map, img)

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
