package moul

import (
	"bytes"
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
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
	"gopkg.in/h2non/bimg.v1"
)

type Collection struct {
	Name     string `json:"name"`
	Src      string `json:"src"`
	Srcset   string `json:"srcset"`
	SrcHd    string `json:"src_hd"`
	Width    int    `json:"width"`
	WidthHd  int    `json:"width_hd"`
	Height   int    `json:"height"`
	HeightHd int    `json:"height_hd"`
}

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

// get image size without open
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

// generate photos
func Generate(path string, sizes []int) {
	files := getImage(path)

	for _, file := range files {
		for _, size := range sizes {
			manipulate(size, file)
		}
		sqipy(file)
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// build actual html
func Build() {
	template := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<title><%= name %></title>

	<style type="text/css">
		*,
		:before,
		:after {
			box-sizing: border-box;
		}
		body {
			margin: 0;
			font-family: -apple-system,BlinkMacSystemFont,"Segoe UI",Roboto,"Helvetica Neue",Arial,sans-serif;
			line-height: 1.5;
			-webkit-font-smoothing: antialiased;
			-moz-osx-font-smoothing: grayscale;
			text-rendering: optimizeLegibility;
		}
		.cover {
			position: relative;
			height: 500px;
		}
		.cover-wrap {
			position: absolute;
			top: 0;
			left: 0;
			width: 100%;
			height: 100%;
			text-align: center;
		}
		.cover-wrap picture {
			position: relative;
			display: block;
			height: 100%;
			z-index: -1;
		}
		.cover-wrap picture img {
			width: 100%;
			height: 100%;
			-o-object-fit: cover;
			object-fit: cover;
		}
		.profile {
			text-align: center;
			margin: 32px 0 16px;
			font-size: 0;
		}
		.profile a img {
			width: 120px;
			height: 120px;
			border-radius: 50%;
			border: 2px solid transparent;
			transition: border 250ms cubic-bezier(0.4, 0, 0.2, 1), box-shadow 250ms cubic-bezier(0.4, 0, 0.2, 1);
		}
		.profile a img:hover {
			box-shadow: 0 1px 2px 0 rgba(60,64,67,0.302), 0 2px 6px 2px rgba(60,64,67,0.149);
		}
		.name {
			margin: 0 0 8px;
			text-align: center;
			font-weight: 900;
			color: #111;
			line-height: 1;
		}
		.info {
			max-width: 500px;
			text-align: center;
			margin: 0 auto 32px;
			color: #444;
		}
		.collection {
			position: relative;
			margin: 10px auto;
		}
		.collection figure {
			margin: 0;
		}
		.collection figure a {
			display: block;
			font-size: 0;
			float: left;
		}
		.collection figure a img {
			position: absolute;
			background-size: cover;
		}
	</style>
	<link rel="stylesheet" href="moul-collection.min.css">
</head>
<body>
    <div class="cover">
        <div class="cover-wrap">
            <picture>
                <img class="js-img" src="<%= coverSrc %>" srcset="<%= coverSrcset %>" alt="cover" sizes="1px">
            </picture>
        </div>
    </div>
    <div class="profile">
        <a href="./photos/profile/1024/<%= pfn %>">
            <img class="js-img" src="./photos/profile/320/<%= pfn %>" srcset="./photos/profile/320/<%= pfn %> 320w, ./photos/profile/svg/<%= pfns %>.svg 10w" alt="<%= name %>" sizes="1px">
        </a>
    </div>
    <h1 class="name"><%= name %></h1>
    <p class="info"><%= bio %></p>
	<div id="collection">
	<div>
	
<input type="hidden" id="photo-collection" value="<%= collection %>">

	<script type="text/javascript">
		const $ = document.querySelector.bind(document);
		const $$ = document.querySelectorAll.bind(document);

		document.addEventListener("DOMContentLoaded", () => {
			const photoTags = $$('.js-img');
			setTimeout(() => {
				photoTags.forEach(photo => {
					const sizes = Math.ceil(
					  (photo.getBoundingClientRect().width / window.innerWidth) * 100
					);
					photo.setAttribute('sizes', sizes + 'vw');
				});
			}, 500);
		});
	</script>
	<script src="moul-collection.min.js"></script>
</body>
</html>`

	// get configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()

	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	// get cover
	base := "./photos/cover/"
	var coverSrcset bytes.Buffer
	cover := getImage("./.moul/photos/cover")
	coverName := filepath.Base(cover[0])
	coverSrc := base + "1920/" + filepath.Base(cover[0])
	coverSizes := []int{2560, 1920, 1280, 960, 640, 480, 320}
	for i, size := range coverSizes {
		fp := base + strconv.Itoa(size) + "/" + coverName + " " + strconv.Itoa(size) + "w, "
		coverSrcset.WriteString(fp)
		le := len(coverSizes) - 1

		if i == le {
			coverSrcset.WriteString(base + "svg/" + strings.TrimSuffix(coverName, filepath.Ext(coverName)) + ".svg 30w")
		}
	}

	// get profile
	profile := getImage("./.moul/photos/profile")
	pfn := filepath.Base(profile[0])
	pfns := strings.TrimSuffix(filepath.Base(profile[0]), filepath.Ext(filepath.Base(profile[0])))

	// get photo collection
	collection := getImage("./.moul/photos/collection")

	mc := make([]Collection, 0)
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
					fs := filepath.Base(collection[index])
					base := "./photos/collection/750/"
					baseHd := "./photos/collection/2048/"

					svg := strings.TrimSuffix(fs, filepath.Ext(fs))
					mc = append(mc, Collection{
						Name:     fs,
						Src:      base + fs,
						Srcset:   base + fs + " 300w, ./photos/collection/svg/" + svg + ".svg 20w",
						Width:    width,
						Height:   height,
						SrcHd:    baseHd + fs,
						WidthHd:  widthHd,
						HeightHd: heightHd,
					})
				}
			}
		}
	}

	mcj, err := json.Marshal(mc)
	check(err)

	// push data to template
	ctx := plush.NewContext()
	ctx.Set("coverSrcset", coverSrcset.String())
	ctx.Set("coverSrc", coverSrc)

	ctx.Set("pfn", pfn)
	ctx.Set("pfns", pfns)

	ctx.Set("name", viper.Get("site.name"))
	ctx.Set("bio", viper.Get("site.bio"))

	ctx.Set("twitter", viper.Get("social.twitter"))
	ctx.Set("facebook", viper.Get("social.facebook"))
	ctx.Set("instagram", viper.Get("social.instagram"))

	ctx.Set("collection", string(mcj))

	s, err := plush.Render(template, ctx)
	check(err)

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)
	m.AddFuncRegexp(regexp.MustCompile("^(application|text)/(x-)?(java|ecma)script$"), js.Minify)
	//m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)

	m.Add("text/html", &html.Minifier{
		KeepDocumentTags: true,
	})

	//mt, err := m.String("text/html", s)
	//check(err)

	ioutil.WriteFile("./.moul/index.html", []byte(s), 0644)
}
