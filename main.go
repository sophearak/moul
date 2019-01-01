package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/denisbrodbeck/sqip"
	"github.com/gobuffalo/plush"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/html"
	"github.com/tdewolff/minify/js"
	"github.com/tdewolff/minify/json"
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
	viper.SetConfigName("config")
	viper.AddConfigPath(".")

	//generate("photos/cover", []int{2560, 1920, 1280, 960, 640, 480, 320})
	//generate("photos/profile", []int{1024, 320})
	//generate("photos/collection", []int{2048, 750})

	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	htmlTemplate := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta http-equiv="X-UA-Compatible" content="ie=edge">
  <title><%= siteName %></title>
  <style type="text/css">
    body {
      margin: 0
    }
  </style>
</head>
<body>

  <h1><%= siteName %></h1>

  <script type="text/javascript">
    document.addEventListener("DOMContentLoaded", () => {
      console.log("DOM ready!");
    })
  </script>
</body>
</html>`

	ctx := plush.NewContext()
	ctx.Set("siteName", viper.Get("site.name"))

	s, err := plush.Render(htmlTemplate, ctx)
	if err != nil {
		log.Fatal(err)
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)

	m.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
	})

	mt, err := m.String("text/html", s)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(mt)
}
