package es

import (
	_ "embed"
	"golang.org/x/sync/errgroup"
	"time"

	"context"
	"github.com/olivere/elastic/v7"
)

var (
	//go:embed user_index.json
	userIndex string
	//go:embed post_index.json
	postIndex string
)

func tryCreateIndex(ctx context.Context, client *elastic.Client, indexName, indexConfig string) error {
	ok, err := client.IndexExists(indexName).Do(ctx)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	_, err = client.CreateIndex(indexName).Body(indexConfig).Do(ctx)
	return err
}

func initES(client *elastic.Client) error {
	const timeout = time.Second * 10
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	var eg errgroup.Group
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, userIndex, userIndex)
	})
	eg.Go(func() error {
		return tryCreateIndex(ctx, client, userIndex, userIndex)
	})
	return eg.Wait()
}
