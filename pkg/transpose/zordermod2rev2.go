package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("zorder-mod2-rev2", ZOrderMod2Rev2)
}

func ZOrderMod2Rev2(n, m int, a, b *matrix.Matrix) {
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
		for rr := v; rr < v+4; rr++ {
			r0 = a.Read(rr, v)
			r1 = a.Read(rr, v+1)
			r2 = a.Read(rr, v+2)
			r3 = a.Read(rr, v+3)
			r4 = a.Read(rr, v+4)
			r5 = a.Read(rr, v+5)
			r6 = a.Read(rr, v+6)
			r7 = a.Read(rr, v+7)

			b.Write(v, rr, r0)
			b.Write(v+1, rr, r1)
			b.Write(v+2, rr, r2)
			b.Write(v+3, rr, r3)

			b.Write(v, rr+4, r4)
			b.Write(v+1, rr+4, r5)
			b.Write(v+2, rr+4, r6)
			b.Write(v+3, rr+4, r7)
		}

		for cc := v; cc < v+4; cc++ {
			r0 = b.Read(cc, v+4)
			r1 = b.Read(cc, v+5)
			r2 = b.Read(cc, v+6)
			r3 = b.Read(cc, v+7)

			r4 = a.Read(v+4, cc)
			r5 = a.Read(v+5, cc)
			r6 = a.Read(v+6, cc)
			r7 = a.Read(v+7, cc)

			b.Write(cc, v+4, r4)
			b.Write(cc, v+5, r5)
			b.Write(cc, v+6, r6)
			b.Write(cc, v+7, r7)

			b.Write(cc+4, v, r0)
			b.Write(cc+4, v+1, r1)
			b.Write(cc+4, v+2, r2)
			b.Write(cc+4, v+3, r3)
		}

		for rr := v + 4; rr < v+8; rr++ {
			r0 = a.Read(rr, v+4)
			r1 = a.Read(rr, v+5)
			r2 = a.Read(rr, v+6)
			r3 = a.Read(rr, v+7)

			b.Write(v+4, rr, r0)
			b.Write(v+5, rr, r1)
			b.Write(v+6, rr, r2)
			b.Write(v+7, rr, r3)
		}
	}
}
