package  netcontext

import "context"

var(
	RootCTX context.Context
)

func Init()  {
	RootCTX = context.Background()
}