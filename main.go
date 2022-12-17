package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"os"
	"sort"

	"github.com/nfnt/resize"
)

// ----- Rangemap implementation from d-schmidt@github.com: rangemap.go
type Range struct {
	L int
	U int
}

type RangeMap struct {
	Keys   []Range
	Values []string
}

func (rm RangeMap) Get(key int) (string, bool) {
	i := sort.Search(len(rm.Keys), func(i int) bool {
		// fmt.Printf("search %v at index %d for %v is %v\n", rm.Keys[i], i, key, key < rm.Keys[i].L)
		return key < rm.Keys[i].L
	})

	i -= 1
	if i >= 0 && i < len(rm.Keys) && key <= rm.Keys[i].U {
		return rm.Values[i], true
	}
	return "", false
}

func main() {
	// ----- flags
	filename := flag.String("file", "", "Name of file to turn into ascii art")
	scalingFactor := flag.Float64("scale", 1.0, "Scaling factor to resize the image")
	flag.Parse()

	if *filename == "" {
		fmt.Println("need valid filename")
		return
	}

	// read the image and construct the brightness array
	imgfile, err := os.Open(*filename)
	if err != nil {
		panic(err)
	}
	defer imgfile.Close()

	imgCfg, err := jpeg.DecodeConfig(imgfile)
	if err != nil {
		panic(err)
	}

	imgfile.Seek(0, 0)
	img, err := jpeg.Decode(imgfile)
	if err != nil {
		panic(err)
	}

	width := imgCfg.Width
	height := imgCfg.Height

	// resize image if flag declared
	if *scalingFactor != 1.0 {
		width = int(float64(imgCfg.Width) * *scalingFactor)
		height = int(float64(imgCfg.Height) * *scalingFactor)
		img = resize.Resize(uint(width),
			uint(height),
			img,
			resize.Lanczos3,
		)
		fmt.Println("Resized dim:", width, height)
	}

	// initialize brightness array
	var brightnessArr [][]int = make([][]int, height)
	for i := 0; i < height; i++ {
		brightnessArr[i] = make([]int, width)
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			r, g, b, _ := img.At(x, y).RGBA()
			average := ((r / 256.0) + (g / 256.0) + (b / 256.0)) / 3.0
			brightnessArr[y][x] = int(average)
		}
	}

	// construct values and range for map search
	brightString := "`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	var brightList []string
	for _, char := range brightString {
		brightList = append(brightList, string(char))
	}

	var rangeList []Range
	chunk := 255.0 / 64.0
	for i := chunk; i <= 255; i += chunk {
		rangeList = append(rangeList, Range{L: int(i - chunk), U: int(i)})
	}

	rangeMap := RangeMap{
		Values: brightList,
		Keys:   rangeList,
	}
	_ = rangeMap

	printArt(rangeMap, brightnessArr)

}

func printArt(rangeMap RangeMap, brightnessArr [][]int) {
	for _, row := range brightnessArr {
		for _, col := range row {
			value, _ := rangeMap.Get(int(col))
			fmt.Printf("%s", value)
			fmt.Printf("%s", value)
			fmt.Printf("%s", value)
		}
		fmt.Printf("\n")
	}
}
