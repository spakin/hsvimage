hsvimage and hsvcolor
=====================

[![Go Report Card](https://goreportcard.com/badge/github.com/spakin/hsvimage)](https://goreportcard.com/report/github.com/spakin/hsvimage)
[![GoDoc](https://godoc.org/github.com/spakin/hsvimage?status.svg)](https://godoc.org/github.com/spakin/hsvimage)

The [Go programming language](https://golang.org/)'s standard library provides support for manipulating graphical images represented using [RGB](https://en.wikipedia.org/wiki/RGB_color_model), [CMYK](https://en.wikipedia.org/wiki/CMYK_color_model), [YCbCr](https://en.wikipedia.org/wiki/YCbCr), and [grayscale](https://en.wikipedia.org/wiki/Grayscale) color models and variations such as 8- vs. 16-bit channels, premultiplied alpha vs. non-premultiplied alpha vs. no alpha, and [paletted color](https://en.wikipedia.org/wiki/Indexed_color).  `hsvimage` augments the Go standard library with support for the [HSV](https://en.wikipedia.org/wiki/HSL_and_HSV) (hue, saturation, and value) color model.

Usage
-----

`hsvimage` provides the following new image types (all of which implement the [`image.Image` interface](https://golang.org/pkg/image/#Image)), color types (all of which implement the [`color.Color` interface](https://golang.org/pkg/image/color/#Color), and color models (all of which implement the [`color.Model` interface](https://golang.org/pkg/image/color/#Model)):

| Image | Color | Color model | Description |
| :---- | :---- | :---------- | :---------- |
| [`NHSVA`](https://godoc.org/github.com/spakin/hsvimage#NHSVA) | [`NHSVA`](https://godoc.org/github.com/spakin/hsvimage/hsvcolor#NHSVA) | [`NHSVAModel`](https://godoc.org/github.com/spakin/hsvimage/hsvcolor#pkg-variables) | Non-alpha-premultiplied HSV + alpha, 8-bit color channels |
| [`NHSVA64`](https://godoc.org/github.com/spakin/hsvimage#NHSVA64) | [`NHSVA64`](https://godoc.org/github.com/spakin/hsvimage/hsvcolor#NHSVA64) | [`NHSVA64Model`](https://godoc.org/github.com/spakin/hsvimage/hsvcolor#pkg-variables) | Non-alpha-premultiplied HSV + alpha, 16-bit color channels |
| [`NHSVAF64`](https://godoc.org/github.com/spakin/hsvimage#NHSVAF64) | [`NHSVAF64`](https://godoc.org/github.com/spakin/hsvimage/hsvcolor#NHSVAF64) | [`NHSVAF64Model`](https://godoc.org/github.com/spakin/hsvimage/hsvcolor#pkg-variables) | Non-alpha-premultiplied HSV + alpha, 64-bit floating-point color channels |


`hsvimage` and `hsvimage/hsvcolor`, which are analogous to Go's [`image`](https://golang.org/pkg/image/) and [`image/color`](https://golang.org/pkg/image/color/), respectively, can be imported in the usual manner:

```Go
import (
	"github.com/spakin/hsvimage"
	"github.com/spakin/hsvimage/hsvcolor"
)
```

Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott+hsv@pakin.org*
