/*
hsvimage implements the image.Image interface with HSV + alpha images.
Most of this code was adapted from code in the Go standard library.
*/
package hsvimage

import (
	"github.com/spakin/hsvimage/hsvcolor"
	"image"
	"image/color"
)

// NHSVA is an in-memory image whose At method returns NHSVA values.
type NHSVA struct {
	// Pix holds the image's pixels, in H, S, V, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

func (p *NHSVA) ColorModel() color.Model { return hsvcolor.NHSVAModel }

func (p *NHSVA) Bounds() image.Rectangle { return p.Rect }

func (p *NHSVA) At(x, y int) color.Color {
	return p.NHSVAAt(x, y)
}

func (p *NHSVA) NHSVAAt(x, y int) hsvcolor.NHSVA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return hsvcolor.NHSVA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	return hsvcolor.NHSVA{s[0], s[1], s[2], s[3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *NHSVA) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

func (p *NHSVA) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := hsvcolor.NHSVAModel.Convert(c).(hsvcolor.NHSVA)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.H
	s[1] = c1.S
	s[2] = c1.V
	s[3] = c1.A
}
func (p *NHSVA) SetNHSVA(x, y int, c hsvcolor.NHSVA) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c.H
	s[1] = c.S
	s[2] = c.V
	s[3] = c.A
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *NHSVA) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &NHSVA{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &NHSVA{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *NHSVA) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0, i1 := 3, p.Rect.Dx()*4
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 4 {
			if p.Pix[i] != 0xff {
				return false
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}

// NewNHSVA returns a new NHSVA image with the given bounds.
func NewNHSVA(r image.Rectangle) *NHSVA {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 4*w*h)
	return &NHSVA{pix, 4 * w, r}
}
