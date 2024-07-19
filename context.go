package grpt

import ctx "context"

type availableSpaceKey struct{}

func SetAvailableSpace(context ctx.Context, size Size) ctx.Context {
	return ctx.WithValue(context, availableSpaceKey{}, size)
}

func GetAvailableSpace(context ctx.Context) Size {
	size, _ := context.Value(availableSpaceKey{}).(Size)
	return size
}
