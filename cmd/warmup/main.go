package main

import (
	"context"
	"flag"
	"os"
	"path/filepath"

	bootstrap "parallel/init"
	"parallel/internal/service"
)

func main() {
	csv := flag.String("csv", "class.csv", "path to CSV file")
	unload := flag.Bool("unload", false, "flush Redis DB")
	flag.Parse()

	exe, _ := os.Executable()
	root := filepath.Dir(exe)

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
