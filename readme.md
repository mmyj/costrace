# constrace

## trace the function
```
    ctx := costrace.New(context.Background(), "trace the function")
    defer func() {
        costrace.Done(ctx)
        fmt.Print(costrace.ToString(ctx))
    }()
    costrace.Trace(ctx, "func1", func(context.Context) {
        time.Sleep(time.Second)
    })
    sctx := costrace.SegmentTrace(ctx, "func2")
    time.Sleep(time.Second)
    costrace.SegmentDone(sctx)
    
    // nest
    sctx = costrace.SegmentTrace(ctx, "nest-func1")
    time.Sleep(time.Second)
    sctx2 := costrace.SegmentTrace(sctx, "nest-func2")
    time.Sleep(time.Second)
    costrace.SegmentDone(sctx2)
    costrace.SegmentDone(sctx)
```
```
trace the function (4013ms 100%)
├─func1 (1004ms 25%)
├─func2 (1002ms 24%)
└─nest-func1 (2005ms 49%)
  └─nest-func2 (1003ms 50%)
```

## trace the parallel function
```
	ctx := costrace.New(context.Background(), "trace the parallel function")
	defer func() {
		costrace.Done(ctx)
		fmt.Print(costrace.ToString(ctx))
	}()

	const parallel = 3
	var wg sync.WaitGroup
	wg.Add(3)

	pctx := costrace.ParallelTrace(ctx, parallel)
	go func() {
		defer wg.Done()
		costrace.Trace(pctx, "func1", func(context.Context) {
			time.Sleep(time.Second)
		})
	}()

	go func() {
		defer wg.Done()
		sctx := costrace.SegmentTrace(pctx, "func2")
		time.Sleep(time.Second)
		costrace.SegmentDone(sctx)
	}()

	go func() {
		defer wg.Done()
		// nest
		sctx := costrace.SegmentTrace(pctx, "nest-func1")
		time.Sleep(time.Second)
		sctx2 := costrace.SegmentTrace(sctx, "nest-func2")
		time.Sleep(time.Second)
		costrace.SegmentDone(sctx2)
		costrace.SegmentDone(sctx)
	}()
	wg.Wait()
	costrace.ParallelDone(pctx)
```

```
trace the parallel function (2005ms 100%)
└─[parallel] (2005ms 100%)
  ├─nest-func1 (2005ms 100%)
  │  └─nest-func2 (1005ms 50%)
  ├─func2 (1000ms 49%)
  └─func1 (1000ms 49%)
```