package main

import (
	"context"
	"crypto/rand"
	"fmt"
	mr "math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/canonical/go-dqlite/app"
	"github.com/pkg/errors"
)

const (
	ip     = "127.0.0.1"
	dir    = "/tmp/seguro-dqlite"
	events = 10
)

func genEvents(length int) (events [][]byte) {
	fmt.Printf("generating events... ")
	events = make([][]byte, length)
	min := 1000
	max := 32000
	eventLength := mr.Intn(max-min) + min // [1,32] KB
	for i := range events {
		events[i] = make([]byte, eventLength)
		rand.Read(events[i])
	}
	fmt.Printf("done\n")
	return
}

// Starts a dqlite master node.
func master(addr string, mrdy chan bool, ardy chan bool, done chan bool) (node *app.App, err error) {
	// Log.
	fmt.Printf("starting master... ")

	// Create data directory.
	dir := filepath.Join(dir, addr)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, errors.Wrapf(err, "can't create %s", dir)
	}

	// Start the master node.
	node, err = app.New(dir, app.WithAddress(addr))
	if err != nil {
		return nil, err
	}

	// Wait for readiness.
	if err := node.Ready(context.Background()); err != nil {
		return nil, err
	}

	// Open a new dqlite database.
	db, err := node.Open(context.Background(), "seguro-dqlite")
	if err != nil {
		return nil, err
	}

	// Create the events table.
	schema := `
	  CREATE TABLE IF NOT EXISTS events 
		(
			key INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, 
			value BLOB, UNIQUE(key)
		)
	`
	if _, err := db.Exec(schema); err != nil {
		return nil, err
	}

	// Signal readiness to replica node.
	mrdy <- true
	fmt.Printf("done\n")

	// Wait for replica node to be ready, then insert test data.
	<-ardy
	fakeEvents := genEvents(events)
	fmt.Printf("writing events to dqlite... ")
	start := time.Now()
	insert := "INSERT INTO events(value) VALUES(?)"
	for i := range fakeEvents {
		if _, err := db.Exec(insert, fakeEvents[i]); err != nil {
			fmt.Println(err.Error())
		}
	}
	elapsed := time.Since(start)
	fmt.Printf("done in %s\n", elapsed)

	// Signal completion to the replica node.
	done <- true

	// Gracefully shutdown the node.
	db.Close()
	node.Handover(context.Background())
	node.Close()

	return node, nil
}

// Starts a dqlite replica node.
func replica(addr string, join string, mrdy chan bool, ardy chan bool, done chan bool) (node *app.App, err error) {
	// Log.
	fmt.Println("starting replica... ")

	// Create data directory.
	dir := filepath.Join(dir, addr)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, errors.Wrapf(err, "can't create %s", dir)
	}

	// Wait for the master to be ready, then start the replica node.
	<-mrdy
	node, err = app.New(dir, app.WithAddress(addr), app.WithCluster([]string{join}))
	if err != nil {
		return nil, err
	}

	// Wait for readiness.
	if err := node.Ready(context.Background()); err != nil {
		return nil, err
	}

	// Signal readiness to the master.
	ardy <- true
	fmt.Println("done!")

	// Wait for the master to finish inserting data.
	<-done

	// Gracefully shutdown the node.
	node.Handover(context.Background())
	node.Close()

	return node, nil
}

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(2)

	mrdy := make(chan bool, 1)
	ardy := make(chan bool, 1)
	done := make(chan bool, 1)

	go func() {
		defer wg.Done()
		master(ip+":9000", mrdy, ardy, done)
	}()

	go func() {
		defer wg.Done()
		replica(ip+":9001", ip+":9000", mrdy, ardy, done)
	}()

	wg.Wait()
}
