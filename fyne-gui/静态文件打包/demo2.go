package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"image"
	"image/color"
	imageStatic "staticpackage/imageStatic"
)

func main() {
	a := app.New()
	w := a.NewWindow("Hello")

	img1 := canvas.NewImageFromResource(theme.FyneLogo())
	img1.FillMode = canvas.ImageFillOriginal

	//img2 := canvas.NewImageFromFile("./image/luffy.jpg")
	img2 := canvas.NewImageFromResource(imageStatic.Resourceluffy)
	img2.FillMode = canvas.ImageFillOriginal

	img4 := canvas.NewImageFromResource(imageStatic.Resource16627Jpg)
	img4.FillMode = canvas.ImageFillOriginal

	image := image.NewNRGBA(image.Rectangle{image.Point{0, 0}, image.Point{100, 100}})
	for i := 0; i < 100; i++ {
		for j := 0; j < 100; j++ {
			image.Set(i, j, color.NRGBA{uint8(i % 256), uint8(j % 256), 0, 255})
		}
	}
	img3 := canvas.NewImageFromImage(image)
	img3.FillMode = canvas.ImageFillOriginal

	container := &fyne.Container{
		Hidden:  false,
		Layout:  layout.NewGridWrapLayout(fyne.NewSize(150, 150)),
		Objects: []fyne.CanvasObject{img1, img2, img3, img4},
	}
	//container := fyne.NewContainerWithLayout(
	//	layout.NewGridWrapLayout(fyne.NewSize(150, 150)),
	//	img1, img2, img3)
	w.SetContent(container)
	w.ShowAndRun()
}
