package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"sync"

	"github.com/redis/go-redis/v9"
)

type WarmUpClient struct {
	client *redis.Client
}

type Records struct {
	Id   string
	Slot string
}

func newWarmUpClient(client *redis.Client) *WarmUpClient {
	return &WarmUpClient{client: client}
}

// Pipeline per worker approach
func (c *WarmUpClient) warmup(ctx context.Context, file_name string, workers int) {
	file, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	records, err := csv.NewReader(file).ReadAll()
	if err != nil {
		panic(err)
	}

	ch := make(chan Records, workers)
	var wg sync.WaitGroup

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for r := range ch {
				if err := c.client.Set(ctx, r.Id, r.Slot, 0).Err(); err != nil {
					fmt.Printf("failed to set %s: %v\n", r.Id, err)
				}
			}
		}()
	}

	for _, row := range records[1:] {
		ch <- Records{Id: row[0], Slot: row[1]}
	}
	close(ch)
	wg.Wait()

	fmt.Println("Bulk import complete")
}

func (c *WarmUpClient) BulkWarmup(ctx context.Context, file_name string) {
	file, err := os.Open(file_name)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	pipeline := c.client.Pipeline()
	for _, row := range records[1:] {
		id := row[0]
		slot := row[1]
		pipeline.Set(ctx, id, slot, 0)

	}
	_, err = pipeline.Exec(ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Bulk import complete")
}

func (c *WarmUpClient) unload(ctx context.Context) {
	err := c.client.FlushDB(ctx).Err()
	if err != nil {
		panic(err)
	}
	fmt.Println("Data flushed")
}
