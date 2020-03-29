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
