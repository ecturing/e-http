package main

import (
	"ews/ehttp"
)

func main() {
	r := ehttp.NewRouter()
	ehttp.Server(r, "/server", server, ehttp.POST)
	ehttp.Confirm(":8080")
}

func server(rq *ehttp.E_Request, rp *ehttp.E_Response) {
	rp.Headers["Content-Type"] = "text/plain"
	rp.DataFrom = rq.Proto
}
