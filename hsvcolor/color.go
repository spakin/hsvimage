// hsvcolor provides HSV color models.
package hsvcolor

import (
	"image/color"
	"math"
)

// NHSVA represents a non-alpha-premultiplied 32-bit HSV color.  Note that all
// color channels range from 0 to 255.  (It is more common for hue to range
// from 0 to 359 and saturation and value to range from 0 to 1, but that's not
// what we do here.)
type NHSVA struct {
	H, S, V, A uint8
}

// min3uint32 returns the minimum of three uint32 values.
func min3uint32(a, b, c uint32) uint32 {
	m := a
	if b < m {
		m = b
	}
	if c < m {
		m = c
	}
	return m
}

// max3uint32 returns the maximum of three uint32 values.
func max3uint32(a, b, c uint32) uint32 {
	m := a
	if b > m {
		m = b
	}
	if c > m {
		m = c
	}
	return m
}

// nhsvaModel converts an arbitrary color to an NHSVA color.
func nhsvaModel(c color.Color) color.Color {
	// Handle the easy cases first: already NHSVA and fully transparent.
	if _, ok := c.(NHSVA); ok {
		return c
	}
	r, g, b, a := c.RGBA() // 32-bit values in the range [0, 65535]
	if a == 0 {
		return NHSVA{0, 0, 0, 0}
	}

	// Convert from premultiplied RGBA to non-premultiplied RGBA.
	r = (r * 65535) / a
	g = (g * 65535) / a
	b = (b * 65535) / a

	// Compute the easy channels: saturation and value.
	cMin := min3uint32(r, g, b)
	cMax := max3uint32(r, g, b)
	delta := cMax - cMin
	v := cMax
	var s uint32
	if cMax > 0 {
		s = (65535 * delta) / cMax
	}

	// Compute hue.
	if delta == 0 {
		return NHSVA{0, 0, uint8(v >> 8), uint8(a >> 8)} // Gray + alpha
	}
	var h360 int // Hue in the range [0, 359]
	ri, gi, bi, di := int(r), int(g), int(b), int(delta)
	switch cMax {
	case r:
		h360 = (60*(gi-bi))/di + 0
	case g:
		h360 = (60*(bi-ri))/di + 120
	case b:
		h360 = (60*(ri-gi))/di + 240
	}
	h360 = (h360 + 360) % 360             // Make positive.
	h := uint32((h360*65536 + 180) / 360) // Scale to [0, 65535].

	// Return an NHSVA color.
	return NHSVA{uint8(h >> 8), uint8(s >> 8), uint8(v >> 8), uint8(a >> 8)}
}

// An NHSVAModel is a model for NHSVA colors.
var NHSVAModel color.Model = color.ModelFunc(nhsvaModel)

// RGBA converts an NHSVA color to alpha-premultiplied RGBA.
func (c NHSVA) RGBA() (r, g, b, a uint32) {
	// Handle the easy case: a grayscale value.
	v16 := uint32(c.V) // 16-bit value in a 32-bit field
	v16 |= v16 << 8
	a16 := uint32(c.A) // 16-bit alpha in a 32-bit field
	a16 |= a16 << 8
	if c.S == 0 {
		v16pm := (v16*a16 + 32768) / 65535
		return v16pm, v16pm, v16pm, a16
	}

	// We work with float64 values primarily out of laziness: most of the
	// conversion formulas on the Web assume real values.
	hf := float64(c.H) / 255.0
	sf := float64(c.S) / 255.0
	vf := float64(c.V) / 255.0
	af := float64(c.A) / 255.0
	cf := vf * sf
	hf6 := hf / 60.0
	xf := cf * (1.0 - math.Abs(math.Mod(hf6, 2.0)-1.0))
	var rf, gf, bf float64
	switch {
	case hf6 < 0.0:
		panic("Internal error in RGBA (hf6 too small)")
	case hf6 <= 1.0:
		rf, gf, bf = cf, xf, 0.0
	case hf6 <= 2.0:
		rf, gf, bf = xf, cf, 0.0
	case hf6 <= 3.0:
		rf, gf, bf = 0.0, cf, xf
	case hf6 <= 4.0:
		rf, gf, bf = 0.0, xf, cf
	case hf6 <= 5.0:
		rf, gf, bf = xf, 0.0, cf
	case hf6 <= 6.0:
		rf, gf, bf = cf, 0.0, xf
	default:
		panic("Internal error in RGBA (hf6 too large)")
	}
	mf := vf - cf
	rf += mf
	gf += mf
	bf += mf

	// Premultiply by alpha then convert from float64 to uint32.
	r16 := uint32(rf * af * 65535.0)
	g16 := uint32(gf * af * 65535.0)
	b16 := uint32(bf * af * 65535.0)
	return r16, g16, b16, a16
}
