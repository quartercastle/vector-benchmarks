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

const runs = 16

func randomListOfFloat64(size int) []float64 {
	l := make([]float64, size)

	for i := range l {
		l[i] = rand.Float64()
	}

	return l
}

func BenchmarkAddition(b *testing.B) {
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
	sample1, sample2 := randomListOfFloat64(n), randomListOfFloat64(n)

	for _, addition := range additions {
		for k := 0.; k <= runs; k++ {
			n := int(math.Pow(2, k))
			b.Run(fmt.Sprintf("%s/%d", addition.name, n), func(b *testing.B) {
				b.ReportAllocs()
				v1, v2 := vec(sample1[:n]).Clone(), vec(sample2[:n]).Clone()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					addition.exec(n, v1, v2)
				}
			})
		}
	}
}

func TestGonumVSVector(t *testing.T) {
	a, b := randomListOfFloat64(1000), randomListOfFloat64(1000)

	r1 := make(vec, len(a))
	v1, v2 := vec(a), vec(b)

	r1.Add(v1, v2)

	r2 := mat.NewVecDense(len(a), nil)
	gv1 := mat.NewVecDense(len(a), a)
	gv2 := mat.NewVecDense(len(b), b)

	r2.AddVec(gv1, gv2)

	for i := range a {
		if math.Abs(r1[i]-r2.AtVec(i)) > 1e-8 {
			t.Errorf(
				"value not equal at index %d; vector = %.8f; gonum = %.8f",
				i, r1[i], r2.AtVec(i),
			)
		}
	}
}

func AvgRawList(winds [][2]float64) []float64 {
	var speed float64
	r := make([]float64, 2)
	for _, wind := range winds {
		r[0], r[1] = r[0]+wind[0], r[1]+wind[1]
		speed += math.Sqrt(wind[0]*wind[0] + wind[1]*wind[1])
	}
	direction := math.Atan2(r[1], r[0])
	speed /= float64(len(winds))

	cos, sin := math.Cos(direction), math.Sin(direction)

	return []float64{
		cos*speed - sin*0,
		sin*speed + cos*0,
	}
}

type vec2d struct {
	x, y float64
}

func AvgRawStruct(winds []vec2d) []float64 {
	var speed float64
	r := vec2d{}
	for _, wind := range winds {
		r.x, r.y = r.x+wind.x, r.y+wind.y
		speed += math.Sqrt(wind.x*wind.x + wind.y*wind.y)
	}
	direction := math.Atan2(r.y, r.x)
	speed /= float64(len(winds))

	cos, sin := math.Cos(direction), math.Sin(direction)

	return []float64{
		cos*speed - sin*0,
		sin*speed + cos*0,
	}
}

func AvgVector(winds []vec) []float64 {
	r := make(vec, 2).Add(winds...)
	direction := math.Atan2(r.Y(), r.X())

	var speed float64
	for _, wind := range winds {
		speed += wind.Magnitude()
	}
	speed /= float64(len(winds))

	return vec{speed}.Rotate(direction)
}

func AvgGonum(winds []*mat.VecDense) []float64 {
	var speed float64
	r := mat.NewVecDense(2, nil)

	for _, wind := range winds {
		r.AddVec(r, wind)
		speed += mat.Norm(wind, 2)
	}

	direction := math.Atan2(r.AtVec(1), r.AtVec(0))
	speed /= float64(len(winds))

	return []float64{
		math.Cos(direction) * speed,
		math.Sin(direction) * speed,
	}
}

func generateWinds(n int) [][2]float64 {
	winds := make([][2]float64, n)

	random := func() float64 {
		return 0 + rand.Float64()*(10-0)
	}

	for i := range winds {
		winds[i][0], winds[i][1] = random(), random()
	}

	return winds
}

func BenchmarkAverageWind(b *testing.B) {
	cases := []struct {
		name string
		exec func(winds [][2]float64, b *testing.B)
	}{
		// {"struct", func(winds [][2]float64, b *testing.B) {
		// 	ws := make([]vec2d, len(winds))
		//
		// 	for i := range winds {
		// 		ws[i] = vec2d{winds[i][0], winds[i][1]}
		// 	}
		//
		// 	b.ReportAllocs()
		// 	b.ResetTimer()
		// 	for i := 0; i < b.N; i++ {
		// 		AvgRawStruct(ws)
		// 	}
		// }},
		{"list", func(winds [][2]float64, b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				AvgRawList(winds)
			}
		}},
		{"gonum", func(winds [][2]float64, b *testing.B) {
			ws := make([]*mat.VecDense, len(winds))

			for i := range winds {
				ws[i] = mat.NewVecDense(2, []float64{winds[i][0], winds[i][0]})
			}

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				AvgGonum(ws)
			}
		}},
		{"vector", func(winds [][2]float64, b *testing.B) {
			ws := make([]vec, len(winds))

			for i := range winds {
				ws[i] = vec([]float64{winds[i][0], winds[i][0]})
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				AvgVector(ws)
			}
		}},
	}

	runs := 24
	n := int(math.Pow(2, float64(runs)))
	winds := generateWinds(n)

	for _, c := range cases {
		for k := 0; k <= runs; k++ {
			n := int(math.Pow(2, float64(k)))
			b.Run(fmt.Sprintf("%s/%d", c.name, n), func(b *testing.B) {
				c.exec(winds[:n], b)
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
