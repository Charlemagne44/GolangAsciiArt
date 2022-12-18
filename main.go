package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"os"
	"sort"

	"github.com/disintegration/imaging"
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
		return key < rm.Keys[i].L
	})

	i -= 1
	if i >= 0 && i < len(rm.Keys) && key <= rm.Keys[i].U {
		return rm.Values[i], true
	}
	return "", false
}

func main() {
	// filename for image to edit
	filename := flag.String("file", "", "Name of file to turn into ascii art")
	// 1.0 scaling factor will create 3 ascii characteres per 1 pixel
	scalingFactor := flag.Float64("scale", 1.0, "Scaling factor to resize the image")
	// good for sharpening real images with lots of noise and vibrance for better ascii result
	contrastFactor := flag.Int("contrast", 0, "Contrast value to apply to image")
	// name for an output text file
	outputName := flag.String("out", "", "name for output to be written")
	// print to terminal option
	printOption := flag.Bool("print", false, "print art to stdout")
	// html output option
	htmlOption := flag.String("html", "", "print to html file")
	flag.Parse()

	if *filename == "" {
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
		img = imaging.Resize(img, width, height, imaging.Lanczos)
	}

	// if contrast flag declared
	if *contrastFactor != 0 {
		img = imaging.AdjustContrast(img, float64(*contrastFactor))
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

	if *printOption {
		printArt(rangeMap, brightnessArr)
	}

	if *outputName != "" {
		writeArt(rangeMap, brightnessArr, *outputName)
	}

	if *htmlOption != "" {
		writeHTML(rangeMap, brightnessArr, *htmlOption)
	}
}

func writeHTML(rangeMap RangeMap, brightnessArr [][]int, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	file.WriteString("<body>\n")
	file.WriteString("<tt>")

	for _, row := range brightnessArr {
		for _, col := range row {
			value, _ := rangeMap.Get(int(col))
			file.WriteString(value)
			file.WriteString(value)
			file.WriteString(value)
		}
		file.WriteString("<br>\n")
	}
	file.WriteString("</tt>")
	file.WriteString("</body>")
}

func writeArt(rangeMap RangeMap, brightnessArr [][]int, fileName string) {
	file, err := os.Create(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, row := range brightnessArr {
		for _, col := range row {
			value, _ := rangeMap.Get(int(col))
			file.WriteString(value)
			file.WriteString(value)
			file.WriteString(value)
		}
		file.WriteString("\n")
	}
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
