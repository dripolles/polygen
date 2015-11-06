package fixtures

//go:generate polygen github.com/dripolles/polygen/generator/fixtures/convertslice.tgo -t"a:float64"

func compilationCheck() []interface{} {
	xs := []float64{.1, .2, .3}

	return convertfloat64slice(xs)
}
