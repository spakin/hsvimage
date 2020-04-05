/*
Package hsvimage implements the image.Image interface with HSV + alpha images.
The code was largely adapted from code in the Go standard library.
*/
package hsvimage

import (
	"github.com/spakin/hsvimage/hsvcolor"
	"image"
	"image/color"
)

// NHSVA is an in-memory image whose At method returns hsvcolor.NHSVA values.
type NHSVA struct {
	// Pix holds the image's pixels, in H, S, V, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel states that an NHSVA image uses the hsvcolor.NHSVA color model.
func (p *NHSVA) ColorModel() color.Model { return hsvcolor.NHSVAModel }

// Bounds returns the image's bounding rectangle.
func (p *NHSVA) Bounds() image.Rectangle { return p.Rect }

// At returns the color at the given image coordinates.
func (p *NHSVA) At(x, y int) color.Color {
	return p.NHSVAAt(x, y)
}

// NHSVAAt returns the color at the given image coordinates as specifically an
// hsvcolor.NHSVA color.
func (p *NHSVA) NHSVAAt(x, y int) hsvcolor.NHSVA {
	if !(image.Point{x, y}.In(p.Rect)) {
		return hsvcolor.NHSVA{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	return hsvcolor.NHSVA{H: s[0], S: s[1], V: s[2], A: s[3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *NHSVA) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

// Set assigns an arbitrary color to a given coordinate.
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

// SetNHSVA assigns an NHSVA color to a given coordinate.
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

// NHSVA64 is an in-memory image whose At method returns hsvcolor.NHSVA64 values.
type NHSVA64 struct {
	// Pix holds the image's pixels, in H, S, V, A order and big-endian
	// format.  The pixel at (x, y) starts at Pix[(y-Rect.Min.Y)*Stride +
	// (x-Rect.Min.X)*8].
	Pix []uint8
	// Stride is the Pix stride (in bytes) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel states that an NHSVA64 image uses the hsvcolor.NHSVA64 color model.
func (p *NHSVA64) ColorModel() color.Model { return hsvcolor.NHSVA64Model }

// Bounds returns the image's bounding rectangle.
func (p *NHSVA64) Bounds() image.Rectangle { return p.Rect }

// At returns the color at the given image coordinates.
func (p *NHSVA64) At(x, y int) color.Color {
	return p.NHSVA64At(x, y)
}

// NHSVA64At returns the color at the given image coordinates as specifically an
// hsvcolor.NHSVA64 color.
func (p *NHSVA64) NHSVA64At(x, y int) hsvcolor.NHSVA64 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return hsvcolor.NHSVA64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+8 : i+8] // Small cap improves performance, see https://golang.org/issue/27857
	return hsvcolor.NHSVA64{
		H: uint16(s[0])<<8 | uint16(s[1]),
		S: uint16(s[2])<<8 | uint16(s[3]),
		V: uint16(s[4])<<8 | uint16(s[5]),
		A: uint16(s[6])<<8 | uint16(s[7]),
	}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *NHSVA64) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*8
}

// Set assigns an arbitrary color to a given coordinate.
func (p *NHSVA64) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := hsvcolor.NHSVA64Model.Convert(c).(hsvcolor.NHSVA64)
	s := p.Pix[i : i+8 : i+8] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8(c1.H >> 8)
	s[1] = uint8(c1.H)
	s[2] = uint8(c1.S >> 8)
	s[3] = uint8(c1.S)
	s[4] = uint8(c1.V >> 8)
	s[5] = uint8(c1.V)
	s[6] = uint8(c1.A >> 8)
	s[7] = uint8(c1.A)
}

// SetNHSVA64 assigns an NHSVA64 color to a given coordinate.
func (p *NHSVA64) SetNHSVA64(x, y int, c hsvcolor.NHSVA64) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+8 : i+8] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = uint8(c.H >> 8)
	s[1] = uint8(c.H)
	s[2] = uint8(c.S >> 8)
	s[3] = uint8(c.S)
	s[4] = uint8(c.V >> 8)
	s[5] = uint8(c.V)
	s[6] = uint8(c.A >> 8)
	s[7] = uint8(c.A)
}

// SubImage returns an image representing the portion of the image p visible
// through r. The returned value shares pixels with the original image.
func (p *NHSVA64) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &NHSVA64{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &NHSVA64{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *NHSVA64) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0, i1 := 6, p.Rect.Dx()*8
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 8 {
			if p.Pix[i+0] != 0xff || p.Pix[i+1] != 0xff {
				return false
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}

// NewNHSVA64 returns a new NHSVA64 image with the given bounds.
func NewNHSVA64(r image.Rectangle) *NHSVA64 {
	w, h := r.Dx(), r.Dy()
	pix := make([]uint8, 8*w*h)
	return &NHSVA64{pix, 8 * w, r}
}

// NHSVAF64 is an in-memory image whose At method returns hsvcolor.NHSVAF64
// values.
type NHSVAF64 struct {
	// Pix holds the image's pixels, in H, S, V, A order. The pixel at
	// (x, y) starts at Pix[(y-Rect.Min.Y)*Stride + (x-Rect.Min.X)*4].
	Pix []float64
	// Stride is the Pix stride (in 64-bit words) between vertically adjacent pixels.
	Stride int
	// Rect is the image's bounds.
	Rect image.Rectangle
}

// ColorModel states that an NHSVAF64 image uses the hsvcolor.NHSVAF64 color
// model.
func (p *NHSVAF64) ColorModel() color.Model { return hsvcolor.NHSVAF64Model }

// Bounds returns the image's bounding rectangle.
func (p *NHSVAF64) Bounds() image.Rectangle { return p.Rect }

// At returns the color at the given image coordinates.
func (p *NHSVAF64) At(x, y int) color.Color {
	return p.NHSVAF64At(x, y)
}

// NHSVAF64At returns the color at the given image coordinates as specifically
// an hsvcolor.NHSVAF64 color.
func (p *NHSVAF64) NHSVAF64At(x, y int) hsvcolor.NHSVAF64 {
	if !(image.Point{x, y}.In(p.Rect)) {
		return hsvcolor.NHSVAF64{}
	}
	i := p.PixOffset(x, y)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	return hsvcolor.NHSVAF64{H: s[0], S: s[1], V: s[2], A: s[3]}
}

// PixOffset returns the index of the first element of Pix that corresponds to
// the pixel at (x, y).
func (p *NHSVAF64) PixOffset(x, y int) int {
	return (y-p.Rect.Min.Y)*p.Stride + (x-p.Rect.Min.X)*4
}

// Set assigns an arbitrary color to a given coordinate.
func (p *NHSVAF64) Set(x, y int, c color.Color) {
	if !(image.Point{x, y}.In(p.Rect)) {
		return
	}
	i := p.PixOffset(x, y)
	c1 := hsvcolor.NHSVAF64Model.Convert(c).(hsvcolor.NHSVAF64)
	s := p.Pix[i : i+4 : i+4] // Small cap improves performance, see https://golang.org/issue/27857
	s[0] = c1.H
	s[1] = c1.S
	s[2] = c1.V
	s[3] = c1.A
}

// SetNHSVAF64 assigns an NHSVAF64 color to a given coordinate.
func (p *NHSVAF64) SetNHSVAF64(x, y int, c hsvcolor.NHSVAF64) {
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
func (p *NHSVAF64) SubImage(r image.Rectangle) image.Image {
	r = r.Intersect(p.Rect)
	// If r1 and r2 are Rectangles, r1.Intersect(r2) is not guaranteed to
	// be inside either r1 or r2 if the intersection is empty. Without
	// explicitly checking for this, the Pix[i:] expression below can
	// panic.
	if r.Empty() {
		return &NHSVAF64{}
	}
	i := p.PixOffset(r.Min.X, r.Min.Y)
	return &NHSVAF64{
		Pix:    p.Pix[i:],
		Stride: p.Stride,
		Rect:   r,
	}
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (p *NHSVAF64) Opaque() bool {
	if p.Rect.Empty() {
		return true
	}
	i0, i1 := 3, p.Rect.Dx()*4
	for y := p.Rect.Min.Y; y < p.Rect.Max.Y; y++ {
		for i := i0; i < i1; i += 4 {
			if p.Pix[i] != 1.0 {
				return false
			}
		}
		i0 += p.Stride
		i1 += p.Stride
	}
	return true
}

// NewNHSVAF64 returns a new NHSVAF64 image with the given bounds.
func NewNHSVAF64(r image.Rectangle) *NHSVAF64 {
	w, h := r.Dx(), r.Dy()
	pix := make([]float64, 4*w*h)
	return &NHSVAF64{pix, 4 * w, r}
}
