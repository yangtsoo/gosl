// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package fun

import (
	"math"

	"github.com/cpmech/gosl/chk"
	"github.com/cpmech/gosl/io"
)

// enums
var (
	OP_JACOBI     = io.NewEnum("Jacobi", "num.op", "J", "Jacobi OrthoPolynomial")     // Jacobi OrthoPolynomial
	OP_LEGENDRE   = io.NewEnum("Legendre", "num.op", "L", "Legendre OrthoPolynomial") // Legendre OrthoPolynomial
	OP_HERMITE    = io.NewEnum("Hermite", "num.op", "H", "Hermite OrthoPolynomial")   // Hermite OrthoPolynomial
	OP_CHEBYSHEV1 = io.NewEnum("Chebyshev1", "num.op", "T", "Chebyshev First Kind")   // Chebyshev1 OrthoPolynomial
	OP_CHEBYSHEV2 = io.NewEnum("Chebyshev2", "num.op", "U", "Chebyshev Second Kind")  // Chebyshev2 OrthoPolynomial
)

// OrthoPolynomial (main) structure ////////////////////////////////////////////////////////////////

// OrthoPolynomial implements orthogonal polynomials
type OrthoPolynomial struct {

	// input
	Type io.Enum // type of orthogonal polynomial
	N    int     // (max) degree of polynomial. Lower order can be quickly obtained after this polynomial with max(N) is generated

	// computed
	c [][]float64 // all c coefficients [N+1][M+1]

	// internal
	poly oPoly // implementation
}

// NewOrthoPolynomial creates a new orthogonal polynomial
//   Type -- is the type: e.g. OP_JACOBI, OP_LEGENDRE, OP_HERMITE
//   N    -- is the (max) degree of the polynomial.
//           Lower order can later be quickly obtained after this
//           polynomial with max(N) is created
//   Prms -- parameters to each particular polynomial
func NewOrthoPolynomial(Type io.Enum, N int, prms Prms) (o *OrthoPolynomial) {
	o = new(OrthoPolynomial)
	o.Type = Type
	o.N = N
	o.poly = newopoly(Type, prms)
	o.c = make([][]float64, o.N+1)
	for n := 1; n <= o.N; n++ {
		o.c[n] = make([]float64, o.poly.M(o.N)+1)
		M := o.poly.M(n)
		for m := 0; m <= M; m++ {
			o.c[n][m] = o.poly.c(n, m)
		}
	}
	return
}

// F computes P(n,x) with n=N (max)
func (o *OrthoPolynomial) F(x float64) (res float64) {
	return o.P(o.N, x)
}

// P computes P(n,x) where n must be ≤ N
func (o *OrthoPolynomial) P(n int, x float64) (res float64) {
	if n > o.N {
		chk.Panic("the degree n must not be greater than max N. %d > %d", n, o.N)
	}
	if n == 0 {
		return 1
	}
	for m := 0; m <= o.poly.M(n); m++ {
		res += o.c[n][m] * o.poly.g(n, m, x)
	}
	res *= o.poly.d(n)
	return
}

// oPoly database //////////////////////////////////////////////////////////////////////////////////

// oPoly defines the functions that OrthoPolynomial must have
type oPoly interface {
	M(n int) int
	d(n int) float64
	c(n, m int) float64
	g(n, m int, x float64) float64
}

// oPolyMaker defines a function that makes new oPolys
type oPolyMaker func(prms Prms) oPoly

// oPolyDB implements a database of oPoly makers
var oPolyDB map[io.Enum]oPolyMaker = make(map[io.Enum]oPolyMaker)

// newopoly finds oPoly or panic
func newopoly(code io.Enum, prms Prms) oPoly {
	if maker, ok := oPolyDB[code]; ok {
		return maker(prms)
	}
	chk.Panic("cannot find OrthoPolynomial named %q in database", code)
	return nil
}

// Jacobi //////////////////////////////////////////////////////////////////////////////////////////

type opJacobi struct {
	alpha float64
	beta  float64
}

func (o *opJacobi) M(n int) int {
	return n
	//return int(math.Floor(float64(n / 2.0))) // L
}

func (o *opJacobi) d(n int) float64 {
	var twon uint64 = 1 << uint64(n) // 1<<n = 2ⁿ
	return 1.0 / float64(twon)
}

func (o *opJacobi) c(n, m int) float64 {
	r := Rbinomial(float64(n)+o.alpha, float64(m))
	s := Rbinomial(float64(n)+o.beta, float64(n-m))
	return r * s
}

func (o *opJacobi) g(n, m int, x float64) float64 {
	return math.Pow(x-1, float64(n-m)) * math.Pow(x+1, float64(m))
}

func newJacobi(prms Prms) oPoly {
	o := new(opJacobi)
	prms.CheckLimits()
	prms.ConnectSet([]*float64{&o.alpha, &o.beta}, []string{"alpha", "beta"}, "newJacobi")
	return o
}

// Legendre //////////////////////////////////////////////////////////////////////////////////////////

type opLegendre struct{}

func (o *opLegendre) M(n int) int {
	return int(math.Floor(float64(n) / 2.0))
}

func (o *opLegendre) d(n int) float64 {
	var twon uint64 = 1 << uint64(n) // 1<<n = 2ⁿ
	return 1.0 / float64(twon)
}

func (o *opLegendre) c(n, m int) float64 {
	r := Rbinomial(float64(n), float64(m))
	s := Rbinomial(float64(2*n-2*m), float64(n))
	return math.Pow(-1, float64(m)) * r * s
}

func (o *opLegendre) g(n, m int, x float64) float64 {
	return math.Pow(x, float64(n-2*m))
}

func newLegendre(prms Prms) oPoly {
	return new(opLegendre)
}

// Hermite //////////////////////////////////////////////////////////////////////////////////////////

type opHermite struct{}

func (o *opHermite) M(n int) int {
	return int(math.Floor(float64(n) / 2.0))
}

func (o *opHermite) d(n int) float64 {
	return Factorial22(n)
}

func (o *opHermite) c(n, m int) float64 {
	r := Factorial22(m)
	s := Factorial22(n - 2*m)
	return math.Pow(-1, float64(m)) / (r * s)
}

func (o *opHermite) g(n, m int, x float64) float64 {
	return math.Pow(2*x, float64(n-2*m))
}

func newHermite(prms Prms) oPoly {
	return new(opHermite)
}

// Chebyshev1 //////////////////////////////////////////////////////////////////////////////////////////

type opChebyshev1 struct{}

func (o *opChebyshev1) M(n int) int {
	return int(math.Floor(float64(n) / 2.0))
}

func (o *opChebyshev1) d(n int) float64 {
	return float64(n) / 2.0
}

func (o *opChebyshev1) c(n, m int) float64 {
	r := Factorial22(n - m - 1)
	s := Factorial22(m)
	t := Factorial22(n - 2*m)
	return math.Pow(-1, float64(m)) * r / (s * t)
}

func (o *opChebyshev1) g(n, m int, x float64) float64 {
	return math.Pow(2*x, float64(n-2*m))
}

func newChebyshev1(prms Prms) oPoly {
	return new(opChebyshev1)
}

// Chebyshev2 //////////////////////////////////////////////////////////////////////////////////////////

type opChebyshev2 struct{}

func (o *opChebyshev2) M(n int) int {
	return int(math.Floor(float64(n) / 2.0))
}

func (o *opChebyshev2) d(n int) float64 {
	return 1.0
}

func (o *opChebyshev2) c(n, m int) float64 {
	r := Factorial22(n - m)
	s := Factorial22(m)
	t := Factorial22(n - 2*m)
	return math.Pow(-1, float64(m)) * r / (s * t)
}

func (o *opChebyshev2) g(n, m int, x float64) float64 {
	return math.Pow(2*x, float64(n-2*m))
}

func newChebyshev2(prms Prms) oPoly {
	return new(opChebyshev2)
}

// add polynomials to database /////////////////////////////////////////////////////////////////////

func init() {
	oPolyDB[OP_JACOBI] = newJacobi
	oPolyDB[OP_LEGENDRE] = newLegendre
	oPolyDB[OP_HERMITE] = newHermite
	oPolyDB[OP_CHEBYSHEV1] = newChebyshev1
	oPolyDB[OP_CHEBYSHEV2] = newChebyshev2
}