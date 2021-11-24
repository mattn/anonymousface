package main

import (
	"embed"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	"os"
	"runtime"

	pigo "github.com/esimov/pigo/core"
	"github.com/nfnt/resize"
	"golang.org/x/image/draw"
)

const name = "anonymousface"

const version = "0.0.3"

var revision = "HEAD"

var (
	maskImg    image.Image
	classifier *pigo.Pigo

	//go:embed data/*
	static embed.FS
)

func main() {
	var showVersion bool
	flag.BoolVar(&showVersion, "V", false, "Print the version")
	flag.Parse()

	if showVersion {
		fmt.Printf("%s %s (rev: %s/%s)\n", name, version, revision, runtime.Version())
		return
	}

	f, err := static.Open("data/mask.png")
	if err != nil {
		log.Fatal("cannot open mask.png:", err)
	}
	defer f.Close()

	maskImg, _, err = image.Decode(f)
	if err != nil {
		log.Fatal("cannot decode mask.png:", err)
	}

	f, err = static.Open("data/facefinder")
	if err != nil {
		log.Fatal("cannot open facefinder:", err)
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		log.Fatal("cannot read facefinder:", err)
	}

	classifier, err = pigo.NewPigo().Unpack(b)
	if err != nil {
		log.Fatal("cannot unpack facefinder:", err)
	}

	img, _, err := image.Decode(os.Stdin)
	if err != nil {
		log.Fatal("cannot decode input image:", err)
		return
	}
	bounds := img.Bounds().Max
	param := pigo.CascadeParams{
		MinSize:     20,
		MaxSize:     2000,
		ShiftFactor: 0.1,
		ScaleFactor: 1.1,
		ImageParams: pigo.ImageParams{
			Pixels: pigo.RgbToGrayscale(pigo.ImgToNRGBA(img)),
			Rows:   bounds.Y,
			Cols:   bounds.X,
			Dim:    bounds.X,
		},
	}
	faces := classifier.RunCascade(param, 0)
	faces = classifier.ClusterDetections(faces, 0.18)

	canvas := image.NewRGBA(img.Bounds())
	draw.Draw(canvas, img.Bounds(), img, image.Point{0, 0}, draw.Over)
	for _, face := range faces {
		pt := image.Point{face.Col - face.Scale/2, face.Row - face.Scale/2}
		fimg := resize.Resize(uint(face.Scale), uint(face.Scale), maskImg, resize.NearestNeighbor)
		draw.Copy(canvas, pt, fimg, fimg.Bounds(), draw.Over, nil)
	}
	err = jpeg.Encode(os.Stdout, canvas, &jpeg.Options{Quality: 100})
	if err != nil {
		log.Fatal("cannot encode output image:", err)
		return
	}
}
