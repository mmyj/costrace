package costrace

import (
	"context"
	"fmt"
	"time"
)

type (
	costRaceKey        int
	costRaceSegmentKey int
)

var (
	key        costRaceKey
	segmentKey costRaceSegmentKey
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
	return context.WithValue(ctx, segmentKey, this)
}

func SegmentDone(ctx context.Context) {
	this, ok := ctx.Value(segmentKey).(*constNode)
	if !ok || this.isDone {
		return
	}
	this.isDone = true
	this.endTime = time.Now()
}

func ToString(ctx context.Context) (ret string) {
	father, ok := ctx.Value(key).(*constNode)
	if !ok {
		return ""
	}
	Done(ctx)
	const fmtStr = "%s%s (%dms %d%%)\n"
	var levelPrint func(level int, node *constNode, prefix string)
	levelPrint = func(level int, node *constNode, prefix string) {
		var (
			lastTabs   string
			noLastTabs string
		)
		noLastTabs = prefix + "├─"
		lastTabs = prefix + "└─"
		for i, child := range node.child {
			tabs := noLastTabs
			if i == len(node.child)-1 {
				tabs = lastTabs
			}
			childCostMs := child.cost().Milliseconds()
			fatherCostMs := node.cost().Milliseconds()
			radio := int64(0)
			if fatherCostMs > 0 {
				radio = childCostMs * 100 / fatherCostMs
			}
			ret += fmt.Sprintf(fmtStr, tabs, child.title, childCostMs, radio)
			if len(child.child) > 0 {
				if i == len(node.child)-1 {
					levelPrint(level+1, child, prefix+"  ")
				} else {
					levelPrint(level+1, child, prefix+"│  ")
				}
			}
		}
	}
	ret += fmt.Sprintf(fmtStr, "", father.title, father.cost().Milliseconds(), 100)
	levelPrint(0, father, "")
	return
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
