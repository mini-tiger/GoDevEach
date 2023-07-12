package image

/**
 * @Author: Tao Jun
 * @Since: 2023/7/12
 * @Desc: resource.go
**/

import (
	_ "embed"
	"fyne.io/fyne/v2"
)

//go:embed 16627.jpg
var jpg16627 []byte

//go:embed luffy.jpg
var luffy []byte

var Resource16627Jpg = &fyne.StaticResource{
	StaticName:    "16627.jpg",
	StaticContent: jpg16627}

var Resourceluffy = &fyne.StaticResource{
	StaticName:    "luffy.jpg",
	StaticContent: luffy}
