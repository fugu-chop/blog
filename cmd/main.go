package main

import (
	"fmt"
	"net/http"
)

func main() {
	svr := http.NewServeMux()

	svr.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello there!")
	})

	fmt.Println("starting server on localhost at port 3000...")
	http.ListenAndServe(":3000", svr)
}
