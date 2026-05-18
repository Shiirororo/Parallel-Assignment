package main

import (
	"context"
	"flag"
	bootstrap "parallel/init"
	"parallel/internal/service"
)

func main() {
	csv := flag.String("csv", "class.csv", "path to CSV file")
	unload := flag.Bool("unload", false, "flush Redis DB")
	flag.Parse()

	rdb, _ := bootstrap.Run()
	client := service.NewWarmUpClient(rdb)
	ctx := context.Background()

	if *unload {
		client.Unload(ctx)
	} else {
		client.BulkWarmup(ctx, *csv)
	}
}
