// Package hsvcolor provides HSV color models.
package hsvcolor

import (
	"image/color"
	"math"
)

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

// nhsvaFloat64ToRGBA is a helper function for NHSVA.RGBA and NHSVA64.RGBA that
// converts float64 versions of H, S, V, and A to RGBA.
func nhsvaFloat64ToRGBA(hf, sf, vf, af float64) (r uint32, g uint32, b uint32, a uint32) {
	// Follow the textbook formulas for converting HSV to RGB.
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
	a16 := uint32(af * 65535.0)
	return r16, g16, b16, a16
}

// NHSVA represents a non-alpha-premultiplied 32-bit HSV color.  Note that all
// color channels range from 0 to 255.  (It is more common for hue to range
// from 0 to 359 and saturation and value to range from 0 to 1, but that's not
// what we do here.)
type NHSVA struct {
	H, S, V, A uint8
}

// nhsvaModel converts an arbitrary color to an NHSVA color.
func nhsvaModel(c color.Color) color.Color {
	// Handle the easy case first: already NHSVA.
	if _, ok := c.(NHSVA); ok {
		return c
	}

	// Produce a 64-bit color then scale it down to 32 bits.
	nhsva64 := nhsva64Model(c).(NHSVA64)
	scale := func(n16 uint16) uint8 {
		return uint8((uint32(n16)*255 + 32768) / 65535)
	}
	return NHSVA{
		H: scale(nhsva64.H),
		S: scale(nhsva64.S),
		V: scale(nhsva64.V),
		A: scale(nhsva64.A),
	}
}

// NHSVAModel is a color model for NHSVA (non-alpha-premultiplied hue,
// saturation, and value plus alpha) colors.
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
	hf := float64(c.H) * 360.0 / 255.0
	sf := float64(c.S) / 255.0
	vf := float64(c.V) / 255.0
	af := float64(c.A) / 255.0
	return nhsvaFloat64ToRGBA(hf, sf, vf, af)
}

// NHSVA64 represents a non-alpha-premultiplied 64-bit HSV color.  Note that
// all color channels range from 0 to 65535.  (It is more common for hue to
// range from 0 to 359 and saturation and value to range from 0 to 1, but
// that's not what we do here.)
type NHSVA64 struct {
	H, S, V, A uint16
}

// nhsva64Model converts an arbitrary color to an NHSVA64 color.
func nhsva64Model(c color.Color) color.Color {
	// Handle the easy cases first: already NHSVA64 and fully transparent.
	if _, ok := c.(NHSVA64); ok {
		return c
	}
	r, g, b, a := c.RGBA() // 32-bit values in the range [0, 65535]
	if a == 0 {
		return NHSVA64{0, 0, 0, 0}
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
		return NHSVA64{0, 0, uint16(v), uint16(a)} // Gray + alpha
	}
	var h360 int // Hue in the range [0, 360]
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
	h := uint32((h360*65535 + 180) / 360) // Scale to [0, 65535].

	// Return an NHSVA color.
	return NHSVA64{uint16(h), uint16(s), uint16(v), uint16(a)}
}

// NHSVA64Model is a color model for NHSVA64 (non-alpha-premultiplied hue,
// saturation, and value plus alpha) colors.
var NHSVA64Model color.Model = color.ModelFunc(nhsva64Model)

// RGBA converts an NHSVA64 color to alpha-premultiplied RGBA.
func (c NHSVA64) RGBA() (r, g, b, a uint32) {
	// Handle the easy case: a grayscale value.
	a16 := uint32(c.A)
	if c.S == 0 {
		v16pm := (uint32(c.V)*a16 + 32768) / 65535
		return v16pm, v16pm, v16pm, a16
	}

	// We work with float64 values primarily out of laziness: most of the
	// conversion formulas on the Web assume real values.
	hf := float64(c.H) * 360.0 / 65535.0
	sf := float64(c.S) / 65535.0
	vf := float64(c.V) / 65535.0
	af := float64(c.A) / 65535.0
	return nhsvaFloat64ToRGBA(hf, sf, vf, af)
}

// NHSVAF64 represents a non-alpha-premultiplied HSV color with each channel
// represented by a 64-bit floating-point number.  In this representation, hue
// is a value in [0, 360); and the remaining channels are values in [0, 1].
type NHSVAF64 struct {
	H, S, V, A float64
}

// nhsvaF64Model converts an arbitrary color to an NHSVAF64 color.
func nhsvaF64Model(c color.Color) color.Color {
	// Handle the easy cases first: already NHSVAF64 and fully transparent.
	if _, ok := c.(NHSVAF64); ok {
		return c
	}
	r, g, b, a := c.RGBA() // 32-bit values in the range [0, 65535]
	if a == 0 {
		return NHSVAF64{0.0, 0.0, 0.0, 0.0}
	}

	// Convert all values to floating point.
	rf := float64(r) / 65535.0
	gf := float64(g) / 65535.0
	bf := float64(b) / 65535.0
	af := float64(a) / 65535.0

	// Convert from premultiplied RGBA to non-premultiplied RGBA.
	rf /= af
	gf /= af
	bf /= af

	// Compute the easy channels: saturation and value.
	cMin := math.Min(math.Min(rf, gf), bf)
	cMax := math.Max(math.Max(rf, gf), bf)
	delta := cMax - cMin
	vf := cMax
	var sf float64
	if cMax > 0.00 {
		sf = delta / cMax
	}

	// Compute hue.
	if delta == 0.0 {
		return NHSVAF64{0.0, 0.0, vf, af} // Gray + alpha
	}
	var hf float64
	switch cMax {
	case rf:
		hf = (gf-bf)/delta + 0.0
	case gf:
		hf = (bf-rf)/delta + 2.0
	case bf:
		hf = (rf-gf)/delta + 4.0
	}
	hf = math.Mod(hf*60.0+360.0, 360.0)

	// Return an NHSVAF64 color.
	return NHSVAF64{hf, sf, vf, af}
}

// NHSVAF64Model is a color model for NHSVAF64 (non-alpha-premultiplied hue,
// saturation, and value plus alpha, with 64-bit floating-point channels)
// colors.
var NHSVAF64Model color.Model = color.ModelFunc(nhsvaF64Model)

// RGBA converts an NHSVAF64 color to alpha-premultiplied RGBA.
func (c NHSVAF64) RGBA() (r, g, b, a uint32) {
	// Force all HSVA values into their expected range: [0, 360) for hue
	// (with wraparound) and [0, 1] for everything else (with clamping).
	clamp01 := func(x float64) float64 { return math.Max(0.0, math.Min(1.0, x)) }
	wrap360 := func(x float64) float64 { return math.Mod(math.Mod(x, 360.0)+360.0, 360.0) }
	hf := wrap360(c.H)
	sf := clamp01(c.S)
	vf := clamp01(c.V)
	af := clamp01(c.A)

	// Handle the easy case: a grayscale value.
	if sf == 0.0 {
		v16pm := uint32(vf * af * 65535.0)
		return v16pm, v16pm, v16pm, uint32(af * 65535.0)
	}

	// Handle all other cases.
	return nhsvaFloat64ToRGBA(hf, sf, vf, af)
}
