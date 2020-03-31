// This file tests HSV color conversions.

package hsvcolor

import (
	"image/color"
	"testing"
)

// TestGrayHSVToRGB confirms that we can convert grayscale HSV values to RGB.
func TestGrayHSVToRGB(t *testing.T) {
	for vi := uint32(0); vi <= 255; vi++ {
		v := uint8(vi)
		hsv := NHSVA{0, 0, v, 255}
		r32, g32, b32, a32 := hsv.RGBA()
		r, g, b, a := uint8(r32>>8), uint8(g32>>8), uint8(b32>>8), uint8(a32>>8)
		if r != v || g != v || b != v || a != 255 {
			t.Fatalf("Incorrectly mapped %#v to {%d, %d, %d, %d}",
				hsv, r, g, b, a)
		}
	}
}

// TestGrayHSVToRGBA confirms that we can convert grayscale HSV values to RGB
// in the context of partial transparency.
func TestGrayHSVToRGBA(t *testing.T) {
	for ai := uint32(0); ai <= 255; ai += 15 {
		a := uint8(ai)
		for vi := uint32(0); vi <= 255; vi += 15 {
			v := uint8(vi)
			hsv := NHSVA{0, 0, v, a}
			rp32, gp32, bp32, a32 := hsv.RGBA() // Premultiplied colors
			var r, g, b uint8                   // Non-premultiplied
			if a32 != 0 {
				// Not fully transparent -- divide by alpha and
				// round.
				a32half := a32 / 2
				r = uint8((255*rp32 + a32half) / a32)
				g = uint8((255*gp32 + a32half) / a32)
				b = uint8((255*bp32 + a32half) / a32)
			} else {
				// Fully transparent -- treat the value as 0.
				v = 0
			}
			if r != v || g != v || b != v || uint8(a32>>8) != a {
				t.Fatalf("Incorrectly mapped %#v to {%d, %d, %d, %d}",
					hsv, r, g, b, a32>>8)
			}
		}
	}
}

// TestGrayToHSV confirms that we can convert grayscale values to HSV.
func TestGrayToHSV(t *testing.T) {
	for vi := uint32(0); vi <= 255; vi++ {
		g := color.Gray{uint8(vi)}
		hsv := NHSVAModel.Convert(g).(NHSVA)
		if hsv.H != 0 || hsv.S != 0 || hsv.V != g.Y || hsv.A != 255 {
			t.Fatalf("Incorrectly mapped %#v to %#v", g, hsv)
		}
	}
}

// An rgbHSVassoc associates an RGB color with an HSV color.
type rgbHSVassoc struct {
	Name string
	RGB  [3]uint8
	HSV  [3]uint8
}

// colorEquivalences associates RGB with HSV values.  In this form, all color
// channels lie in [0, 255].
var colorEquivalences = []rgbHSVassoc{
	{"black", [3]uint8{0, 0, 0}, [3]uint8{0, 0, 0}},
	{"white", [3]uint8{255, 255, 255}, [3]uint8{0, 0, 255}},
	{"red", [3]uint8{255, 0, 0}, [3]uint8{0, 255, 255}},
	{"green", [3]uint8{0, 255, 0}, [3]uint8{85, 255, 255}},
	{"blue", [3]uint8{0, 0, 255}, [3]uint8{170, 255, 255}},
	{"yellow", [3]uint8{255, 255, 0}, [3]uint8{43, 255, 255}},
	{"cyan", [3]uint8{0, 255, 255}, [3]uint8{128, 255, 255}},
	{"magenta", [3]uint8{255, 0, 255}, [3]uint8{213, 255, 255}},
	{"dark blue", [3]uint8{0, 0, 128}, [3]uint8{170, 255, 128}},
	{"pale yellow", [3]uint8{255, 255, 192}, [3]uint8{43, 63, 255}},
}

// TestNRGBToNHSV confirms that we can convert non-premultiplied RGB to
// non-premultiplied HSV, with no transparency in either.
func TestNRGBToNHSV(t *testing.T) {
	for _, cEq := range colorEquivalences {
		nrgba := color.NRGBA{cEq.RGB[0], cEq.RGB[1], cEq.RGB[2], 255}
		nhsva := NHSVAModel.Convert(nrgba).(NHSVA)
		if nhsva.H != cEq.HSV[0] || nhsva.S != cEq.HSV[1] || nhsva.V != cEq.HSV[2] || nhsva.A != 255 {
			t.Fatalf("Incorrectly mapped %s from %v to %v (expected %v)", cEq.Name, nrgba, nhsva, cEq.HSV)
		}
	}
}

// TestNRGBAToNHSVA confirms that we can convert non-premultiplied RGB to
// non-premultiplied HSV, with transparency preserved.
func TestNRGBAToNHSVA(t *testing.T) {
	for ai := uint32(0); ai <= 255; ai += 15 {
		a := uint8(ai)
		for _, cEq := range colorEquivalences {
			nrgba := color.NRGBA{cEq.RGB[0], cEq.RGB[1], cEq.RGB[2], a}
			nhsva := NHSVAModel.Convert(nrgba).(NHSVA)
			if a == 0 {
				// Special case for fully transparent colors
				if nhsva.H != 0 || nhsva.S != 0 || nhsva.V != 0 || nhsva.A != 0 {
					t.Fatalf("Incorrectly mapped %s from %v to %v (expected [0, 0, 0, 0])", cEq.Name, nrgba, nhsva)
				}
				continue
			}
			if nhsva.H != cEq.HSV[0] || nhsva.S != cEq.HSV[1] || nhsva.V != cEq.HSV[2] || nhsva.A != a {
				t.Fatalf("Incorrectly mapped %s from %v to %v (expected %v + %d)", cEq.Name, nrgba, nhsva, cEq.HSV, a)
			}
		}
	}
}

// TestNHSVToNRGB confirms that we can convert non-premultiplied HSV to
// non-premultiplied RGB, with no transparency in either.
func TestNHSVToNRGB(t *testing.T) {
	// Because HSV to RGB conversions are inexact, we define a
	// "close enough" metric.
	near := func(a uint8, b uint8) bool {
		diff := int(a) - int(b)
		if diff < 0 {
			diff = -diff
		}
		return diff < 4
	}

	// Test a selection of color conversions for being close enough.
	for _, cEq := range colorEquivalences {
		nhsva := NHSVA{cEq.HSV[0], cEq.HSV[1], cEq.HSV[2], 255}
		r16, g16, b16, a16 := nhsva.RGBA() // Same as non-premultiplied because alpha is 255.
		r := uint8(r16 >> 8)
		g := uint8(g16 >> 8)
		b := uint8(b16 >> 8)
		a := uint8(a16 >> 8)
		if !near(r, cEq.RGB[0]) || !near(g, cEq.RGB[1]) || !near(b, cEq.RGB[2]) || a != 255 {
			t.Fatalf("Incorrectly mapped %s from %v to [%d %d %d %d] (expected %v + 255)", cEq.Name, nhsva, r, g, b, a, cEq.RGB)
		}
	}
}
