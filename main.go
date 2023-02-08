package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/disintegration/gift"
	"github.com/nfnt/resize"
)

func main() {
	http.HandleFunc("/smol/", func(w http.ResponseWriter, r *http.Request) {
		// Extract the id and addons from the URL
		id := strings.TrimPrefix(r.URL.Path, "/smol/")
		id = strings.Split(id, ".")[0]
		addons := r.URL.Query().Get("addons")

		// Load the base image
		baseImage, err := LoadImage(fmt.Sprintf("/images/smols/%s.png", id))
		if err != nil {
			http.Error(w, "Error loading base image", http.StatusInternalServerError)
			return
		}

		// Load the addons
		var addonImages []image.Image
		var finalImage image.Image

		if addons != "" {
			addonList := strings.Split(addons, ",")
			for _, addon := range addonList {
				addonImage, err := LoadImage(fmt.Sprintf("/images/addons/%s.png", addon))
				if err != nil {
					http.Error(w, "Error loading addon image", http.StatusInternalServerError)
					return
				}
				addonImages = append(addonImages, addonImage)
			}

			// Compose the final image
			finalImage := ComposeImages(baseImage, addonImages)

		} else {
			finalImage = baseImage

		}
		// Write the final image to the response
		err = WriteImage(w, baseImage)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal(err)
	}
}

// LoadImage loads an image file from disk into memory
func LoadImage(filepath string) (image.Image, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, nil
}

// ComposeImages composes multiple images into a single image
func ComposeImages(baseImage image.Image, addonImages []image.Image) image.Image {
	g := gift.New(gift.Overlay(gift.Rectangle(baseImage.Bounds())))

	for _, addonImage := range addonImages {
		resizedAddonImage := resize.Resize(uint(baseImage.Bounds().Dx()), 0, addonImage, resize.Lanczos3)
		g = gift.New(gift.Overlay(gift.Rectangle(resizedAddonImage.Bounds())), g)
	}

	finalImage := image.NewNRGBA(g.Bounds(baseImage.Bounds()))
	g.Draw(finalImage, baseImage)

	return finalImage
}

func WriteImage(w http.ResponseWriter, img image.Image) error {
	w.Header().Set("Content-Type", "image/png")
	return png.Encode(w, img)
}
