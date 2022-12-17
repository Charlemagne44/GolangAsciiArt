package main

import (
	"flag"
	"fmt"
	"sort"

	"gopkg.in/gographics/imagick.v2/imagick"
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
	flag.Parse()

	if *filename == "" {
		fmt.Println("need valid filename")
		return
	}

	// ----- load the imagick library
	imagick.Initialize()
	defer imagick.Terminate()

	mw := imagick.NewMagickWand()
	defer mw.Destroy()

	err := mw.ReadImage(*filename)
	if err != nil {
		panic(err)
	}

	width := mw.GetImageWidth()
	height := mw.GetImageHeight()
	fmt.Println("Width:", width, "Height:", height)

	// ----- create 2d array of pixel tuple native types
	fmt.Println("Creating pixel matrix")
	var imgArr [][]*imagick.PixelWand
	iterator := mw.NewPixelIterator()
	defer iterator.Destroy()
	for i := 0; i < int(height); i++ {
		pixels := iterator.GetNextIteratorRow()
		imgArr = append(imgArr, pixels)
		iterator.SyncIterator()
	}

	// create the brigntess array
	fmt.Println("Creating the brightness array")
	brightnessArr := make([][]int, height)
	for i := 0; i < int(height); i++ {
		brightnessArr[i] = make([]int, width)
	}
	for i, pixelRow := range imgArr {
		for j, pixel := range pixelRow {
			// fmt.Println("value being added to barr:", int((pixel.GetRed()*255)+(pixel.GetGreen()*255)+(pixel.GetBlue()*255))/3)
			brightnessArr[i][j] = int((pixel.GetRed()*255)*(pixel.GetGreen()*255)*(pixel.GetBlue()*255)) / 3
		}
	}
	fmt.Println("Created the brightness array!")
	fmt.Println("brightness arr size:", len(brightnessArr), len(brightnessArr[0]))

	// ----- create the ascii brightness map
	fmt.Println("Creating ascii brightness maps")

	// construct values and range for map search
	brightString := "`^\",:;Il!i~+_-?][}{1)(|\\/tfjrxnuvczXYUJCLQ0OZmwqpdbkhao*#MW&8%B@$"
	var brightList []string
	for _, char := range brightString {
		brightList = append(brightList, string(char))
	}
	// fmt.Println("brightlist:", brightList)

	var rangeList []Range
	chunk := 255.0 / 67.0
	for i := chunk; i <= 255; i += chunk {
		rangeList = append(rangeList, Range{L: int(i - chunk), U: int(i)})
	}

	fmt.Println("Created ascii brightness map!")

	rangeMap := RangeMap{
		Values: brightList,
		Keys:   rangeList,
	}
	_ = rangeMap

	// ----- Print out the characters to terminal with the range map and brightness values
	// for _, row := range brightnessArr {
	// 	for _, col := range row {
	// 		value, _ := rangeMap.Get(int(col))
	// 		fmt.Printf("%s", value)
	// 	}
	// 	fmt.Printf("\n")
	// }
}
