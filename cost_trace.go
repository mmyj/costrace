package costrace

import (
	"context"
	"time"
)

type (
	costRaceKey int
)

var (
	key costRaceKey
)

type constNode struct {
	startTime time.Time
	endTime   time.Time
	isDone    bool
	title     string
	child     []*constNode

	parallel bool
	childCh  chan *constNode
}

func (n *constNode) cost() time.Duration {
	return n.endTime.Sub(n.startTime)
}

func New(ctx context.Context, title string) context.Context {
	return context.WithValue(ctx, key, &constNode{title: title, startTime: time.Now()})
}

func Done(ctx context.Context) {
	this, ok := ctx.Value(key).(*constNode)
	if !ok || this.isDone {
		return
	}
	this.isDone = true
	this.endTime = time.Now()
}

func Trace(ctx context.Context, title string, fn func(ctx context.Context)) {
	father, ok := ctx.Value(key).(*constNode)
	if !ok {
		fn(ctx)
		return
	}
	this := &constNode{title: title, startTime: time.Now()}
	fn(context.WithValue(ctx, key, this))
	this.endTime = time.Now()
	if father.parallel {
		father.childCh <- this
		return
	}
	father.child = append(father.child, this)
}

func SegmentTrace(ctx context.Context, title string) context.Context {
	father, ok := ctx.Value(key).(*constNode)
	if !ok {
		return ctx
	}

	this := &constNode{title: title, startTime: time.Now()}
	if father.parallel {
		father.childCh <- this
	} else {
		father.child = append(father.child, this)
	}
	return context.WithValue(ctx, key, this)
}

func SegmentDone(ctx context.Context) {
	this, ok := ctx.Value(key).(*constNode)
	if !ok || this.isDone {
		return
	}
	this.isDone = true
	this.endTime = time.Now()
}

func ParallelTrace(ctx context.Context, parallel int) context.Context {
	father, ok := ctx.Value(key).(*constNode)
	if !ok {
		return ctx
	}

	this := &constNode{title: "[parallel]", startTime: time.Now(), parallel: true}
	this.childCh = make(chan *constNode, parallel)
	father.child = append(father.child, this)
	return context.WithValue(ctx, key, this)
}

func ParallelDone(ctx context.Context) {
	this, ok := ctx.Value(key).(*constNode)
	if !ok || this.isDone {
		return
	}
	this.isDone = true
	this.endTime = time.Now()
	close(this.childCh)
	for child := range this.childCh {
		this.child = append(this.child, child)
	}
}
