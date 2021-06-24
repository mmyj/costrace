package costrace

import (
	"context"
	"fmt"
	"time"
)

type costRaceKey int

var key costRaceKey

type constNode struct {
	startTime time.Time
	endTime   time.Time
	title     string
	child     []*constNode
}

func (n *constNode) cost() time.Duration {
	return n.endTime.Sub(n.startTime)
}

func New(ctx context.Context, title string) context.Context {
	return context.WithValue(ctx, key, &constNode{title: title, startTime: time.Now()})
}

func Done(ctx context.Context) {
	father, ok := ctx.Value(key).(*constNode)
	if !ok {
		return
	}
	father.endTime = time.Now()
}

// Trace 计算单个函数耗时
func Trace(ctx context.Context, title string, fn func()) {
	father, ok := ctx.Value(key).(*constNode)
	if !ok {
		fn()
		return
	}

	child := &constNode{title: title, startTime: time.Now()}
	fn()
	child.endTime = time.Now()
	father.child = append(father.child, child)
}

// TraceInto 需要计算函数内部耗时
func TraceInto(ctx context.Context, title string, fn func(ctx context.Context)) {
	father, ok := ctx.Value(key).(*constNode)
	if !ok {
		fn(ctx)
		return
	}

	this := &constNode{title: title, startTime: time.Now()}
	fn(context.WithValue(ctx, key, this))
	this.endTime = time.Now()
	father.child = append(father.child, this)
}

func ToString(ctx context.Context) (ret string) {
	father, ok := ctx.Value(key).(*constNode)
	if !ok {
		return ""
	}
	const fmtStr = "%s%s (%dms)\n"
	var levelPrint func(level int, node *constNode, prefix string)
	levelPrint = func(level int, node *constNode, prefix string) {
		var (
			lastTabs   string
			noLastTabs string
		)
		noLastTabs = prefix + "├"
		lastTabs = prefix + "└"
		for i, child := range node.child {
			if i == len(node.child)-1 {
				ret += fmt.Sprintf(fmtStr, lastTabs, child.title, child.cost().Milliseconds())
			} else {
				ret += fmt.Sprintf(fmtStr, noLastTabs, child.title, child.cost().Milliseconds())
			}
			if len(child.child) > 0 {
				if i == len(node.child)-1 {
					levelPrint(level+1, child, prefix+"\t")
				} else {
					levelPrint(level+1, child, prefix+"│\t")
				}
			}
		}
	}
	ret += fmt.Sprintf(fmtStr, "", father.title, father.cost().Milliseconds())
	levelPrint(0, father, "")
	return
}
