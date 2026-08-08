package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("hilbert-order", Hilbert)
}

func Hilbert(n, m int, a, b *matrix.Matrix) {
	bits := 0
	for max(n, m)>>bits > 0 {
		bits++
	}
	size := 1 << bits
	for d := 0; d < size*size; d++ {
		i, j := d2xy(size, d)
		if i < n && j < m {
			tmp := a.Read(i, j)
			b.Write(j, i, tmp)
		}
	}
}

func d2xy(size, d int) (x, y int) {
	for s := 1; s < size; s *= 2 {
		rx := 1 & (d / 2)
		ry := 1 & (d ^ rx)
		x, y = rot(s, x, y, rx, ry)
		x += s * rx
		y += s * ry
		d /= 4
	}
	return x, y
}

func rot(n, x, y, rx, ry int) (int, int) {
	if ry == 0 {
		if rx == 1 {
			x = n - 1 - x
			y = n - 1 - y
		}
		x, y = y, x
	}
	return x, y
}
