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
	{"green", [3]uint8{0, 255, 0}, [3]uint8{(120 * 256) / 360, 255, 255}},
	{"blue", [3]uint8{0, 0, 255}, [3]uint8{(240 * 256) / 360, 255, 255}},
	{"yellow", [3]uint8{255, 255, 0}, [3]uint8{(60 * 256) / 360, 255, 255}},
	{"cyan", [3]uint8{0, 255, 255}, [3]uint8{(180 * 256) / 360, 255, 255}},
	{"magenta", [3]uint8{255, 0, 255}, [3]uint8{(300 * 256) / 360, 255, 255}},
	{"dark blue", [3]uint8{0, 0, 128}, [3]uint8{(240 * 256) / 360, 255, 128}},
	{"pale yellow", [3]uint8{255, 255, 192}, [3]uint8{(60 * 256) / 360, 63, 255}},
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
