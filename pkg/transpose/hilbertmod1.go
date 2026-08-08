package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("hilbert-order-mod1", HilbertMod1)
}

func HilbertMod1(n, m int, a, b *matrix.Matrix) {
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

	var r0, r1, r2, r3, r4, r5, r6, r7 int32
	for v := 0; v < 64; v += 8 {
		for i := range 8 {
			for j := range 8 {
				switch j {
				case 0:
					r0 = a.Read(v+i, v+j)
				case 1:
					r1 = a.Read(v+i, v+j)
				case 2:
					r2 = a.Read(v+i, v+j)
				case 3:
					r3 = a.Read(v+i, v+j)
				case 4:
					r4 = a.Read(v+i, v+j)
				case 5:
					r5 = a.Read(v+i, v+j)
				case 6:
					r6 = a.Read(v+i, v+j)
				case 7:
					r7 = a.Read(v+i, v+j)
				}
			}
			for j := range 8 {
				switch j {
				case 0:
					b.Write(v+j, v+i, r0)
				case 1:
					b.Write(v+j, v+i, r1)
				case 2:
					b.Write(v+j, v+i, r2)
				case 3:
					b.Write(v+j, v+i, r3)
				case 4:
					b.Write(v+j, v+i, r4)
				case 5:
					b.Write(v+j, v+i, r5)
				case 6:
					b.Write(v+j, v+i, r6)
				case 7:
					b.Write(v+j, v+i, r7)
				}
			}
		}
	}
}
