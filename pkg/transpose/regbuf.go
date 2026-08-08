package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("regbuf", RegBuf)
}

func RegBuf(n, m int, a, b *matrix.Matrix) {
	var r0, r1, r2, r3, r4, r5, r6, r7 int32
	size := m * n
	for i := 0; i < size; i += 8 {
		for j := range 8 {
			switch j {
			case 0:
				r0 = a.Read(i/m, i%m)
			case 1:
				r1 = a.Read((i+1)/m, (i+1)%m)
			case 2:
				r2 = a.Read((i+2)/m, (i+2)%m)
			case 3:
				r3 = a.Read((i+3)/m, (i+3)%m)
			case 4:
				r4 = a.Read((i+4)/m, (i+4)%m)
			case 5:
				r5 = a.Read((i+5)/m, (i+5)%m)
			case 6:
				r6 = a.Read((i+6)/m, (i+6)%m)
			case 7:
				r7 = a.Read((i+7)/m, (i+7)%m)
			}
		}
		for j := range 8 {
			switch j {
			case 0:
				b.Write(i%m, i/m, r0)
			case 1:
				b.Write((i+1)%m, (i+1)/m, r1)
			case 2:
				b.Write((i+2)%m, (i+2)/m, r2)
			case 3:
				b.Write((i+3)%m, (i+3)/m, r3)
			case 4:
				b.Write((i+4)%m, (i+4)/m, r4)
			case 5:
				b.Write((i+5)%m, (i+5)/m, r5)
			case 6:
				b.Write((i+6)%m, (i+6)/m, r6)
			case 7:
				b.Write((i+7)%m, (i+7)/m, r7)
			}
		}
	}
}
