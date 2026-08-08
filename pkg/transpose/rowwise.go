package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("row-wise", RowWise)
}

func RowWise(n, m int, a, b *matrix.Matrix) {
	for i := range n {
		for j := range m {
			tmp := a.Read(i, j)
			b.Write(j, i, tmp)
		}
	}
}
