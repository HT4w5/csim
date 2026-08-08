package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("tiled8regbuf", Tiled8RegBuf)
}

func Tiled8RegBuf(n, m int, a, b *matrix.Matrix) {
	var r0, r1, r2, r3, r4, r5, r6, r7 int32
	const tile = 8
	for i0 := 0; i0 < n; i0 += tile {
		for j0 := 0; j0 < m; j0 += tile {
			iEnd := min(i0+tile, n)
			jEnd := min(j0+tile, m)
			for i := i0; i < iEnd; i++ {
				for j := 0; j+j0 < jEnd; j++ {
					switch j {
					case 0:
						r0 = a.Read(i, j+j0)
					case 1:
						r1 = a.Read(i, j+j0)
					case 2:
						r2 = a.Read(i, j+j0)
					case 3:
						r3 = a.Read(i, j+j0)
					case 4:
						r4 = a.Read(i, j+j0)
					case 5:
						r5 = a.Read(i, j+j0)
					case 6:
						r6 = a.Read(i, j+j0)
					case 7:
						r7 = a.Read(i, j+j0)
					}
				}
				for j := 0; j+j0 < jEnd; j++ {
					switch j {
					case 0:
						b.Write(j+j0, i, r0)
					case 1:
						b.Write(j+j0, i, r1)
					case 2:
						b.Write(j+j0, i, r2)
					case 3:
						b.Write(j+j0, i, r3)
					case 4:
						b.Write(j+j0, i, r4)
					case 5:
						b.Write(j+j0, i, r5)
					case 6:
						b.Write(j+j0, i, r6)
					case 7:
						b.Write(j+j0, i, r7)
					}
				}
			}
		}
	}
}
