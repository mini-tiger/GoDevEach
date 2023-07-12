package main

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
)

func main() {
	a := app.New()
	w := a.NewWindow("Canvas")

	rect := canvas.NewRectangle(color.White)

	text := canvas.NewText("Hello Text", color.White)
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}

	line := canvas.NewLine(color.White)
	line.StrokeWidth = 5

	circle := canvas.NewCircle(color.White)
	circle.StrokeColor = color.Gray{0x99}
	circle.StrokeWidth = 5

	image := canvas.NewImageFromResource(theme.FyneLogo())
	image.FillMode = canvas.ImageFillOriginal

	gradient := canvas.NewHorizontalGradient(color.White, color.Transparent)

	container := fyne.NewContainerWithLayout(
		layout.NewGridWrapLayout(fyne.NewSize(150, 150)),
		rect, text, line, circle, image, gradient)
	w.SetContent(container)
	w.ShowAndRun()
}
