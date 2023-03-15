package main

import (
	"build_dir/main/submain"
	"build_dir/public"
	"fmt"
)

func main() {
	fmt.Println(public.Util1())
	fmt.Println(submain.Sub())
}
