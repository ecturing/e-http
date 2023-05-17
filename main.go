package main

import (
	"ews/ehttp"
)

func main() {
	ehttp.Server("/server", server)
	ehttp.Confirm(":8080")
}

func server(rq *ehttp.Request, rp *ehttp.Response) {

}
