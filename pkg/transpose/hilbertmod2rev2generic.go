package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("hilbert-order-mod2-rev2-generic", HilbertMod2Rev2Generic)
}

func HilbertMod2Rev2Generic(n, m int, a, b *matrix.Matrix) {
	bits := 0
	for max(n, m)>>bits > 0 {
		bits++
	}
	size := 1 << bits
	for d := 0; d < size*size; d++ {
		i, j := d2xy(size, d)
		if i/8 == j/8 {
			continue
		}
		if i < n && j < m {
			tmp := a.Read(i, j)
			b.Write(j, i, tmp)
		}
	}

	for v := 0; v < min(n, m); v += 8 {
		transposeBlock8(n, m, a, b, v, v)
	}
}
