package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
)

func main() {
	ctx := context.Background()
	urlExample := "postgres://postgres:Impossible@localhost5432/goproject"
	conn, err := pgx.Connect(ctx, urlExample)
	if err != nil {
		panic(err)
	}
	defer conn.Close(context.Background())

	if err := conn.Ping(ctx); err != nil {
		panic(err)
	}
	log.Println("Pinged DB")
}