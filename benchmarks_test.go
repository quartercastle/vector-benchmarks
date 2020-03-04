package test

import (
	"fmt"
	"math"
	"math/rand"
	"testing"

	"github.com/kvartborg/vector"
	"gonum.org/v1/gonum/mat"
)

type vec = vector.Vector

const runs = 12

func BenchmarkAddition(b *testing.B) {
	randomListOfFloat64 := func(size int) []float64 {
		l := make([]float64, size)

		for i := range l {
			l[i] = rand.Float64()
		}

		return l
	}

	additions := []struct {
		name string
		exec func(size int, a, b []float64)
	}{
		{"gonum", func(size int, a, b []float64) {
			result := mat.NewVecDense(size, nil)
			v1 := mat.NewVecDense(size, a)
			v2 := mat.NewVecDense(size, b)
			result.AddVec(v1, v2)
		}},
		{"vector (gonum style)", func(size int, a, b []float64) {
			result := make(vec, size)
			result.Add(a, b)
		}},
		{"vector (immutable)", func(size int, a, b []float64) {
			vector.Add(a, b)
		}},
		{"vector (inline)", func(size int, a, b []float64) {
			vec(a).Add(b)
		}},
	}

	n := int(math.Pow(2, runs))
	v1, v2 := randomListOfFloat64(n), randomListOfFloat64(n)

	for _, addition := range additions {
		for k := 0.; k <= runs; k++ {
			n := int(math.Pow(2, k))
			b.Run(fmt.Sprintf("%s/%d", addition.name, n), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					addition.exec(n, v1[:n], v2[:n])
				}
			})
		}
	}
}

// func BenchmarkSubtraction(b *testing.B) {
// 	randomListOfFloat64 := func(size int) []float64 {
// 		l := make([]float64, size)
//
// 		for i := range l {
// 			l[i] = rand.Float64()
// 		}
//
// 		return l
// 	}
//
// 	additions := []struct {
// 		name string
// 		exec func(size int, a, b []float64)
// 	}{
// 		{"gonum", func(size int, a, b []float64) {
// 			result := mat.NewVecDense(size, nil)
// 			v1 := mat.NewVecDense(size, a)
// 			v2 := mat.NewVecDense(size, b)
//
// 			result.SubVec(v1, v2)
// 		}},
// 		{"vector", func(size int, a, b []float64) {
// 			result := make(vec, size)
// 			v1, v2 := vec(a), vec(b)
//
// 			result.Sub(v1, v2)
// 		}},
// 	}
//
// 	for _, addition := range additions {
// 		for k := 0.; k <= 12; k++ {
// 			n := int(math.Pow(2, k))
// 			v1, v2 := randomListOfFloat64(n), randomListOfFloat64(n)
// 			b.Run(fmt.Sprintf("%s/%d", addition.name, n), func(b *testing.B) {
// 				for i := 0; i < b.N; i++ {
// 					addition.exec(n, v1, v2)
// 				}
// 			})
// 		}
// 	}
// }
//
// func BenchmarkScale(b *testing.B) {
// 	randomListOfFloat64 := func(size int) []float64 {
// 		l := make([]float64, size)
//
// 		for i := range l {
// 			l[i] = rand.Float64()
// 		}
//
// 		return l
// 	}
//
// 	additions := []struct {
// 		name string
// 		exec func(size int, a, b []float64)
// 	}{
// 		{"gonum", func(size int, a, b []float64) {
// 			result := mat.NewVecDense(size, nil)
// 			v := mat.NewVecDense(size, a)
//
// 			result.ScaleVec(2, v)
// 		}},
// 		{"vector", func(size int, a, b []float64) {
// 			result := make(vec, size)
// 			v := vec(a)
//
// 			result.Add(v).Scale(2)
// 		}},
// 	}
//
// 	for _, addition := range additions {
// 		for k := 0.; k <= 12; k++ {
// 			n := int(math.Pow(2, k))
// 			v1, v2 := randomListOfFloat64(n), randomListOfFloat64(n)
// 			b.Run(fmt.Sprintf("%s/%d", addition.name, n), func(b *testing.B) {
// 				for i := 0; i < b.N; i++ {
// 					addition.exec(n, v1, v2)
// 				}
// 			})
// 		}
// 	}
// }

// func BenchmarkMultiAddition(b *testing.B) {
// 	randomListOfFloat64 := func(size int) []float64 {
// 		l := make([]float64, size)
//
// 		for i := range l {
// 			l[i] = rand.Float64()
// 		}
//
// 		return l
// 	}
//
// 	additions := []struct {
// 		name string
// 		exec func(size int, vs [][]float64)
// 	}{
// 		{"gonum", func(size int, in [][]float64) {
// 			result := mat.NewVecDense(size, nil)
// 			vs := []*mat.VecDense{
// 				mat.NewVecDense(size, in[0]),
// 				mat.NewVecDense(size, in[1]),
// 				mat.NewVecDense(size, in[2]),
// 				mat.NewVecDense(size, in[3]),
// 				mat.NewVecDense(size, in[4]),
// 				mat.NewVecDense(size, in[5]),
// 				mat.NewVecDense(size, in[6]),
// 				mat.NewVecDense(size, in[7]),
// 			}
//
// 			for i := range vs {
// 				result.AddVec(result, vs[i])
// 			}
// 		}},
// 		{"vector", func(size int, in [][]float64) {
// 			result := make(vec, size)
// 			vs := []vec{
// 				vec(in[0]), vec(in[1]),
// 				vec(in[2]), vec(in[3]),
// 				vec(in[4]), vec(in[5]),
// 				vec(in[6]), vec(in[7]),
// 			}
//
// 			result.Add(vs...)
// 		}},
// 	}
//
// 	for _, addition := range additions {
// 		for k := 0.; k <= 16; k++ {
// 			n := int(math.Pow(2, k))
// 			vs := [][]float64{
// 				randomListOfFloat64(n),
// 				randomListOfFloat64(n),
// 				randomListOfFloat64(n),
// 				randomListOfFloat64(n),
// 				randomListOfFloat64(n),
// 				randomListOfFloat64(n),
// 				randomListOfFloat64(n),
// 				randomListOfFloat64(n),
// 			}
// 			b.Run(fmt.Sprintf("%s/%d", addition.name, n), func(b *testing.B) {
// 				for i := 0; i < b.N; i++ {
// 					addition.exec(n, vs)
// 				}
// 			})
// 			runtime.GC()
// 		}
// 	}
// }
