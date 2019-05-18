package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/gobuffalo/packr/v2"
	copy2 "github.com/otiai10/copy"
	"github.com/sophearak/moul/moul"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new [string folder name]",
	Short: "Create new photo collection",
	Long:  ``,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(args[0]); !os.IsNotExist(err) {
			fmt.Println("Folder already exists")
		}

		err := os.MkdirAll(args[0], os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}

		folders := []string{args[0] + "/photos", args[0] + "/.moul/assets", args[0] + "/photos/cover", args[0] + "/photos/profile", args[0] + "/photos/collection"}
		for _, folder := range folders {
			os.MkdirAll(folder, os.ModePerm)
		}

		coverSizes := []string{"2560", "1920", "1280", "1200", "1024", "960", "640", "480", "320", "svg"}
		for _, cs := range coverSizes {
			os.MkdirAll(args[0]+"/.moul/photos/cover/"+cs, os.ModePerm)
		}

		profileSizes := []string{"1024", "320", "svg"}
		for _, ps := range profileSizes {
			os.MkdirAll(args[0]+"/.moul/photos/profile/"+ps, os.ModePerm)
		}

		collectionSizes := []string{"2048", "750", "svg"}
		for _, cs := range collectionSizes {
			os.MkdirAll(args[0]+"/.moul/photos/collection/"+cs, os.ModePerm)
		}

		box := packr.New("assets", "../assets")

		dpng, _ := box.FindString("default-skin.76672929.png")
		ioutil.WriteFile(args[0]+"/.moul/assets/default-skin.76672929.png", []byte(dpng), 0644)

		dsvg, _ := box.FindString("default-skin.a5214274.svg")
		ioutil.WriteFile(args[0]+"/.moul/assets/default-skin.a5214274.svg", []byte(dsvg), 0644)

		pgif, _ := box.FindString("preloader.f75eb900.gif")
		ioutil.WriteFile(args[0]+"/.moul/assets/preloader.f75eb900.gif", []byte(pgif), 0644)

		mcss, _ := box.FindString("moul-collection.min.css")
		ioutil.WriteFile(args[0]+"/.moul/assets/moul-collection.min.css", []byte(mcss), 0644)

		mjs, _ := box.FindString("moul-collection.min.js")
		ioutil.WriteFile(args[0]+"/.moul/assets/moul-collection.min.js", []byte(mjs), 0644)

		conf, _ := box.FindString("config.json")
		ioutil.WriteFile(args[0]+"/config.json", []byte(conf), 0644)
	},
}

var devCmd = &cobra.Command{
	Use:   "dev",
	Short: "Generate photo collection",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		moul.Generate("photos/cover", []int{2560, 1920, 1280, 1200, 1024, 960, 640, 480, 320})
		moul.Generate("photos/profile", []int{1024, 320})
		moul.Generate("photos/collection", []int{2048, 750})

		moul.Build()

		fs := http.FileServer(http.Dir("./.moul/"))
		http.Handle("/", fs)

		fmt.Println("Serve:  http://localhost:12345")
		fmt.Println("Ctrl + C to exit")
		http.ListenAndServe(":12345", nil)
	},
}

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Building collection for deployment",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat("./dist"); !os.IsNotExist(err) {
			removeContents("./dist")
		}
		os.MkdirAll("./dist", os.ModePerm)
		err := copy2.Copy("./.moul", "./dist")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Photo collection is built")
	},
}

var igCmd = &cobra.Command{
	Use:   "ig",
	Short: "Generate photos for Instagram grid",
	Long:  "moul ig photo.jpg",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		const s = 1080

		photo := args[0]
		e := filepath.Ext(photo)
		fn := filepath.Base(photo)

		// only working with jpg and png
		if e == ".jpg" || e == ".jpeg" || e == ".png" {
			w, h := moul.GetImageDimension(photo)

			// width must be > 3240px
			if w < 3240 || h < 3240 {
				fmt.Println("Minimum W x H: 3240px x 3240px")
				os.Exit(1)
			}
			src, err := imaging.Open(photo)
			if err != nil {
				log.Fatal(err)
			}

			img := imaging.Resize(src, 3240, 3240, imaging.Lanczos)
			// resize
			tl := imaging.CropAnchor(img, s, s, imaging.TopLeft)
			t := imaging.CropAnchor(img, s, s, imaging.Top)
			tr := imaging.CropAnchor(img, s, s, imaging.TopRight)
			ml := imaging.CropAnchor(img, s, s, imaging.Left)
			m := imaging.CropAnchor(img, s, s, imaging.Center)
			mr := imaging.CropAnchor(img, s, s, imaging.Right)
			bl := imaging.CropAnchor(img, s, s, imaging.BottomLeft)
			b := imaging.CropAnchor(img, s, s, imaging.Bottom)
			br := imaging.CropAnchor(img, s, s, imaging.BottomRight)

			photoName := strings.TrimSuffix(fn, filepath.Ext(fn))
			os.Mkdir(photoName, os.ModePerm)

			imaging.Save(br, filepath.Join(photoName, "1"+filepath.Ext(fn)))
			imaging.Save(b, filepath.Join(photoName, "2"+filepath.Ext(fn)))
			imaging.Save(bl, filepath.Join(photoName, "3"+filepath.Ext(fn)))
			imaging.Save(mr, filepath.Join(photoName, "4"+filepath.Ext(fn)))
			imaging.Save(m, filepath.Join(photoName, "5"+filepath.Ext(fn)))
			imaging.Save(ml, filepath.Join(photoName, "6"+filepath.Ext(fn)))
			imaging.Save(tr, filepath.Join(photoName, "7"+filepath.Ext(fn)))
			imaging.Save(t, filepath.Join(photoName, "8"+filepath.Ext(fn)))
			imaging.Save(tl, filepath.Join(photoName, "9"+filepath.Ext(fn)))
		} else {
			fmt.Println("Support only jpg and png!")
			os.Exit(1)
		}
	},
}

func removeContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// Execute func
func Execute() {
	var rootCmd = &cobra.Command{Use: "moul", Short: "The minimalist photo collection generator"}
	rootCmd.AddCommand(newCmd)
	rootCmd.AddCommand(devCmd)
	rootCmd.AddCommand(buildCmd)
	rootCmd.AddCommand(igCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
