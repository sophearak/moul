package main

import (
	"encoding/json"
	"fmt"
	"github.com/sophearak/moul/moul"
	"image"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/denisbrodbeck/sqip"
	"gopkg.in/h2non/bimg.v1"
)

// resize image
func manipulate(size int, path string) {
	options := bimg.Options{
		Width: size,
	}
	buffer, err := bimg.Read(path)
	if err != nil {
		log.Fatal(err)
	}
	newImage, err := bimg.NewImage(buffer).Process(options)
	if err != nil {
		log.Fatal(err)
	}
	out := ".moul/" + filepath.Dir(path) + "/" + strconv.Itoa(size) + "/" + filepath.Base(path)
	bimg.Write(out, newImage)
}

// get image from given folder
func getImage(path string) []string {
	var files []string
	// folder to walk through
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		if filepath.Ext(path) == ".jpg" || filepath.Ext(path) == ".png" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		panic(err)
	}
	return files
}

// make placeholder image
func sqipy(path string) {
	workSize := 256
	count := 1
	mode := 0
	alpha := 128
	repeat := 0
	workers := runtime.NumCPU()
	background := ""
	svg, width, height, err := sqip.Run(path, workSize, count, mode, alpha, repeat, workers, background)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(width, height, workers)

	out := strings.TrimSuffix(".moul/"+path, filepath.Ext(path)) + ".svg"
	dest := filepath.Dir(out) + "/svg/" + filepath.Base(out)
	if err := sqip.SaveFile(dest, svg); err != nil {
		log.Fatal(err)
	}
}

func getImageDimension(path string) (int, int) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}

	image, _, err := image.DecodeConfig(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", path, err)
	}
	return image.Width, image.Height
}

func generate(path string, sizes []int) {
	files := getImage(path)

	for _, file := range files {
		for _, size := range sizes {
			manipulate(size, file)
		}
		sqipy(file)
	}
}

func main() {

	//generate("photos/cover", []int{2560, 1920, 1280, 960, 640, 480, 320})
	//generate("photos/profile", []int{1024, 320})
	//generate("photos/collection", []int{2048, 750})

	collection := getImage("./.moul/photos/collection")

	mc := make([]moul.Collection, 0)

	// to be clean up
	for index, photo := range collection {
		for i, p := range collection {
			if index != i && filepath.Base(photo) == filepath.Base(p) {
				fsindex, err := os.Stat(photo)
				if err != nil {
					log.Fatalln(err)
				}
				fsi, err := os.Stat(p)
				if err != nil {
					log.Fatalln(err)
				}
				if fsindex.Size() < fsi.Size() {
					width, height := getImageDimension(collection[index])
					widthHd, heightHd := getImageDimension(collection[i])
					svg := strings.TrimSuffix(filepath.Base(collection[index]), filepath.Ext(collection[index]))
					mc = append(mc, moul.Collection{
						Name:     filepath.Base(collection[index]),
						Src:      collection[index],
						Svg:      ".moul/photos/collection/svg/" + svg + ".svg",
						Width:    width,
						Height:   height,
						SrcHd:    collection[i],
						WidthHd:  widthHd,
						HeightHd: heightHd,
					})
				}
			}
		}
	}

	mcjson, err := json.Marshal(mc)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(string(mcjson))

}
