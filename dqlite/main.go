package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/canonical/go-dqlite/app"
	"github.com/canonical/go-dqlite/client"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

func GenRandomBytes(size int) (blk []byte, err error) {
	blk = make([]byte, size)
	_, err = rand.Read(blk)
	return
}

func main() {
	var api string
	var db string
	var join *[]string
	var dir string
	var verbose bool

	// Configure CLI.
	cmd := &cobra.Command{
		Use:   "seguro-dqlite",
		Short: "dqlite benchmark for urbit integration",
		Long: `A performance benchmark for dqlite to determine its viability
		for integration with the Urbit binary as part of the Seguro project.
		https://github.com/wexpert/seguro`,
		RunE: func(cmd *cobra.Command, args []string) error {
			dir := filepath.Join(dir, db)
			if err := os.MkdirAll(dir, 0755); err != nil {
				return errors.Wrapf(err, "can't create %s", dir)
			}
			logFunc := func(l client.LogLevel, format string, a ...interface{}) {
				if !verbose {
					return
				}
				log.Printf(fmt.Sprintf("%s: %s: %s\n", api, l.String(), format), a...)
			}
			// Configure dqlite cluster.
			app, err := app.New(dir, app.WithAddress(db), app.WithCluster(*join), app.WithLogFunc(logFunc))
			if err != nil {
				return err
			}

			if err := app.Ready(context.Background()); err != nil {
				return err
			}

			// Open new dqlite database.
			db, err := app.Open(context.Background(), "seguro-dqlite")
			if err != nil {
				return err
			}

			// Create a table.
			if _, err := db.Exec(schema); err != nil {
				return err
			}

			// Sleep for a little.
			time.Sleep(60 * time.Minute)

			db.Close()

			app.Handover(context.Background())
			app.Close()

			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&db, "db", "d", "", "address used for internal database replication")
	join = flags.StringSliceP("join", "j", nil, "database addresses of existing nodes")
	flags.StringVarP(&dir, "dir", "D", "/tmp/seguro-dqlite", "data directory")
	flags.BoolVarP(&verbose, "verbose", "v", false, "verbose logging")

	cmd.MarkFlagRequired("db")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

const (
	//  SQL statements.
	schema = "CREATE TABLE IF NOT EXISTS events (key INT, value BLOB, UNIQUE(key))"
	query  = "SELECT value FROM events WHERE key = ?"
	insert = "INSERT INTO events(key, value) VALUES(?, ?)"
)
