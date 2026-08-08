package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("tiled8", Tiled8)
}

func Tiled8(n, m int, a, b *matrix.Matrix) {
	const tile = 8
	for i0 := 0; i0 < n; i0 += tile {
		for j0 := 0; j0 < m; j0 += tile {
			iEnd := min(i0+tile, n)
			jEnd := min(j0+tile, m)
			for i := i0; i < iEnd; i++ {
				for j := j0; j < jEnd; j++ {
					tmp := a.Read(i, j)
					b.Write(j, i, tmp)
				}
			}
		}
	}
}
