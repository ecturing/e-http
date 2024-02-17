package netcontext

import "context"

func worker(ctx context.Context,f func()) {
	select{
	case <-ctx.Done():
		return
	default:f()
	}
}