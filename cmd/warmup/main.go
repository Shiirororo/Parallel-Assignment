package main

import (
	"context"
	"flag"
	"path/filepath"
	"runtime"

	bootstrap "parallel/init"
	"parallel/internal/service"
)

func main() {
	csv := flag.String("csv", "class.csv", "path to CSV file")
	unload := flag.Bool("unload", false, "flush Redis DB")
	flag.Parse()

	_, file, _, _ := runtime.Caller(0)
	root := filepath.Join(filepath.Dir(file), "../..")

	config := bootstrap.LoadConfigFrom(root)
	rdb := bootstrap.InitRedis(config)

	client := service.NewWarmUpClient(rdb)
	ctx := context.Background()

	csvPath := filepath.Join(root, *csv)
	if *unload {
		client.Unload(ctx)
	} else {
		client.BulkWarmup(ctx, csvPath)
	}
}
