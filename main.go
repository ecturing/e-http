package main

import (
	"ews/ehttp"
)

func main() {
	r := ehttp.NewRouter()
	ehttp.Server(r,"/server", server,ehttp.GET)
	ehttp.Confirm(":8080")
}

func server(rq *ehttp.Request, rp *ehttp.Response) {
	rp.DataFrom = "reponse"
}
