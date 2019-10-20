package moul

import (
	"encoding/json"
	"fmt"
	"image"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/disintegration/imaging"
	"github.com/fsnotify/fsnotify"
	"github.com/generaltso/vibrant"
	"github.com/gobuffalo/plush"
	"github.com/spf13/viper"
	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/css"
	"github.com/tdewolff/minify/v2/html"
)

// Collection struct
type Collection struct {
	Name     string `json:"name"`
	Src      string `json:"src"`
	Color    string `json:"color"`
	SrcHd    string `json:"src_hd"`
	Width    int    `json:"width"`
	WidthHd  int    `json:"width_hd"`
	Height   int    `json:"height"`
	HeightHd int    `json:"height_hd"`
}

func getDominantDarkColor(path string) string {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	img, _, err := image.Decode(f)
	if err != nil {
		log.Fatal(err)
	}

	colors, _ := vibrant.NewPalette(img, 256)
	c := colors.ExtractAwesome()

	vc, ok := c["DarkMuted"]

	if !ok {
		vc, ok = c["Muted"]
		if !ok {
			vc, ok = c["Muted"]
			if !ok {
				vc, ok = c["Vibrant"]
			}
		}
	}

	return vc.Color.RGBHex()
}

// resize image
func manipulate(size int, path string) {
	src, err := imaging.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	out := ".moul/" + filepath.Dir(path) + "/" + strconv.Itoa(size) + "/" + filepath.Base(path)

	newImage := imaging.Resize(src, size, 0, imaging.Lanczos)

	err = imaging.Save(newImage, out)
	if err != nil {
		log.Fatal(err)
	}
}

// GetImageDimension func get image size without open
func GetImageDimension(path string) (int, int) {
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

// Generate photos
func Generate(path string, sizes []int) {
	s := spinner.New(spinner.CharSets[35], 100*time.Millisecond)
	s.Prefix = "Crafting "
	files := getImage(path)

	if len(files) == 0 {
		if path != "photos/collection" {
			fmt.Println("Add a photo inside " + path)
		} else {
			fmt.Println("Add photos inside " + path)
		}
		os.Exit(1)
	}

	s.Start()
	for _, file := range files {
		for _, size := range sizes {
			manipulate(size, file)
		}
	}
	s.Stop()
}

// Build actual html
func Build() {
	// get cover
	cover := getImage("./.moul/photos/cover")
	coverName := filepath.Base(cover[0])
	coverColor := getDominantDarkColor("./.moul/photos/cover/620/" + coverName)

	// get profile
	profile := getImage("./.moul/photos/profile")
	profileName := filepath.Base(profile[0])
	profileColor := getDominantDarkColor("./.moul/photos/profile/320/" + profileName)

	// get photo collection
	collection := getImage("./.moul/photos/collection")

	// get configuration
	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}

	mc := make([]Collection, 0)
	for _, photo := range collection {
		if strings.Contains(photo, "2048") {
			fs := filepath.Base(photo)
			widthHd, heightHd := GetImageDimension(photo)
			width, height := GetImageDimension(".moul/photos/collection/750/" + fs)
			base := "./photos/collection/750/"
			baseHd := "./photos/collection/2048/"
			customBg := viper.Get("background")

			var color string
			if customBg != nil {
				color = getDominantDarkColor(photo)
			} else {
				color = "rgba(0, 0, 0, .95)"
			}

			mc = append(mc, Collection{
				Name:     fs,
				Src:      base + fs,
				Color:    color,
				Width:    width,
				Height:   height,
				SrcHd:    baseHd + fs,
				WidthHd:  widthHd,
				HeightHd: heightHd,
			})
		}
	}

	mcj, _ := json.Marshal(mc)

	buildHTML(coverName, coverColor, profileName, profileColor, string(mcj))
}

func buildHTML(coverName, coverColor, profileName, profileColor, collection string) {
	template := `<!DOCTYPE html>
<html lang="en">
<head>
	<meta charset="utf-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<meta http-equiv="X-UA-Compatible" content="ie=edge">
	<base href="/">
	<title><%= name %></title>
	<meta name="description" content="<%= bio %>">

	<meta name="twitter:card" content="summary">
	<meta name="twitter:title" content="<%= name %>">
	<meta name="twitter:description" content="<%= bio %>">
	<%= if (twitter) { %>
	<meta name="twitter:site" content="<%= twitter %>">
	<meta name="twitter:creator" content="<%= twitter %>">
	<% } %>
	<meta name="twitter:image:src" content="<%= url %>/photos/cover/1024/<%= coverName %>">

	<meta name="og:title" content="<%= name %>">
	<meta name="og:description" content="<%= bio %>">
	<meta name="og:image" content="<%= url %>/photos/cover/1200/<%= coverName %>">
	<meta name="og:url" content="<%= url %>">
	<meta name="og:site_name" content="<%= name %>">
	<meta name="og:type" content="website">
	<style>
		.social {
			display: flex;
			justify-content: center;
			align-items: center;
			margin: 0 0 32px;
		}
		.social a {
			color: #111;
			line-height: 0;
			margin: 0 8px 16px;
			transition: color 150ms cubic-bezier(0.4, 0, 0.2, 1);
		}
		.social a:hover {
			color: #888;
		}
	</style>

	<link rel="stylesheet" href="assets/index.css">
</head>
<body>
	<header>
		<div class="banner">
			<picture>
				<source media="(max-width: 600px)" srcset="./photos/cover/620/<%= coverName %>">
				<source media="(min-width: 601px)" srcset="./photos/cover/1280/<%= coverName %>">
				<source media="(min-width: 1201px)" srcset="./photos/cover/2560/<%= coverName %>">
				<img
					data-src="./photos/cover/1280/<%= coverName %>"
					alt="cover"
					class="lazy"
					style="background: <%= coverColor %>">
			</picture>
		</div>
	</header>
	<div class="profile">
		<a href="./photos/profile/1024/<%= profileName %>">
			<img
				data-src="./photos/profile/320/<%= profileName %>"
				alt="profile photo"
				class="lazy">
		</a>
	</div>
	<h1><%= name %></h1>
	<%= if (len(twitter) > 0 || len(youtube) > 0 || len(facebook) > 0 || len(instagram) > 0) { %>
		<%= if (bio && len(bio) > 0) { %>
		<p style="margin: 0 auto 16px"><%= bio %></p>
		<% } %>
	<div class="social">
		<%= if (twitter) { %>
			<a href="https://twitter.com/<%= twitter %>">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-twitter"><path d="M23 3a10.9 10.9 0 0 1-3.14 1.53 4.48 4.48 0 0 0-7.86 3v1A10.66 10.66 0 0 1 3 4s-4 9 5 13a11.64 11.64 0 0 1-7 2c9 5 20 0 20-11.5a4.5 4.5 0 0 0-.08-.83A7.72 7.72 0 0 0 23 3z"></path></svg>
			</a>
		<% } %>
		<%= if (youtube) { %>
			<a href="https://www.youtube.com/<%= youtube %>">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-youtube"><path d="M22.54 6.42a2.78 2.78 0 0 0-1.94-2C18.88 4 12 4 12 4s-6.88 0-8.6.46a2.78 2.78 0 0 0-1.94 2A29 29 0 0 0 1 11.75a29 29 0 0 0 .46 5.33A2.78 2.78 0 0 0 3.4 19c1.72.46 8.6.46 8.6.46s6.88 0 8.6-.46a2.78 2.78 0 0 0 1.94-2 29 29 0 0 0 .46-5.25 29 29 0 0 0-.46-5.33z"></path><polygon points="9.75 15.02 15.5 11.75 9.75 8.48 9.75 15.02"></polygon></svg>
			</a>
		<% } %>
		<%= if (facebook) { %>
			<a href="https://www.facebook.com/<%= facebook %>">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-facebook"><path d="M18 2h-3a5 5 0 0 0-5 5v3H7v4h3v8h4v-8h3l1-4h-4V7a1 1 0 0 1 1-1h3z"></path></svg>
			</a>
		<% } %>
		<%= if (instagram) { %>
			<a href="https://www.instagram.com/<%= instagram %>">
				<svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" class="feather feather-instagram"><rect x="2" y="2" width="20" height="20" rx="5" ry="5"></rect><path d="M16 11.37A4 4 0 1 1 12.63 8 4 4 0 0 1 16 11.37z"></path><line x1="17.5" y1="6.5" x2="17.5" y2="6.5"></line></svg>
			</a>
		<% } %>
	</div>
	<% } else { %>
		<%= if (bio && len(bio) > 0) { %>
			<p><%= bio %></p>
		<% } %>
	<% } %>

	<div id="moul"></div>
	<input type="hidden" id="photo-collection" value="<%= collection %>">

	<script src="assets/index.js"></script>
	<div class="pswp" tabindex="-1" role="dialog" aria-hidden="true">
		<div class="pswp__bg"></div>
		<div class="pswp__scroll-wrap">
			<div class="pswp__container">
				<div class="pswp__item"></div>
				<div class="pswp__item"></div>
				<div class="pswp__item"></div>
			</div>
			<div class="pswp__ui pswp__ui--hidden">
				<div class="pswp__top-bar">
					<div class="pswp__counter"></div>
					<button class="pswp__button pswp__button--close" title="Close (Esc)"></button>
					<button class="pswp__button pswp__button--share" title="Share"></button>
					<button class="pswp__button pswp__button--fs" title="Toggle fullscreen"></button>
					<button class="pswp__button pswp__button--zoom" title="Zoom in/out"></button>
					<div class="pswp__preloader">
						<div class="pswp__preloader__icn">
							<div class="pswp__preloader__cut">
								<div class="pswp__preloader__donut"></div>
							</div>
						</div>
					</div>
				</div>
				<div class="pswp__share-modal pswp__share-modal--hidden pswp__single-tap">
					<div class="pswp__share-tooltip"></div>
				</div>
				<button class="pswp__button pswp__button--arrow--left" title="Previous (arrow left)"></button>
				<button class="pswp__button pswp__button--arrow--right" title="Next (arrow right)"></button>
				<div class="pswp__caption">
					<div class="pswp__caption__center"></div>
				</div>
			</div>
		</div>
	</div>
</body>
</html>`

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		buildHTML(coverName, coverColor, profileName, profileColor, collection)
		fmt.Println("Updated")
	})

	// push data to template
	ctx := plush.NewContext()
	ctx.Set("coverName", coverName)
	ctx.Set("coverColor", coverColor)

	ctx.Set("profileName", profileName)
	ctx.Set("profileColor", profileColor)

	ctx.Set("url", viper.Get("site.url"))
	ctx.Set("name", viper.Get("site.name"))
	ctx.Set("bio", viper.Get("site.bio"))

	ctx.Set("twitter", viper.Get("social.twitter"))
	ctx.Set("youtube", viper.Get("social.youtube"))
	ctx.Set("facebook", viper.Get("social.facebook"))
	ctx.Set("instagram", viper.Get("social.instagram"))

	ctx.Set("collection", collection)

	s, err := plush.Render(template, ctx)
	if err != nil {
		log.Fatal(err)
	}

	m := minify.New()
	m.AddFunc("text/css", css.Minify)
	m.AddFunc("text/html", html.Minify)

	mt, err := m.String("text/html", s)

	ioutil.WriteFile("./.moul/index.html", []byte(mt), 0644)
}
