package examples

import "fmt"

type Foo struct {
}

//go:generate polygen github.com/dripolles/polygen/examples/fooprocess.tgo fooint -t"a:int"
//go:generate polygen github.com/dripolles/polygen/examples/fooprocess.tgo foofloat -t"a:float"

func (f *Foo) BarInt(xs []int) int {
	xs := []int{1, 2, 3, 4}
	res, err := f.processint(xs, 0, 100)
	fmt.Println(res, err)

	return 0
}
