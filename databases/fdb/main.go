package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
)

func main() {
	// Different API versions may expose different runtime behaviors.
	fdb.MustAPIVersion(710)

	// Open the default database from the system cluster.
	db := fdb.MustOpenDefault()

	// Generate some mock events.
	events := make([][]byte, 4096)
	for i := range events {
		e := make([]byte, 1024)
		rand.Read(e)
		events[i] = e
	}

	// Auto-incrementing keys are unavailable in fdb.
	max := 0.0
	min := 999.99
	sum := 0.0
	var errors = 0
	// for i, e := range events {
	// 	txStart := time.Now()
	// 	_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
	// 		tr.Set(fdb.Key(fmt.Sprint(i)), e)
	// 		return
	// 	})
	// 	if err == nil {
	// 		txElapsed := float64(time.Since(txStart).Nanoseconds()) / 1000000.0
	// 		sum = sum + txElapsed
	// 		if txElapsed > max {
	// 			max = txElapsed
	// 		}
	// 		if txElapsed < min {
	// 			min = txElapsed
	// 		}
	// 	} else {
	// 		errors++
	// 	}
	// }
	for i := 0; i < len(events)-5; i += 5 {
		txStart := time.Now()
		_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
			tr.Set(fdb.Key(fmt.Sprint(i)), events[i])
			tr.Set(fdb.Key(fmt.Sprint(i+1)), events[i+1])
			tr.Set(fdb.Key(fmt.Sprint(i+2)), events[i+2])
			tr.Set(fdb.Key(fmt.Sprint(i+3)), events[i+3])
			tr.Set(fdb.Key(fmt.Sprint(i+4)), events[i+4])
			return
		})
		if err == nil {
			txElapsed := float64(time.Since(txStart).Nanoseconds()) / 1000000.0
			sum = sum + txElapsed
			if (txElapsed / 5.0) > max {
				max = txElapsed / 5.0
			}
			if (txElapsed / 5.0) < min {
				min = txElapsed / 5.0
			}
		} else {
			errors++
		}
	}
	avg := float64(sum) / float64(len(events))
	fmt.Printf("n %d\n", len(events))
	fmt.Printf("n_err %d\n", errors)
	fmt.Printf("avg [ms] %f\n", avg)
	fmt.Printf("max [ms] %f\n", max)
	fmt.Printf("min [ms] %f\n", min)
}
