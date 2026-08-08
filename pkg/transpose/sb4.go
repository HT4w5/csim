package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func transposeBlock8(n, m int, a, b *matrix.Matrix, r, c int) {
	rEnd := min(r+8, n)
	cEnd := min(c+8, m)

	if rEnd-r < 8 || cEnd-c < 8 {
		for i := r; i < rEnd; i++ {
			for j := c; j < cEnd; j++ {
				b.Write(j, i, a.Read(i, j))
			}
		}
		return
	}

	var r0, r1, r2, r3, r4, r5, r6, r7 int32
	for rr := r; rr < r+4; rr++ {
		r0 = a.Read(rr, c)
		r1 = a.Read(rr, c+1)
		r2 = a.Read(rr, c+2)
		r3 = a.Read(rr, c+3)
		r4 = a.Read(rr, c+4)
		r5 = a.Read(rr, c+5)
		r6 = a.Read(rr, c+6)
		r7 = a.Read(rr, c+7)

		b.Write(c, rr, r0)
		b.Write(c+1, rr, r1)
		b.Write(c+2, rr, r2)
		b.Write(c+3, rr, r3)

		b.Write(c, rr+4, r4)
		b.Write(c+1, rr+4, r5)
		b.Write(c+2, rr+4, r6)
		b.Write(c+3, rr+4, r7)
	}

	for cc := c; cc < c+4; cc++ {
		r0 = b.Read(cc, r+4)
		r1 = b.Read(cc, r+5)
		r2 = b.Read(cc, r+6)
		r3 = b.Read(cc, r+7)

		r4 = a.Read(r+4, cc)
		r5 = a.Read(r+5, cc)
		r6 = a.Read(r+6, cc)
		r7 = a.Read(r+7, cc)

		// Write bottom-right first.
		b.Write(cc, r+4, r4)
		b.Write(cc, r+5, r5)
		b.Write(cc, r+6, r6)
		b.Write(cc, r+7, r7)

		b.Write(cc+4, r, r0)
		b.Write(cc+4, r+1, r1)
		b.Write(cc+4, r+2, r2)
		b.Write(cc+4, r+3, r3)
	}

	for rr := r + 4; rr < r+8; rr++ {
		r0 = a.Read(rr, c+4)
		r1 = a.Read(rr, c+5)
		r2 = a.Read(rr, c+6)
		r3 = a.Read(rr, c+7)

		b.Write(c+4, rr, r0)
		b.Write(c+5, rr, r1)
		b.Write(c+6, rr, r2)
		b.Write(c+7, rr, r3)
	}
}
