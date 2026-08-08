package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("tiled16x4", Tiled16x4)
}

func Tiled16x4(n, m int, a, b *matrix.Matrix) {
	const tileX = 16
	const tileY = 4
	for i0 := 0; i0 < n; i0 += tileX {
		for j0 := 0; j0 < m; j0 += tileY {
			iEnd := min(i0+tileX, n)
			jEnd := min(j0+tileY, m)
			for i := i0; i < iEnd; i++ {
				for j := j0; j < jEnd; j++ {
					tmp := a.Read(i, j)
					b.Write(j, i, tmp)
				}
			}
		}
	}
}
