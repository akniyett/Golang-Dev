package main

import (
	"context"
	"6/internal/cache/redis_cache"
	"6/internal/http"
	"6/internal/store/postgres"
	"log"
)

const (
	port = ":8081"
	cacheDB = 1
	cacheExpTime = 1800
	cachePort = "localhost:6379"
)

func main() {
	urlDB := "postgres://postgres:Impossible@localhost5432/goproject"
	store := postgres.NewDB()
	if err := store.Connect(urlDB); err != nil {
		panic(err)
	}
	defer store.Close()


	cache := redis_cache.NewRedisCache(cachePort, cacheDB, cacheExpTime)
	

	srv := http.NewServer(context.Background(),
		http.WithAddress(port),
		http.WithStore(store),
		http.WithCache(cache),
	)
	if err := srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()
}