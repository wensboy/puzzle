package palette

import (
	"fmt"
	"testing"
)

func Test_palette(t *testing.T) {
	p := NewPalette()
	fmt.Printf("bad key is %s\n", _default_palette.Put(RGB_DEFAULT, RGB_DEFAULT).Sprint("_system_"))
	fmt.Printf("bad key is %s\n", p.Put(RGB_SKYBLUE, RGB_BLACK).Sprint("_system_"))
}

func Test_basic_color(t *testing.T) {
	fmt.Printf("%s\n", Red("red message"))
	fmt.Printf("%s\n", Green("green message"))
	fmt.Printf("%s\n", Blue("blue message"))
	fmt.Printf("%s\n", White("white message"))
	fmt.Printf("%s\n", Black("black message"))
	fmt.Printf("%s\n", Yellow("yellow message"))
}
