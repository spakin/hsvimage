hsvimage and hsvcolor
=====================

[![Go Report Card](https://goreportcard.com/badge/github.com/spakin/hsvimage)](https://goreportcard.com/report/github.com/spakin/hsvimage)
[![GoDoc](https://godoc.org/github.com/spakin/hsvimage?status.svg)](https://godoc.org/github.com/spakin/hsvimage)

The [Go programming language](https://golang.org/)'s standard library provides support for manipulating graphical images represented using [RGB](https://en.wikipedia.org/wiki/RGB_color_model), [CMYK](https://en.wikipedia.org/wiki/CMYK_color_model), [YCbCr](https://en.wikipedia.org/wiki/YCbCr), and [grayscale](https://en.wikipedia.org/wiki/Grayscale) color models and variations such as 8- vs. 16-bit channels, premultiplied alpha vs. non-premultiplied alpha vs. no alpha, and [paletted color](https://en.wikipedia.org/wiki/Indexed_color).  `hsvimage` augments the Go standard library with support for the [HSV](https://en.wikipedia.org/wiki/HSL_and_HSV) color model.

Usage
-----

`hsvimage` provides an [`hsvimage.NHSVA`](https://godoc.org/github.com/spakin/hsvimage#NHSVA) image type (**N**on-alpha-premultiplied **H**ue, **S**aturation, and **V**alue with **A**lpha) that implements the [`image.Image`](https://golang.org/pkg/image/#Image) interface.  The underlying [NHSVA color model](https://godoc.org/github.com/spakin/hsvimage/hsvcolor#pkg-variables) and [NHSVA data type](https://godoc.org/github.com/spakin/hsvimage/hsvcolor#NHSVA) are provided by `hsvimage/hsvcolor`.  `hsvimage` and `hsvimage/hsvcolor` can be imported in the usual manner:


```Go
import (
	"github.com/spakin/hsvimage"
	"github.com/spakin/hsvimage/hsvcolor"
)
```

Author
------

[Scott Pakin](http://www.pakin.org/~scott/), *scott+hsv@pakin.org*
