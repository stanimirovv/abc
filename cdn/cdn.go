package main

import (
	"fmt"
	"net/http"
	)

func main() {
            fmt.Printf("Hello, world!\n")
	    http.Handle("/cdn/", http.StripPrefix("/cdn", http.FileServer(http.Dir("<path>"))))
	    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
				// todo add insert method
				fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	    })
	    //todo maybe add remove method? 
	    http.ListenAndServe(":6543", nil)
}

