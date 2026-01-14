// Package palette color correction for content.
package palette

import (
	"github.com/fatih/color"
)

var (
	_default_palette *Palette

	_default_fg_rgb = RGB_DEFAULT
	_default_fg_hex = HEX_DEFAULT
	_default_bg_rgb = RGB_DEFAULT
	_default_bg_hex = HEX_DEFAULT
)

type (
	// Palette record color usage and mixing.
	Palette struct {
		FgRgb RGB
		BgRgb RGB
		FgHex string
		BgHex string
		c     *color.Color
	}
)

func init() {
	_default_palette = NewPalette()
}

func NewPalette() *Palette {
	p := &Palette{
		FgRgb: _default_fg_rgb,
		FgHex: _default_fg_hex,
		BgRgb: _default_bg_rgb,
		BgHex: _default_bg_hex,
	}
	return p
}

func Red(args ...any) string {
	return _default_palette.Put(RGB_RED, RGB_DEFAULT).Sprint(args...)
}

func Green(args ...any) string {
	return _default_palette.Put(RGB_GREEN, RGB_DEFAULT).Sprint(args...)
}

func Blue(args ...any) string {
	return _default_palette.Put(RGB_BLUE, RGB_DEFAULT).Sprint(args...)
}

func SkyBlue(args ...any) string {
	return _default_palette.Put(RGB_SKYBLUE, RGB_DEFAULT).Sprint(args...)
}

func Yellow(args ...any) string {
	return _default_palette.Put(RGB_YELLOW, RGB_DEFAULT).Sprint(args...)
}

func White(args ...any) string {
	return _default_palette.Put(RGB_WHITE, RGB_DEFAULT).Sprint(args...)
}

func Black(args ...any) string {
	return _default_palette.Put(RGB_BLACK, RGB_DEFAULT).Sprint(args...)
}

func Reset() {
	_default_palette.Reset()
}

func Put(fg, bg RGB) *color.Color {
	return _default_palette.Put(fg, bg)
}

func (p *Palette) mixed() {
	p.c = color.New()
	if p.FgRgb != RGB_DEFAULT {
		p.c.AddRGB(int(p.FgRgb.R), int(p.FgRgb.G), int(p.FgRgb.B))
	}
	if p.BgRgb != RGB_DEFAULT {
		p.c.AddBgRGB(int(p.BgRgb.R), int(p.BgRgb.G), int(p.BgRgb.B))
	}
}

func (p *Palette) Reset() {
	p.FgRgb = _default_fg_rgb
	p.FgHex = _default_fg_hex
	p.BgRgb = _default_bg_rgb
	p.BgHex = _default_bg_hex
	p.mixed()
}

func (p *Palette) Put(fg, bg RGB) *color.Color {
	p.FgRgb = fg
	p.BgRgb = bg
	p.mixed()
	return p.c
}
