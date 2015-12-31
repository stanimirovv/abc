package main

import (
	// General
	"github.com/golang/glog"
	_ "github.com/iambc/xerrors"
	"flag"

	//API
	"net/http"

	//DB
	"github.com/jmoiron/sqlx"
	_ "database/sql"
	_ "github.com/lib/pq"
)



// sample usage
func main() {
    flag.Parse()
    http.HandleFunc("/api", func(res http.ResponseWriter, req *http.Request) {
					//check if function exists, parse parameters to a map[string]string
					values := req.URL.Query()
					glog.Info("values: ", values)
					command, is_passed := values[`command`]
					if !is_passed {
					    res.Write([]byte(`Invalid params!`))
					    return
					}
					if command[0] == `get_boards` {
					    res.Write([]byte(`List of boards:`))
					}

					if command[0] == `get_threads` {
					    res.Write([]byte(`List of threads:`))
					}

					if command[0] == `get_post_for_thread` {
					    res.Write([]byte(`List of posts:`))

					}

					if command[0] == `post_to_thread` {
					    res.Write([]byte(`Posted new thread!!`))
					}
    })

    //Test db connect
    _, err := sqlx.Connect("postgres", "user=abc_api password=123 dbname=abc_dev_cluster sslmode=disable")

    http.ListenAndServe(`:8089`, nil)
}



