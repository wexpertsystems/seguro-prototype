package benchmark

import (
	"context"
	"database/sql"
	"math/rand"
	"time"
)

type work int
type workerType int

func (w work) String() string {
	switch w {
	case exec:
		return "exec"
	case query:
		return "query"
	case none:
		return "none"
	default:
		return "unknown"
	}
}

const (
	// The type of query to perform
	none  work = iota
	exec  work = iota // a `write`
	query work = iota // a `read`

	kvWriter       workerType = iota
	kvReader       workerType = iota
	kvReaderWriter workerType = iota

	kvReadSql  = "SELECT value FROM model WHERE key = ?"
	kvWriteSql = "INSERT INTO model(value) VALUES(?)"
)

// A worker performs the queries to the database and keeps around some state
// in order to do that. `lastWork` and `lastArgs` refer to the previously
// executed operation and can be used to determine the next work the worker
// should perform. `kvKeys` tells the worker the highest key integer it has inserted in the
// database (keys are auto-incrementing integers).
type worker struct {
	workerType   workerType
	lastWork     work
	lastArgs     []interface{}
	tracker      *tracker
	kvKeySizeB   int
	kvValueSizeB int
	kvLastKey    int
}

func (w *worker) randExistingKey() (int, error) {
	return rand.Intn(w.kvLastKey-1) + 1, nil
}

// A random byte slice (mock event).
func (w *worker) randValue() []byte {
	v := make([]byte, w.kvValueSizeB)
	rand.Read(v)
	return v
}

// Returns the type of work to execute and a sql statement with arguments
func (w *worker) getWork() (work, string, []interface{}) {
	switch w.workerType {
	case kvWriter:
		v := w.randValue()
		return exec, kvWriteSql, []interface{}{v}
	case kvReaderWriter:
		read := rand.Intn(2) == 0
		if read && w.kvLastKey != 0 {
			k, _ := w.randExistingKey()
			return query, kvReadSql, []interface{}{k}
		}
		v := w.randValue()
		return exec, kvWriteSql, []interface{}{v}
	default:
		return none, "", []interface{}{}
	}
}

// Retrieve a query and execute it against the database
func (w *worker) doWork(ctx context.Context, db *sql.DB) {
	var err error
	var str string

	work, q, args := w.getWork()
	w.lastWork = work
	w.lastArgs = args

	switch work {
	case exec:
		w.kvLastKey = w.kvLastKey + 1
		defer w.tracker.measure(time.Now(), work, &err)
		_, err = db.ExecContext(ctx, q, args...)
		if err != nil {
			w.kvLastKey = w.kvLastKey - 1
		}
	case query:
		defer w.tracker.measure(time.Now(), work, &err)
		err = db.QueryRowContext(ctx, q, args...).Scan(&str)
	default:
		return
	}
}

func (w *worker) run(ctx context.Context, db *sql.DB) {
	for {
		if ctx.Err() != nil {
			return
		}

		w.doWork(ctx, db)
	}
}

func (w *worker) report() map[work]report {
	return w.tracker.report()
}

func newWorker(workerType workerType, o *options) *worker {
	return &worker{
		workerType:   workerType,
		kvKeySizeB:   o.kvKeySizeB,
		kvValueSizeB: o.kvValueSizeB,
		tracker:      newTracker(),
	}
}
