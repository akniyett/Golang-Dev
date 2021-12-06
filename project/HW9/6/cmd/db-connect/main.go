// package main

// import (
// 	"context"
// 	"github.com/jackc/pgx/v4"
// 	"log"
// )

// func main() {
// 	ctx := context.Background()
// 	urlExample := "postgres://postgres:Impossible@localhost5432/goproject"
// 	conn, err := pgx.Connect(ctx, urlExample)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer conn.Close(context.Background())

// 	if err := conn.Ping(ctx); err != nil {
// 		panic(err)
// 	}
// 	log.Println("Pinged DB")
// }

package main

import (
	"encoding/json"
	"fmt"
	"6/internal/models"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	//ctx := context.Background()
	urlExample := "postgres://postgres:Impossible@localhost5432/goproject"
	conn, err := sqlx.Connect("pgx", urlExample)
	if err != nil {
		panic(err)
	}

	if err = conn.Ping(); err != nil {
		panic(err)
	}

	accessories := make([]*models.Accessory, 0)


	getAccessoryQuery := `SELECT * FROM accessories`
	if err := conn.Select(&acessories, getAccessoryQuery); err != nil {
		panic(err)
	}


	res, err := json.Marshal(acessories)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))

	accessories := make([]*models.Accessory, 0)


	getClothingQuery := `SELECT * FROM clothings`
	if err := conn.Select(&clothings, getClothingQuery); err != nil {
		panic(err)
	}


	res, err := json.Marshal(clothings)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(res))



}