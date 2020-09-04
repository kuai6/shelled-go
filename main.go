package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jmoiron/sqlx"

	"shelled-backend/shelled"
	"shelled-backend/storage/sqlite3"
	"shelled-backend/udp"
)

func main() {
	c, err := LoadConfig("")
	if err != nil {
		log.Fatal(err)
	}

	db, err := sqlx.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	if err := sqlite3.InitStorage(db); err != nil {
		log.Fatalf("failed to init db: %s", err)
	}

	ds := sqlite3.NewDeviceStorage(db)
	bs := sqlite3.NewBankStorage(db)

	app := new(shelled.Shelled)
	if err := app.Init(ds, bs); err != nil {
		log.Fatal(err)
	}

	addr, err := net.ResolveUDPAddr("udp4", fmt.Sprintf("%s:%d", c.UDP.Host, c.UDP.Port))
	if err != nil {
		log.Fatalf("could not create udp server: resolve ip error: %s", err)
	}

	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatalf("udp server: listen error: %s", err)
	}
	defer func() {
		if conn != nil {
			_ = conn.Close()
		}
	}()

	client := udp.NewClient(conn)

	handlers := map[udp.CMD]udp.Handler{
		udp.CMD_HELLO: udp.NewHelloHandler(app, client),
		udp.CMD_PING:  udp.NewPingHandler(app, client),
	}

	server := udp.NewServer(conn, handlers)

	ctx := context.Background()

	stop := make(chan bool)

	go server.Run(ctx, stop)

	for {
		time.Sleep(1 * time.Second)
	}

	if err := app.Destroy(); err != nil {
		log.Fatal(err)
	}
}
