package transpose

import "github.com/HT4w5/csim/pkg/matrix"

type Func struct {
	Name string
	F    matrix.TransposeFunc
}

var Functions []Func

var Sizes = [][2]int{{32, 32}, {61, 67}, {64, 64}}

func register(name string, f matrix.TransposeFunc) {
	Functions = append(Functions, Func{Name: name, F: f})
}
