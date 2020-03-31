// This file tests HSV images.

package hsvimage

import (
	"github.com/spakin/hsvimage/hsvcolor"
	"image"
	"image/color"
	"testing"
)

// cmp was copied verbatim from the Go standard library's image_test.go file.
func cmp(cm color.Model, c0, c1 color.Color) bool {
	r0, g0, b0, a0 := cm.Convert(c0).RGBA()
	r1, g1, b1, a1 := cm.Convert(c1).RGBA()
	return r0 == r1 && g0 == g1 && b0 == b1 && a0 == a1
}

// TestImage was copied almost literally from the Go standard library's
// image_test.go file.
func TestImage(t *testing.T) {
	m := NewNHSVA(image.Rect(0, 0, 10, 10))
	if !image.Rect(0, 0, 10, 10).Eq(m.Bounds()) {
		t.Fatalf("%T: want bounds %v, got %v", m, image.Rect(0, 0, 10, 10), m.Bounds())
	}
	if !cmp(m.ColorModel(), image.Transparent, m.At(6, 3)) {
		t.Fatalf("%T: at (6, 3), want a zero color, got %v", m, m.At(6, 3))
	}
	m.Set(6, 3, image.Opaque)
	if !cmp(m.ColorModel(), image.Opaque, m.At(6, 3)) {
		t.Fatalf("%T: at (6, 3), want a non-zero color, got %v", m, m.At(6, 3))
	}
	if !m.SubImage(image.Rect(6, 3, 7, 4)).(*NHSVA).Opaque() {
		t.Fatalf("%T: at (6, 3) was not opaque", m)
	}
	m = m.SubImage(image.Rect(3, 2, 9, 8)).(*NHSVA)
	if !image.Rect(3, 2, 9, 8).Eq(m.Bounds()) {
		t.Fatalf("%T: sub-image want bounds %v, got %v", m, image.Rect(3, 2, 9, 8), m.Bounds())
	}
	if !cmp(m.ColorModel(), image.Opaque, m.At(6, 3)) {
		t.Fatalf("%T: sub-image at (6, 3), want a non-zero color, got %v", m, m.At(6, 3))
	}
	if !cmp(m.ColorModel(), image.Transparent, m.At(3, 3)) {
		t.Fatalf("%T: sub-image at (3, 3), want a zero color, got %v", m, m.At(3, 3))
	}
	m.Set(3, 3, image.Opaque)
	if !cmp(m.ColorModel(), image.Opaque, m.At(3, 3)) {
		t.Fatalf("%T: sub-image at (3, 3), want a non-zero color, got %v", m, m.At(3, 3))
	}
	// Test that taking an empty sub-image starting at a corner does not panic.
	m.SubImage(image.Rect(0, 0, 0, 0))
	m.SubImage(image.Rect(10, 0, 10, 0))
	m.SubImage(image.Rect(0, 10, 0, 10))
	m.SubImage(image.Rect(10, 10, 10, 10))
}

// TestSimpleColors checks that we can create an NHSVA image with simple,
// easily convertible colors and read the pixels back as RGBA.
func TestSimpleColors(t *testing.T) {
	// Define the set of colors to use in the image.
	hsvColors := []hsvcolor.NHSVA{
		{H: 0, S: 0, V: 0, A: 255},       // Black
		{H: 0, S: 0, V: 255, A: 255},     // White
		{H: 0, S: 255, V: 255, A: 255},   // Red
		{H: 85, S: 255, V: 255, A: 255},  // Green
		{H: 170, S: 255, V: 255, A: 255}, // Blue
		{H: 43, S: 255, V: 255, A: 255},  // Yellow
		{H: 0, S: 255, V: 255, A: 128},   // Half-transparent red
		{H: 85, S: 64, V: 255, A: 255},   // Pale green
		{H: 170, S: 255, V: 64, A: 255},  // Dark blue
		{H: 205, S: 82, V: 143, A: 255},  // French lilac
	}
	rgbColors := []color.RGBA{
		{R: 0, G: 0, B: 0, A: 255},       // Black
		{R: 255, G: 255, B: 255, A: 255}, // White
		{R: 255, G: 0, B: 0, A: 255},     // Red
		{R: 0, G: 255, B: 0, A: 255},     // Green
		{R: 0, G: 0, B: 255, A: 255},     // Blue
		{R: 252, G: 255, B: 0, A: 255},   // Yellow (with rounding error)
		{R: 128, G: 0, B: 0, A: 128},     // Half-transparent red
		{R: 191, G: 255, B: 191, A: 255}, // Pale green
		{R: 0, G: 0, B: 64, A: 255},      // Dark blue
		{R: 135, G: 97, B: 143, A: 255},  // French lilac (with rounding error)
	}
	nc := len(hsvColors)

	// Draw an image with NHSVA colors.
	const wd = 100
	const ht = 100
	img := NewNHSVA(image.Rect(0, 0, wd, ht))
	i := 0
	for y := 0; y < ht; y++ {
		for x := 0; x < wd; x++ {
			cwHSV := hsvColors[i%nc]
			img.Set(x, y, cwHSV)
			i++
		}
	}

	// Check that the RGBA colors we read are as expected.
	i = 0
	for y := 0; y < ht; y++ {
		for x := 0; x < wd; x++ {
			cwHSV := hsvColors[i%nc]
			crHSV := img.NHSVAAt(x, y)
			if crHSV != cwHSV {
				t.Fatalf("Wrote %v but read %v at (%d, %d)", cwHSV, crHSV, x, y)
			}
			r, g, b, a := img.At(x, y).RGBA()
			crRGB := color.RGBA{uint8(r >> 8), uint8(g >> 8), uint8(b >> 8), uint8(a >> 8)}
			if crRGB != rgbColors[i%nc] {
				t.Fatalf("Expected %v but saw %v at (%d, %d)", rgbColors[i%nc], crRGB, x, y)
			}
			i++
		}
	}
}
