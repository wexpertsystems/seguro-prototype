package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/apple/foundationdb/bindings/go/src/fdb"
	"github.com/spf13/cobra"
)

const (
	defaultBatchSize = 5
	defaultValueSize = 1024
	defaultWorkload  = "single"
	defaultNumEvents = 4096
	docString        = "For benchmarking FoundationDB.\n\n"
	// Maximum value size supported by FoundationDB is 100,000 bytes,
	// but we'll use 10,000 as it's recommended for optimal performance.
	// See: https://apple.github.io/foundationdb/largeval.html
	maxValueSize = 10000
)

// Generates mock events without fragmentation.
func generateEvents(valueSize int, numEvents int) [][]byte {
	events := make([][]byte, numEvents)
	for i := range events {
		e := make([]byte, valueSize)
		rand.Read(e)
		events[i] = e
	}
	return events
}

/*
Accepts an event ([]byte) and fragments it into one or more values
of sizes less than or equal to the maxValueSize.

Fragment header format: <fragment index>:<total fragments>:

For example, when...

(eventSize <= (maxValueSize - 4))

The event will be written using only one k/v pair in the database
with a four-byte header encoded to "0:1:" (0x0 0x3a 0x1 0x3a).

(eventSize == 3 * (maxValueSize - 4))

The event will be written using three k/v pairs in the database
with headers encoded to "[0..2]:3:" ([0x0..0x2] 0x3a 0x3 0x3a)

In the case where an event requires more than 2^7-1 (127) database
values to be stored, the high bit of the first byte in the header will be
set, indicating that the lower 7 bits encode the byte-width of the

For example, when...

(eventSize == 256 * (maxValueSize - 4))

The event will be written using 256 k/v pairs in the database, each
with a four-byte header of "["

NOTE: This fragmentation scheme relies on ordered, monotonic integer keys, such
that any given event with a value-length of two or more will be
written in the database with a key k := k + (n-1) where k is a constant
and n specifies which number fragment it is in the whole event fragmentation
sequence (i.e., 1st, 2nd, 3rd, ... nth, and so on...).

*/
func fragmentEvent(e []byte) [][]byte {
	eventSize := len(e)
	totalNumValues, valueDataSize := getTotalNumValues(eventSize, maxValueSize)
	fragments := make([][]byte, totalNumValues)
	for i, f := range fragments {
		header := generateFragmentHeader(eventSize, i, maxValueSize)
		// Insert the header.
		for b := range header {
			f = append(f, byte(b))
		}
		// Insert the data.
		if i < totalNumValues-valueDataSize {
			firstIndex := i * valueDataSize
			lastIndex := (i+1)*valueDataSize - 1
			data := f[firstIndex:lastIndex]
			for b := range data {
				f = append(f, byte(b))
			}
		}
	}
	return fragments
}

/*
Generates a fragment header based on a given event's byte-length.

Fragment header format: <fragment index>:<total fragments>:

If the high bit of the first byte is unset, the lower 7 bits encode the length
of the value. If it is set, they encode the number of bytes that encode the
length of the value.
*/
func generateFragmentHeader(eventSize int, fragmentIndex int, maxValueSize int) []byte {
	headerSize := getFragmentHeaderSize(eventSize, maxValueSize)
	header := make([]byte, headerSize)
	if headerSize == 4 {
		header[0] = byte(fragmentIndex)
		header[1] = ':'
		header[2] = byte(eventSize)
		header[3] = ':'
	} else if headerSize == 12 {
		// TODO
		fmt.Println("IMPLEMENT ME")
	}
	return header
}

/*
Returns the byte-length of a fragment header based on an eventSize
and maxValueSize supported by the database. The byte-length returned is a sum
of the number of bytes required to store a single fragment's header in the format:

<fragment index>:<total fragments>:

For simplicity, we ensure the byte-lengths of the fragment index and total fragments
segments are equivalent. Thus, the return value of this function will be
(2 * <segment width> + 2), as each of the two ':' characters in the header require
one byte of storage.

In the future, it is possible to reduce the space required for headers by reducing
the segment width of the fragment index to be the minimum number of bytes required
to store its actual value. For example, a fragment with an index of 127 or less would
only need 1 byte of storage, even if its corresponding total fragment size is much
larger.
*/
func getFragmentHeaderSize(eventSize int, maxValueSize int) int {
	if maxValueSize < 1 {
		return 0
	}

	sizeRatio := float64(eventSize) / float64(maxValueSize)
	valueCount := math.Ceil(sizeRatio)

	var segmentWidth int
	if valueCount < 0x1<<7 {
		segmentWidth = 1
	} else {
		segmentWidth = 5
	}

	return 2*segmentWidth + 2
}

/*
Returns the total number of database values required to
store an event, based on its eventSize and the maxValueSize of the
database.
*/
func getTotalNumValues(eventSize int, maxValueSize int) (int, int) {
	headerSize := getFragmentHeaderSize(eventSize, maxValueSize)
	valueDataSize := maxValueSize - headerSize // remaining space left in value for data after subtracting the header size
	totalNumValues := math.Ceil(float64(eventSize) / float64(valueDataSize))
	return int(totalNumValues), valueDataSize
}

/*
Main benchmark function.
*/
func main() {
	var batchSize int
	var numEvents int
	var valueSize int
	var workload string

	cmd := &cobra.Command{
		Use:   "fdb-benchmark",
		Short: "FoundationDB benchmark for Seguro.",
		Long:  docString,
		RunE: func(cmd *cobra.Command, args []string) error {
			// Different API versions may expose different runtime behaviors.
			fdb.MustAPIVersion(710)

			// Open the default database from the host system's cluster.
			db := fdb.MustOpenDefault()

			// Generate some mock events without fragmentation.
			events := generateEvents(valueSize, numEvents)

			var errors = 0
			max := 0.0
			min := 999.9
			sum := 0.0
			var txStart time.Time

			switch workload {
			case "single":
				for i, event := range events {
					txStart = time.Now()
					_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
						// Fragment raw mock event before writing to the database.
						fragments := fragmentEvent(event)
						for j, e := range fragments {
							tr.Set(fdb.Key(fmt.Sprint(i+j)), e)
						}
						return
					})
					if err == nil {
						txElapsed := float64(time.Since(txStart).Nanoseconds()) / 1000000.0
						sum = sum + txElapsed
						divisor := float64(batchSize)
						if (txElapsed / divisor) > max {
							max = txElapsed / divisor
						}
						if (txElapsed / divisor) < min {
							min = txElapsed / divisor
						}
					} else {
						errors++
					}
				}
			case "batch":
				for i := 0; i < len(events)-batchSize; i += batchSize {
					txStart = time.Now()
					_, err := db.Transact(func(tr fdb.Transaction) (ret interface{}, err error) {
						for j := 0; j < batchSize; j++ {
							// Fragment raw mock event before writing to the database.
							fragments := fragmentEvent(events[j])
							for k := range fragments {
								tr.Set(fdb.Key(fmt.Sprint(i+j+k)), fragments[k])
							}
						}
						return
					})
					if err == nil {
						txElapsed := float64(time.Since(txStart).Nanoseconds()) / 1000000.0
						sum = sum + txElapsed
						divisor := float64(batchSize)
						if (txElapsed / divisor) > max {
							max = txElapsed / divisor
						}
						if (txElapsed / divisor) < min {
							min = txElapsed / divisor
						}
					} else {
						errors++
					}
				}
			}

			// Print results.
			now := time.Now()
			avg := float64(sum) / float64(len(events))
			fmt.Printf("%s\n", now)
			fmt.Printf("workload: %s\n", workload)
			if workload == "batch" {
				fmt.Printf("batch size: %d\n", batchSize)
			}
			fmt.Printf("event size %d bytes\n", valueSize)
			fmt.Printf("max db value size %d bytes\n", maxValueSize)
			fmt.Printf("n %d\n", len(events))
			fmt.Printf("n_err %d\n", errors)
			fmt.Printf("avg [ms] %f\n", avg)
			fmt.Printf("max [ms] %f\n", max)
			fmt.Printf("min [ms] %f\n", min)
			return nil
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&workload, "workload", "w", defaultWorkload, "The workload to run: \"single\", or \"batch\".")
	flags.IntVar(&valueSize, "value-size", defaultValueSize, "Size of the mock event values in bytes.")
	flags.IntVar(&numEvents, "events", defaultNumEvents, "Number of events to write to the database.")
	flags.IntVar(&batchSize, "batch-size", defaultBatchSize, "Number of events per batch.")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
