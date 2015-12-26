package main

import (
	"net/http"
	"flag"
	"github.com/golang/glog"
	"os"
	"mime/multipart"
	"io"
	)

func main() {
    flag.Parse()

    /* TODO add parametrization */
    glog.Info("Starting File Server on port 6543")
    http.Handle("/cdn/", http.StripPrefix("/cdn", http.FileServer(http.Dir("."))))
    glog.Info("Starting API Server on port 6543 path")
    http.HandleFunc("/api/upload", func(res http.ResponseWriter, req *http.Request) {
				var err error
				defer func() {
				    if nil != err {
					//http.Error(res, err.Error(), status)
					glog.Error("Error string: errStr", err)
					res.Write([]byte(`{"status":"error", "message":"` + err.Error()  + `"}`))
				    }
				}()
				// max bytes in mem at a time  
				const _24K = (1 << 20) * 24
				err = req.ParseMultipartForm(_24K) 
				if nil != err {
				    return
				}
				for _, fheaders := range req.MultipartForm.File {
				    for _, hdr := range fheaders {
					// open uploaded  
					var infile multipart.File
					infile, err = hdr.Open()
					if nil != err {
					     return
					}
					// open destination  
					var outfile *os.File
					outfile, err = os.Create(hdr.Filename)
					if nil != err {
					     return
					}
					// 32K buffer copy  
					var written int64
					written, err = io.Copy(outfile, infile)
					if nil != err || 0 == written {
					     return
					}
					//res.Write([]byte("uploaded file:" + hdr.Filename + ";length:" + strconv.Itoa(int(written))))
					res.Write([]byte(`{"status":"ok", "resource_url":"` + `path_to_storage` + hdr.Filename + `"}`))
				    }
				}
    })
    http.ListenAndServe(":6543", nil)
}

