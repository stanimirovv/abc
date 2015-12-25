package main

import (
	"net/http"
	"flag"
	"github.com/golang/glog"
	"os"
	"mime/multipart"
	"io"
	"strconv"
	)

func main() {
    flag.Parse()

    glog.Info("Starting File Server on port 6543")
    http.Handle("/cdn/", http.StripPrefix("/cdn", http.FileServer(http.Dir("."))))
    glog.Info("Starting API Server on port 6543 path")
    /* TODO add error handling! */
    http.HandleFunc("/api/upload", func(res http.ResponseWriter, req *http.Request) {
				var (
				    status int
				    err  error
				)
				defer func() {
				    if nil != err {
					http.Error(res, err.Error(), status)
				    }
				}()
				// max bytes in mem at a time  
				const _24K = (1 << 20) * 24
				if err = req.ParseMultipartForm(_24K); nil != err {
				    status = http.StatusInternalServerError
				    return
				}
				for _, fheaders := range req.MultipartForm.File {
				    for _, hdr := range fheaders {
					// open uploaded  
					var infile multipart.File
					infile, err = hdr.Open()
					if nil != err {
					     status = http.StatusInternalServerError
					     return
					}
					// open destination  
					var outfile *os.File
					outfile, err = os.Create(hdr.Filename)
					if nil != err {
					     status = http.StatusInternalServerError
					     return
					}
					// 32K buffer copy  
					var written int64
					written, err = io.Copy(outfile, infile)
					if nil != err {
					     status = http.StatusInternalServerError
					     return
					}
					res.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
				    }
				}
    })
    http.ListenAndServe(":6543", nil)
}

