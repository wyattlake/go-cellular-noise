package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"sort"
	"time"
)

const (
	//Width and height values
	WIDTH           = 500
	HEIGHT          = 500
	CELL_COUNT      = 50
	BOARD_SIZE      = 5.0
	MOVEMENT_FACTOR = 1.0
)

func initialize_map(noise_map *[HEIGHT][WIDTH]uint8) {
	for _, y := range noise_map {
		for i, _ := range y {
			y[i] = 1
		}
	}
}

func abs(num float64) float64 {
	var zero float64 = 0.0
	if num < zero {
		return -num
	} else {
		return num
	}
}

func distance(p1, p2 [2]float64) float64 {
	delta_x := abs(p1[0] - p2[0])
	delta_y := abs(p1[1] - p2[1])
	return (delta_x * delta_x) + (delta_y * delta_y)
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	} else if v > hi {
		return hi
	} else {
		return v
	}
}

func min(d1, d2 float64) float64 {
	if d1 < d2 {
		return d1
	} else {
		return d2
	}
}

func cellular_noise(noise_map *[HEIGHT][WIDTH]uint8) {
	cells := make([][2]float64, CELL_COUNT)
	for i := range cells {
		t := time.Now()
		nsec := t.UnixNano()
		s := rand.NewSource(int64(nsec))
		r := rand.New(s)
		cells[i] = [2]float64{r.Float64() * BOARD_SIZE, r.Float64() * BOARD_SIZE}
	}
	for yi, yv := range noise_map {
		for xi, _ := range yv {
			current_point := [2]float64{float64(xi) / float64(WIDTH) * BOARD_SIZE, float64(yi) / float64(HEIGHT) * BOARD_SIZE}
			distances := make([]float64, CELL_COUNT)
			const cell_distance_id = 0
			for i, v := range cells {
				distance := distance(current_point, v)
				distances[i] = distance
			}
			sort.Float64s(distances)
			noise_map[yi][xi] = uint8(clamp(distances[cell_distance_id], 0.0, 1.0) * 255)
		}
	}
}

func cellular_noise_gif(noise_map [HEIGHT][WIDTH]uint8, frames_per_rotation int) {
	cells := make([][2]float64, CELL_COUNT)
	radii := make([]float64, CELL_COUNT)
	start_pos := make([]float64, CELL_COUNT)
	for i := range cells {
		t := time.Now()
		nsec := t.UnixNano()
		s := rand.NewSource(int64(nsec))
		r := rand.New(s)
		cells[i] = [2]float64{r.Float64() * BOARD_SIZE, r.Float64() * BOARD_SIZE}
		s = rand.NewSource(int64(nsec) + 1)
		r = rand.New(s)
		radii[i] = (1.0 - (r.Float64() * 0.15)) * MOVEMENT_FACTOR
		s = rand.NewSource(int64(nsec) - 1)
		r = rand.New(s)
		start_pos[i] = r.Float64() * 10
	}
	frames := make([][HEIGHT][WIDTH]uint8, frames_per_rotation)
	for t := 0; t < frames_per_rotation; t++ {
		for yi, yv := range noise_map {
			for xi, _ := range yv {
				current_point := [2]float64{float64(xi) / float64(WIDTH) * BOARD_SIZE, float64(yi) / float64(HEIGHT) * BOARD_SIZE}
				distances := make([]float64, CELL_COUNT)
				const cell_distance_id = 0
				for i, v := range cells {
					distance := distance(current_point, [2]float64{v[0] + (math.Sin(2 * math.Pi * ((float64(t) / float64(frames_per_rotation)) + start_pos[i]) * radii[i])), v[1] + (math.Cos(2*math.Pi*((float64(t)/float64(frames_per_rotation))+start_pos[i])) * radii[i])})
					distances[i] = distance
				}
				sort.Float64s(distances)
				noise_map[yi][xi] = uint8(clamp(distances[cell_distance_id], 0.0, 1.0) * 255)
			}
		}
		frames[t] = noise_map
	}
	for i, noise_map := range frames {
		frame := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))
		write_image(noise_map, frame)
		name := fmt.Sprintf("frame%d.png", i)
		fmt.Println(name)
		f, err := os.Create(name)
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = png.Encode(f, frame)
		if err != nil {
			log.Fatal(err)
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

	cellular_noise_gif(noise_map, 60)
	// img := image.NewRGBA(image.Rect(0, 0, WIDTH, HEIGHT))

	// write_image(noise_map, img)

	// //Creates an image
	// f, err := os.Create("result.png")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer f.Close()

	// // Encodes and saves to png
	// err = png.Encode(f, img)
	// if err != nil {
	// 	log.Fatal(err)
	// }
}
