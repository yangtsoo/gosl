// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package plt

import (
	"math"

	"github.com/cpmech/gosl/io"
)

// AutoScale rescales plot area
func AutoScale(P [][]float64) {
	if len(P) < 1 {
		return
	}
	xmin, ymin := P[0][0], P[0][1]
	xmax, ymax := xmin, ymin
	for _, p := range P {
		if p[0] < xmin {
			xmin = p[0]
		}
		if p[1] < ymin {
			ymin = p[1]
		}
		if p[0] > xmax {
			xmax = p[0]
		}
		if p[1] > ymax {
			ymax = p[1]
		}
	}
	io.Ff(&bufferPy, "plt.axis([%g, %g, %g, %g])\n", xmin, xmax, ymin, ymax)
}

// Arrow adds arrow to plot
//   styles:
//     Curve           -        None
//     CurveB          ->       head_length=0.4,head_width=0.2
//     BracketB        -[       widthB=1.0,lengthB=0.2,angleB=None
//     CurveFilledB    -|>      head_length=0.4,head_width=0.2
//     CurveA          <-       head_length=0.4,head_width=0.2
//     CurveAB         <->      head_length=0.4,head_width=0.2
//     CurveFilledA    <|-      head_length=0.4,head_width=0.2
//     CurveFilledAB   <|-|>    head_length=0.4,head_width=0.2
//     BracketA        ]-       widthA=1.0,lengthA=0.2,angleA=None
//     BracketAB       ]-[      widthA=1.0,lengthA=0.2,angleA=None,widthB=1.0,lengthB=0.2,angleB=None
//     Fancy           fancy    head_length=0.4,head_width=0.4,tail_width=0.4
//     Simple          simple   head_length=0.5,head_width=0.5,tail_width=0.2
//     Wedge           wedge    tail_width=0.3,shrink_factor=0.5
//     BarAB           |-|      widthA=1.0,angleA=None,widthB=1.0,angleB=None
func Arrow(xi, yi, xf, yf float64, args *A) {
	style := "simple"
	scale := 20.0
	if args.Style != "" {
		style = args.Style
	}
	if args.Scale > 0 {
		scale = args.Scale
	}
	uid := genUid()
	io.Ff(&bufferPy, "pc%d = pat.FancyArrowPatch((%g,%g),(%g,%g),shrinkA=0,shrinkB=0,path_effects=[pff.Stroke(joinstyle='miter')],arrowstyle='%s',mutation_scale=%g", uid, xi, yi, xf, yf, style, scale)
	updateBufferAndClose(&bufferPy, args, false, false)
	io.Ff(&bufferPy, "plt.gca().add_patch(pc%d)\n", uid)
}

// Circle adds circle to plot
func Circle(xc, yc, r float64, args *A) {
	uid := genUid()
	io.Ff(&bufferPy, "pc%d = pat.Circle((%g,%g), %g", uid, xc, yc, r)
	updateBufferAndClose(&bufferPy, args, false, false)
	io.Ff(&bufferPy, "plt.gca().add_patch(pc%d)\n", uid)
}

// Arc adds arc to plot
//  minAlpha and maxAlpha are in degrees
func Arc(xc, yc, r, minAlpha, maxAlpha float64, args *A) {
	uid := genUid()
	r2 := 2.0 * r
	θ1 := minAlpha * 180.0 / math.Pi
	θ2 := maxAlpha * 180.0 / math.Pi
	io.Ff(&bufferPy, "pc%d = pat.Arc((%g,%g),%g,%g,angle=0,theta1=%g,theta2=%g", uid, xc, yc, r2, r2, θ1, θ2)
	updateBufferAndClose(&bufferPy, args, false, false)
	io.Ff(&bufferPy, "plt.gca().add_patch(pc%d)\n", uid)
}

// Polyline draws a polyline. P[npts][2]
func Polyline(P [][]float64, args *A) {
	if len(P) < 1 {
		return
	}
	uid := genUid()
	io.Ff(&bufferPy, "dat%d = [[pth.Path.MOVETO, [%g, %g]]", uid, P[0][0], P[0][1])
	for _, p := range P {
		io.Ff(&bufferPy, ", [pth.Path.LINETO, [%g, %g]]", p[0], p[1])
	}
	closed := true
	if args != nil {
		closed = args.Closed
	}
	if closed {
		io.Ff(&bufferPy, ", [pth.Path.CLOSEPOLY, [0, 0]]")
	}
	io.Ff(&bufferPy, "]\n")
	io.Ff(&bufferPy, "commands%d, vertices%d = zip(*dat%d)\n", uid, uid, uid)
	io.Ff(&bufferPy, "ph%d = pth.Path(vertices%d, commands%d)\n", uid, uid, uid)
	io.Ff(&bufferPy, "pc%d = pat.PathPatch(ph%d", uid, uid)
	updateBufferAndClose(&bufferPy, args, false, false)
	io.Ff(&bufferPy, "plt.gca().add_patch(pc%d)\n", uid)
}
