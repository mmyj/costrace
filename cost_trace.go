package costrace

import (
	"context"
	"fmt"
	"time"
)

const key = "cost_race_key"

type constNode struct {
	startTime time.Time
	endTime   time.Time
	title     string
	child     []*constNode
}

func (n *constNode) cost() int64 {
	return n.endTime.Sub(n.startTime).Milliseconds()
}

func Init(ctx context.Context, title string) context.Context {
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
	const fmtStr = "%s%s:%dms\n"
	var levelPrint func(level int, node *constNode)
	levelPrint = func(level int, node *constNode) {
		var (
			tabs       string
			lastTabs   string
			noLastTabs string
		)
		if level > 0 {
			for i := 0; i < level; i++ {
				tabs += "│\t"
			}
		}
		noLastTabs = tabs + "├"
		lastTabs = tabs + "└"
		for i, child := range node.child {
			if i == len(node.child)-1 {
				ret += fmt.Sprintf(fmtStr, lastTabs, child.title, child.cost())
			} else {
				ret += fmt.Sprintf(fmtStr, noLastTabs, child.title, child.cost())
			}
			if len(child.child) > 0 {
				levelPrint(level+1, child)
			}
		}
	}
	ret += fmt.Sprintf(fmtStr, "", father.title, father.cost())
	levelPrint(0, father)
	return
}
