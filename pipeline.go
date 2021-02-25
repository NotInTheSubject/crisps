package crisps

import (
	"context"
	"fmt"

	"github.com/NotInTheSubject/crisps/actor"
	"github.com/NotInTheSubject/crisps/condition"
)

type connectedActor struct {
	name  string
	actor actor.Actor
}

// Builder builds a pipeline
type Builder struct {
	connectedActors []connectedActor
	// actorMap map[string]connectedActor
}

// ConditionFunc represent some condition
type ConditionFunc condition.Func

type caseState struct {
	condition ConditionFunc
	actor     actor.Actor
}

func Case(name string, cond ConditionFunc, actor actor.Actor) caseState {
	return caseState{
		condition: cond,
		actor:     actor,
	}
}

func (b *Builder) Append(name string, a actor.Actor) *Builder {
	b.connectedActors = append(b.connectedActors, connectedActor{name: name, actor: a})
	return b
}

type DumpFunc func(ctx context.Context, nameTrace []string)

func (b *Builder) Dump(f DumpFunc) *Builder {
	trace := []string{}
	for _, a := range b.connectedActors{
		trace = append(trace, a.name)
	}
	return b.Append("dump", actor.ActorFunc(func(ctx context.Context) context.Context {
		f(ctx, trace)
		return ctx
	}))
}

func Cycle(condition ConditionFunc, a actor.Actor) actor.Actor {
	return actor.ActorFunc(func(ctx context.Context) context.Context {
		for condition(ctx) {
			ctx = a.Act(ctx)
		}
		return ctx
	})
}

func Switch(cases ...caseState) actor.Actor {
	a := actor.ActorFunc(func(ctx context.Context) context.Context {
		for _, actCase := range cases {
			if actCase.condition(ctx) {
				return actCase.actor.Act(ctx)
			}
		}
		panic(fmt.Sprintf("cannot found true condition for ctx \"%+v\"", ctx))
	})
	return a
}

func (b *Builder) Build() actor.Actor {
	return actor.ActorFunc(func(ctx context.Context) context.Context {
		for _, item := range b.connectedActors {
			ctx = item.actor.Act(ctx)
		}
		return ctx
	})
}
