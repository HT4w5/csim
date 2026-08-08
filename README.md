# csim
Cache simulator and matrix transpose cache usage visualizer. Written while doing the famous cache lab.

## Usage
#### Writing solutions
Write your own transpose solutions in `./pkg/transpose/<solution>.go`. E.g.
```go
package transpose

import "github.com/HT4w5/csim/pkg/matrix"

func init() {
	register("row-wise", RowWise)
}

func RowWise(n, m int, a, b *matrix.Matrix) {
	for i := range n {
		for j := range m {
			tmp := a.Read(i, j)
			b.Write(j, i, tmp)
		}
	}
}
```
##### Provided matrices
- `a[n][m]`: source int32 matrix. Has `n` rows and `m` columns.
- `b[m][n]`: target int32 matrix. Has `m` rows and `n` columns.
- You're supposed to transpose a into b.

##### Matrix interface
```go
func (mx *Matrix) Read(x, y int) int32
func (mx *Matrix) Write(x, y int, v int32)
```
- `x`: matrix row index.
- `y`: matrix column index.

#### Evaluating solutions
```plaintext
Usage of csim:
  -E int
        lines per set (default 1)
  -b int
        block offset bits (default 5)
  -build string
        build report pages and exit if set
  -listen string
        listen address (default "127.0.0.1:8081")
  -s int
        set index bits (default 5)
```

Preview reports in browser:
```shell
go run cmd/csim/main.go
```

Build reports:
```shell
go run cmd/csim/main.go -build ./out
```

## Reports
Generated report pages are hosted on [pages](https://ht4w5.github.io/csim).

## Credits
- [go-echarts](https://github.com/go-echarts/go-echarts): MIT.

