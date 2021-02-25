package actor

import (
	"context"
)

// Actor is represents an object whose can somehow Act
type Actor interface {
	Act(context.Context) context.Context
}

type ActorFunc func(context.Context) context.Context

func (a ActorFunc) Act(ctx context.Context) context.Context {
	return a(ctx)
}

func Concat(actrors ...Actor) Actor {
	return ActorFunc(func(ctx context.Context) context.Context {
		for _, actor := range actrors {
			ctx = actor.Act(ctx)
		}
		return ctx
	})
}
