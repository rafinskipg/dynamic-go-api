package handler

import (
	"image"
	"image/png"
	"log"
	"net/http"
	"strings"

	"embed"

	"github.com/disintegration/gift"
	"github.com/nfnt/resize"
)

// content holds our static web server content.
//
//go:embed images/*
var content embed.FS

func Handler(w http.ResponseWriter, r *http.Request) {
	// Extract the id and addons from the URL
	id := r.URL.Query().Get("id")
	addons := r.URL.Query().Get("addons")

	log.Printf("Generating image: id: %s, addons: %s", id, addons)

	// packagePath, _ := filepath.Abs("./api/images/smols/")
	// imagePath := filepath.Join(packagePath, fmt.Sprintf("%s.png", id))

	// Load the base image
	baseImage, err := LoadImage("images/smols/" + id + ".png")
	if err != nil {
		log.Printf("Error loading base image: %s", err)
		http.Error(w, "Error loading base image", http.StatusInternalServerError)
		return
	}

	// Load the addons
	var addonImages []image.Image
	var finalImage image.Image

	if addons != "" {
		addonList := strings.Split(addons, ",")
		//addonPackagePath, _ := filepath.Abs("./images/addons/")
		for _, addon := range addonList {
			// addonImagePath := filepath.Join(addonPackagePath, fmt.Sprintf("%s.png", addon))
			addonImagePath := "images/addons/" + addon + ".png"
			addonImage, err := LoadImage(addonImagePath)
			if err != nil {
				http.Error(w, "Error loading addon image", http.StatusInternalServerError)
				return
			}
			addonImages = append(addonImages, addonImage)
		}

		// Compose the final image
		finalImage = ComposeImages(baseImage, addonImages)

	} else {
		finalImage = baseImage

	}
	// Write the final image to the response
	err = WriteImage(w, finalImage)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}

// LoadImage loads an image file from disk into memory
func LoadImage(filepath string) (image.Image, error) {
	// f, err := os.Open(filepath)
	f, err := content.Open(filepath)

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
	dstImage := image.NewNRGBA(baseImage.Bounds())
	g := gift.New()

	// Draw the base image
	g.Draw(dstImage, baseImage)

	for _, addonImage := range addonImages {
		resizedAddonImage := resize.Resize(uint(baseImage.Bounds().Dx()), 0, addonImage, resize.Lanczos3)
		g.DrawAt(dstImage, resizedAddonImage, image.Pt(0, 0), gift.OverOperator)
	}

	return dstImage
}

func WriteImage(w http.ResponseWriter, img image.Image) error {
	w.Header().Set("Content-Type", "image/png")
	return png.Encode(w, img)
}
