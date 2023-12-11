package main

import (
	"crypto/md5"
	"fmt"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
)

func main() {
	workdir, err := os.Getwd()
	fmt.Printf("workdir :%s\n", workdir)
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%s %s | Range = %s", r.Method, r.URL, r.Header.Get("Range"))

			next.ServeHTTP(w, r)
		})
	})

	r.PathPrefix("/fs/").Handler(http.StripPrefix("/fs", http.FileServer(http.Dir(workdir))))
	r.PathPrefix("/md5/").Handler(http.StripPrefix("/md5", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		offset, err := strconv.ParseInt(r.Header.Get("X-Checksum-Offset"), 10, 64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		length, err := strconv.ParseInt(r.Header.Get("X-Checksum-Length"), 10, 64)
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		f, err := os.Open(path.Clean(r.URL.Path))
		if err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() { _ = f.Close() }()

		if _, err = f.Seek(offset, io.SeekStart); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		h := md5.New()
		rr := io.LimitReader(f, length)
		if _, err = io.Copy(h, rr); err != nil {
			fmt.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_, _ = w.Write(h.Sum(nil))
	})))

	log.Printf("Server working in %q and listening on :8080 ...", workdir)
	if err = http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}
}
