package costrace

import (
	"context"
	"fmt"
)

// ToString format a tracer to string
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
