package costrace_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mmyj/costrace"
)

func a(ctx context.Context) {
	time.Sleep(time.Millisecond * 10)
	costrace.Trace(ctx, "a2", func() {
		a2()
	})
	costrace.TraceInto(ctx, "a1", func(ctx2 context.Context) {
		a1(ctx2)
	})
}
func a1(ctx context.Context) {
	time.Sleep(time.Millisecond * 60)
	costrace.Trace(ctx, "a3", func() {
		a3()
	})
}
func a2() {
	time.Sleep(time.Millisecond * 70)
}
func a3() {
	time.Sleep(time.Millisecond * 70)
}
func b(ctx context.Context) {
	time.Sleep(time.Millisecond * 20)
	costrace.Trace(ctx, "b1", func() {
		b1()
	})
	costrace.Trace(ctx, "b2", func() {
		b2()
	})
}
func b1() {
	time.Sleep(time.Millisecond * 60)
}
func b2() {
	time.Sleep(time.Millisecond * 70)
}
func c() {
	time.Sleep(time.Millisecond * 30)
}

func Test1(t *testing.T) {
	ctx := costrace.New(context.Background(), "Test1")
	defer func() {
		costrace.Done(ctx)
		fmt.Print(costrace.ToString(ctx))
	}()
	costrace.Trace(ctx, "c", func() {
		c()
	})
	costrace.TraceInto(ctx, "a", func(ctx2 context.Context) {
		a(ctx2)
	})
	costrace.TraceInto(ctx, "b", func(ctx2 context.Context) {
		b(ctx2)
	})
}

func Test2(t *testing.T) {
	ctx := costrace.New(context.Background(), "Test2")
	defer func() {
		costrace.Done(ctx)
		fmt.Print(costrace.ToString(ctx))
	}()
	time.Sleep(time.Millisecond * 100)
}
