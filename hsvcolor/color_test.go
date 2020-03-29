// This file tests HSV color conversions.

package hsvcolor

import (
	"testing"
)

// TestGraysToRGB confirms that we can convert grayscale HSV values to RGB.
func TestGraysToRGB(t *testing.T) {
	for v32 := uint32(0); v32 <= 255; v32++ {
		v := uint8(v32)
		hsv := NHSVA{0, 0, v, 255}
		r32, g32, b32, a32 := hsv.RGBA()
		r, g, b, a := uint8(r32>>8), uint8(g32>>8), uint8(b32>>8), uint8(a32>>8)
		if r != v || g != v || b != v || a != 255 {
			t.Fatalf("%#v incorrectly mapped to {%d, %d, %d, %d}",
				hsv, r, g, b, a)
		}
	}
}
