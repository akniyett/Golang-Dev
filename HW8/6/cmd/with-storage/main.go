package main

import (
	"context"
	"6/internal/http"
	"6/internal/store/postgres"
)

func main() {
	urlExample := "postgres://localhost:5432/clothings"
	store := postgres.NewDB()

	if err := store.Connect(urlExample); err != nil {
		panic(err)
	}
	defer store.Close()

	srv := http.NewServer(context.Background(), ":8080", store)
	if err := srv.Run(); err != nil {
		panic(err)
	}

	srv.WaitForGracefulTermination()
}
