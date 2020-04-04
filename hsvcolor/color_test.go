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

// Because color conversions with 8-bit color channels are inexact, we define a
// "close enough" metric.
func near(a uint8, b uint8) bool {
	diff := int(a) - int(b)
	if diff < 0 {
		diff = -diff
	}
	return diff < 4
}

// TestNHSVToNRGB confirms that we can convert non-premultiplied HSV to
// non-premultiplied RGB, with no transparency in either.
func TestNHSVToNRGB(t *testing.T) {
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

// TestNHSVAToNRGBA confirms that we can convert non-premultiplied HSV to
// premultiplied RGB, with transparency preserved.
func TestNHSVAToNRGBA(t *testing.T) {
	for ai := uint32(0); ai <= 255; ai += 15 {
		aOrig := uint8(ai)
		for _, cEq := range colorEquivalences {
			nhsva := NHSVA{cEq.HSV[0], cEq.HSV[1], cEq.HSV[2], aOrig}
			rp16, gp16, bp16, a16 := nhsva.RGBA()
			if a16 == 0 {
				// Special case for fully transparent colors.
				if rp16 != 0 || gp16 != 0 || bp16 != 0 || aOrig != 0 {
					t.Fatalf("Incorrectly mapped full-transparent %s from %v to [%d %d %d %d] (expected [0 0 0 0])", cEq.Name, nhsva, rp16, gp16, bp16, a16)
				}
				continue
			}
			var r, g, b uint8
			a := uint8(a16 >> 8)
			a16half := a16 / 2
			r = uint8((255*rp16 + a16half) / a16)
			g = uint8((255*gp16 + a16half) / a16)
			b = uint8((255*bp16 + a16half) / a16)
			if !near(r, cEq.RGB[0]) || !near(g, cEq.RGB[1]) || !near(b, cEq.RGB[2]) || a != aOrig {
				t.Fatalf("Incorrectly mapped %s from %v to [%d %d %d %d] (expected %v + %d)", cEq.Name, nhsva, r, g, b, a, cEq.RGB, aOrig)
			}
		}
	}
}

// TestGrayHSV64ToRGB confirms that we can convert 64-bit grayscale HSV values
// to RGB.
func TestGrayHSV64ToRGB(t *testing.T) {
	for vi := uint32(0); vi <= 65535; vi++ {
		v := uint16(vi)
		hsv := NHSVA64{0, 0, v, 65535}
		r32, g32, b32, a32 := hsv.RGBA()
		r, g, b, a := uint16(r32), uint16(g32), uint16(b32), uint16(a32)
		if r != v || g != v || b != v || a != 65535 {
			t.Fatalf("Incorrectly mapped %#v to {%d, %d, %d, %d}",
				hsv, r, g, b, a)
		}
	}
}

// TestGrayHSV64ToRGBA confirms that we can convert grayscale 64-bit HSV values
// to RGB in the context of partial transparency.
func TestGrayHSV64ToRGBA(t *testing.T) {
	for ai := uint32(0); ai <= 65535; ai += 3855 {
		a := uint16(ai)
		for vi := uint32(0); vi <= 65535; vi += 3855 {
			v := uint16(vi)
			if ai > 0 {
				// Round v so as to make the conversion exact.
				v = uint16((((vi*ai+32767)/65535)*65535 + ai/2) / ai)
			}
			hsv := NHSVA64{0, 0, v, a}
			rp32, gp32, bp32, a32 := hsv.RGBA() // Premultiplied colors
			var r, g, b uint16                  // Non-premultiplied
			if a32 != 0 {
				// Not fully transparent -- divide by alpha and
				// round.
				a32half := a32 / 2
				r = uint16((65535*rp32 + a32half) / a32)
				g = uint16((65535*gp32 + a32half) / a32)
				b = uint16((65535*bp32 + a32half) / a32)
			} else {
				// Fully transparent -- treat the value as 0.
				v = 0
			}
			if r != v || g != v || b != v || uint16(a32) != a {
				t.Fatalf("Incorrectly mapped %#v to {%d, %d, %d, %d}",
					hsv, r, g, b, a32)
			}
		}
	}
}

// TestGrayToHSV64 confirms that we can convert grayscale values to 64-bit HSV.
func TestGrayToHSV64(t *testing.T) {
	for vi := uint32(0); vi <= 65535; vi++ {
		g := color.Gray16{uint16(vi)}
		hsv := NHSVA64Model.Convert(g).(NHSVA64)
		if hsv.H != 0 || hsv.S != 0 || hsv.V != g.Y || hsv.A != 65535 {
			t.Fatalf("Incorrectly mapped %#v to %#v", g, hsv)
		}
	}
}

// An rgbHSV64assoc associates an RGB color with a 64-bit HSV color.
type rgbHSV64assoc struct {
	Name string
	RGB  [3]uint16
	HSV  [3]uint16
}

// colorEquivalences64 associates RGB with HSV values.  In this form, all color
// channels lie in [0, 65535].
var colorEquivalences64 = []rgbHSV64assoc{
	{"black", [3]uint16{0, 0, 0}, [3]uint16{0, 0, 0}},
	{"white", [3]uint16{65535, 65535, 65535}, [3]uint16{0, 0, 65535}},
	{"red", [3]uint16{65535, 0, 0}, [3]uint16{0, 65535, 65535}},
	{"green", [3]uint16{0, 65535, 0}, [3]uint16{21845, 65535, 65535}},
	{"blue", [3]uint16{0, 0, 65535}, [3]uint16{43690, 65535, 65535}},
	{"yellow", [3]uint16{65535, 65535, 0}, [3]uint16{10923, 65535, 65535}},
	{"cyan", [3]uint16{0, 65535, 65535}, [3]uint16{32768, 65535, 65535}},
	{"magenta", [3]uint16{65535, 0, 65535}, [3]uint16{54613, 65535, 65535}},
	{"dark blue", [3]uint16{0, 0, 32768}, [3]uint16{43690, 65535, 32768}},
	{"pale yellow", [3]uint16{65535, 65535, 49344}, [3]uint16{10923, 16191, 65535}},
}

// TestNRGBToNHSV64 confirms that we can convert non-premultiplied 64-bit RGB
// to non-premultiplied 64-bit HSV, with no transparency in either.
func TestNRGBToNHSV64(t *testing.T) {
	for _, cEq := range colorEquivalences64 {
		nrgba := color.NRGBA64{cEq.RGB[0], cEq.RGB[1], cEq.RGB[2], 65535}
		nhsva := NHSVA64Model.Convert(nrgba).(NHSVA64)
		if nhsva.H != cEq.HSV[0] || nhsva.S != cEq.HSV[1] || nhsva.V != cEq.HSV[2] || nhsva.A != 65535 {
			t.Fatalf("Incorrectly mapped %s from %v to %v (expected %v)", cEq.Name, nrgba, nhsva, cEq.HSV)
		}
	}
}

// Because color conversions with 16-bit color channels are inexact, we define a
// "close enough" metric.
func near16(a uint16, b uint16) bool {
	diff := int(a) - int(b)
	if diff < 0 {
		diff = -diff
	}
	return diff < 18
}

// TestNRGBAToNHSVA64 confirms that we can convert non-premultiplied 64-bit RGB
// to non-premultiplied 64-bit HSV, with transparency preserved.
func TestNRGBAToNHSVA64(t *testing.T) {
	for ai := uint32(0); ai <= 65535; ai += 3855 {
		a := uint16(ai)
		for _, cEq := range colorEquivalences64 {
			nrgba := color.NRGBA64{cEq.RGB[0], cEq.RGB[1], cEq.RGB[2], a}
			nhsva := NHSVA64Model.Convert(nrgba).(NHSVA64)
			if a == 0 {
				// Special case for fully transparent colors
				if nhsva.H != 0 || nhsva.S != 0 || nhsva.V != 0 || nhsva.A != 0 {
					t.Fatalf("Incorrectly mapped %s from %v to %v (expected [0, 0, 0, 0])", cEq.Name, nrgba, nhsva)
				}
				continue
			}
			if !near16(nhsva.H, cEq.HSV[0]) || !near16(nhsva.S, cEq.HSV[1]) || !near16(nhsva.V, cEq.HSV[2]) || nhsva.A != a {
				t.Fatalf("Incorrectly mapped %s from %v to %v (expected %v + %d)", cEq.Name, nrgba, nhsva, cEq.HSV, a)
			}
		}
	}
}

// TestNHSVToNRGB64 confirms that we can convert non-premultiplied 64-bit HSV
// to non-premultiplied 64-bit RGB, with no transparency in either.
func TestNHSVToNRGB64(t *testing.T) {
	for _, cEq := range colorEquivalences64 {
		nhsva := NHSVA64{cEq.HSV[0], cEq.HSV[1], cEq.HSV[2], 65535}
		r16, g16, b16, a16 := nhsva.RGBA() // Same as non-premultiplied because alpha is 65535.
		r := uint16(r16)
		g := uint16(g16)
		b := uint16(b16)
		a := uint16(a16)
		if !near16(r, cEq.RGB[0]) || !near16(g, cEq.RGB[1]) || !near16(b, cEq.RGB[2]) || a != 65535 {
			t.Fatalf("Incorrectly mapped %s from %v to [%d %d %d %d] (expected %v + 65535)", cEq.Name, nhsva, r, g, b, a, cEq.RGB)
		}
	}
}

// TestNHSVAToNRGBA64 confirms that we can convert non-premultiplied 64-bit HSV
// to premultiplied 64-bit RGB, with transparency preserved.
func TestNHSVAToNRGBA64(t *testing.T) {
	for ai := uint32(0); ai <= 65535; ai += 3855 {
		aOrig := uint16(ai)
		for _, cEq := range colorEquivalences64 {
			nhsva := NHSVA64{cEq.HSV[0], cEq.HSV[1], cEq.HSV[2], aOrig}
			rp16, gp16, bp16, a16 := nhsva.RGBA()
			if a16 == 0 {
				// Special case for fully transparent colors.
				if rp16 != 0 || gp16 != 0 || bp16 != 0 || aOrig != 0 {
					t.Fatalf("Incorrectly mapped full-transparent %s from %v to [%d %d %d %d] (expected [0 0 0 0])", cEq.Name, nhsva, rp16, gp16, bp16, a16)
				}
				continue
			}
			var r, g, b uint16
			a := uint16(a16)
			a16half := a16 / 2
			r = uint16((65535*rp16 + a16half) / a16)
			g = uint16((65535*gp16 + a16half) / a16)
			b = uint16((65535*bp16 + a16half) / a16)
			if !near16(r, cEq.RGB[0]) || !near16(g, cEq.RGB[1]) || !near16(b, cEq.RGB[2]) || a != aOrig {
				t.Fatalf("Incorrectly mapped %s from %v to [%d %d %d %d] (expected %v + %d)", cEq.Name, nhsva, r, g, b, a, cEq.RGB, aOrig)
			}
		}
	}
}
