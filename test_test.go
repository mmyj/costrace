package costrace_test

import (
	"context"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/mmyj/costrace"
)

func a(ctx context.Context) {
	time.Sleep(time.Millisecond * 10)
	costrace.Trace(ctx, "a2", func(context.Context) {
		a2()
	})
	costrace.Trace(ctx, "a1", func(ctx context.Context) {
		a1(ctx)
	})
}
func a1(ctx context.Context) {
	time.Sleep(time.Millisecond * 60)
	costrace.Trace(ctx, "a3", func(context.Context) {
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
	costrace.Trace(ctx, "b1", func(context.Context) {
		b1()
	})
	costrace.Trace(ctx, "b2", func(context.Context) {
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

func TestTrace(t *testing.T) {
	ctx := costrace.New(context.Background(), "Test1")
	defer func() {
		costrace.Done(ctx)
		fmt.Print(costrace.ToString(ctx))
	}()
	costrace.Trace(ctx, "c", func(context.Context) {
		c()
	})
	costrace.Trace(ctx, "a", func(ctx context.Context) {
		a(ctx)
	})
	costrace.Trace(ctx, "b", func(ctx context.Context) {
		b(ctx)
	})
}

func TestSegmentTrace(t *testing.T) {
	ctx := costrace.New(context.Background(), "Test3")
	defer func() {
		costrace.Done(ctx)
		fmt.Print(costrace.ToString(ctx))
	}()
	ctxSeg := costrace.SegmentTrace(ctx, "cost of segment")
	time.Sleep(time.Millisecond * 10)
	costrace.SegmentDone(ctxSeg)
	costrace.Trace(ctx, "a", func(ctx context.Context) {
		a(ctx)
	})
}

func TestGoroutineTrace(t *testing.T) {
	ctx := costrace.New(context.Background(), "TestGoroutineTrace")
	defer func() {
		costrace.Done(ctx)
		fmt.Print(costrace.ToString(ctx))
	}()
	var wg sync.WaitGroup
	wg.Add(2)

	parallelCtx := costrace.ParallelTrace(ctx, 2)
	go func(ctx context.Context) {
		costrace.Trace(ctx, "a", func(ctx context.Context) {
			a(ctx)
		})
		wg.Done()
	}(parallelCtx)
	go func(ctx context.Context) {
		ctxSeg := costrace.SegmentTrace(ctx, "cost of segment")
		time.Sleep(time.Millisecond * 10)
		costrace.SegmentDone(ctxSeg)
		wg.Done()
	}(parallelCtx)
	wg.Wait()
	costrace.ParallelDone(parallelCtx)
}
