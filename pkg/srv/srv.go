package srv

import (
	"bytes"
	"io"
	"net/http"
)

func ListenAndServe(addr string, f func(w io.Writer) error) error {
	var buf bytes.Buffer

	if err := f(&buf); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(buf.Bytes())
	})

	return Serve(addr, mux)
}

func Serve(addr string, mux *http.ServeMux) error {
	srv := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return srv.ListenAndServe()
}
