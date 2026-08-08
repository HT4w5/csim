package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("zorder", ZOrder)
}

func ZOrder(n, m int, a, b *matrix.Matrix) {
	bits := 0
	for max(n, m)>>bits > 0 {
		bits++
	}
	for z := 0; z < 1<<(2*bits); z++ {
		var i, j int
		for k := 0; k < bits; k++ {
			i |= ((z >> (2 * k)) & 1) << k
			j |= ((z >> (2*k + 1)) & 1) << k
		}
		if i < n && j < m {
			tmp := a.Read(i, j)
			b.Write(j, i, tmp)
		}
	}
}
