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

    /* TODO add parametrization 
	- server path to cdn TODO
	- server upload url TODO
    */

    /*
    * issues: files with same name can be overridden
    * there is literally no security
    * They will be solved on the second iteration of the cdn
    */

    portNumber, err := strconv.Atoi(os.Getenv("ABC_CDN_PORTNUM"))
    if nil != err {
	glog.Fatal("Error: ", err, portNumber)
    }

    if os.Getenv("ABC_CDN_DIR") == "" {
	 glog.Fatal("Undefinfed var ABC_CDN_DIR!")
    }

    if os.Getenv("ABC_CDN_ENDPOINT_URL") == "" {
	 glog.Fatal("Undefinfed var ABC_CDN_ENDPOINT_URL!")
    }
    ABC_CDN_ENDPOINT_URL := os.Getenv("ABC_CDN_ENDPOINT_URL")

    glog.Info(`Starting File Server on port:`, os.Getenv("ABC_CDN_PORTNUM"), `path: `, os.Getenv("ABC_CDN_DIR") )
    http.Handle("/cdn/", http.StripPrefix("/cdn", http.FileServer(http.Dir(os.Getenv("ABC_CDN_DIR")))))
    glog.Info("Starting API Server...")
    http.HandleFunc("/api/upload", func(res http.ResponseWriter, req *http.Request) {
				var err error
				defer func() {
				    if nil != err {
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
					outfile, err = os.Create(os.Getenv("ABC_CDN_DIR") + hdr.Filename)
					if nil != err {
					     return
					}
					var written int64
					written, err = io.Copy(outfile, infile)
					if nil != err || 0 == written {
					     return
					}
					glog.Info("Succesfully uploaded file: ", ABC_CDN_ENDPOINT_URL + hdr.Filename)
					res.Write([]byte(`{"status":"ok", "resource_url":"` + ABC_CDN_ENDPOINT_URL + hdr.Filename + `"}`))
				    }
				}
    })
    http.ListenAndServe(`:` + os.Getenv("ABC_CDN_PORTNUM"), nil)
}

