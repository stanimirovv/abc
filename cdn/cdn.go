package main

import (
	//"net/http"
	"flag"
	"github.com/golang/glog"	
	)

func main() {
    flag.Parse()
    
    glog.Info("Something something INfo")
    glog.Error("Something something Error")
    glog.Warning("Something something Warning")
/*
	    printHello();
	    http.Handle("/cdn/", http.StripPrefix("/cdn", http.FileServer(http.Dir("<path>"))))
	    http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
				// todo add insert method
				fmt.Fprintf(w, "Hello, %q", r.URL.Path)
	    })
	    //todo maybe add remove method? 
	    http.ListenAndServe(":6543", nil)
*/
}

