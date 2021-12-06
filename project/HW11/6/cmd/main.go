package main

import (
	"context"
	"6/internal/cache/redis_cache"
	"6/internal/http"
	"6/internal/store/postgres"
	"6/internal/http"
	"6/internal/message_broker/kafka"
	"6/internal/store/postgres"
	lru "github.com/hashicorp/golang-lru"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	port = ":8081"
	cacheDB = 1
	cacheExpTime = 1800
	cachePort = "localhost:6379"
)

func main() {




	ctx, cancel := context.WithCancel(context.Background())
	go CatchTermination(cancel)

	dbURL := "postgres://postgres:Impossible@localhost5432/goproject"
	store := postgres.NewDB()
	if err := store.Connect(dbURL); err != nil {
		panic(err)
	}
	defer store.Close()

	cache, err := lru.New2Q(10)
	if err != nil {
		panic(err)
	}

	brokers := []string{"localhost:29092"}
	broker := kafka.NewBroker(brokers, cache, "peer3")
	if err = broker.Connect(ctx); err != nil {
		panic(err)
	}
	defer broker.Close()

	srv := http.NewServer(context.Background(),
		http.WithAddress(port),
		http.WithStore(store),
		http.WithCache(cache),
		http.WithBroker(broker),
	)
	if err = srv.Run(); err != nil {
		log.Println(err)
	}

	srv.WaitForGracefulTermination()

	// urlDB := "postgres://postgres:Impossible@localhost5432/goproject"
	// store := postgres.NewDB()
	// if err := store.Connect(urlDB); err != nil {
	// 	panic(err)
	// }
	// defer store.Close()


	// cache := redis_cache.NewRedisCache(cachePort, cacheDB, cacheExpTime)
	

	// srv := http.NewServer(context.Background(),
	// 	http.WithAddress(port),
	// 	http.WithStore(store),
	// 	http.WithCache(cache),
	// )
	// if err := srv.Run(); err != nil {
	// 	log.Println(err)
	// }

	// srv.WaitForGracefulTermination()
}

func CatchTermination(cancel context.CancelFunc) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	log.Print("[WARN] caught termination signal")
	cancel()
}




