package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("tiled8sb4rev2-generic", Tiled8Subblock4Rev2Generic)
}

func Tiled8Subblock4Rev2Generic(n, m int, a, b *matrix.Matrix) {
	for r := 0; r < n; r += 8 {
		for c := 0; c < m; c += 8 {
			transposeBlock8(n, m, a, b, r, c)
		}
	}
}
