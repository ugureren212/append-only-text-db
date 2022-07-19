package main

import (
	"bufio"
	"bytes"
	"log"
	"os"
	"strings"
	"testing"
)

// TODO: Add test data file
// TODO: Change bench mark measurements from per nano seconds to bytes seconds per
// TODO: Implement a flag to change path for server.go

// TODO: Implement TestReadFileBackwards()
func TestReadFileBackwards(t *testing.T) {
}

// Creates a temp log file and appends key value pair
// Reads key value pair to see if it is appended to file
func TestAppendToFile(t *testing.T) {
	t.Parallel()

	appendToFile("test-key", "test-value", false)

	// Create new file
	file, err := os.Open("/tmp/text.log")
	if err != nil {
		t.Error("Cant access /tmp/text.log")
	}

	// Reads file to check key value pair
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		items := bytes.Split([]byte(scanner.Text()), []byte(","))
		if string(items[0]) != "test-key" {
			t.Errorf("key value pair is not correct. Got %v %v", string(items[0]), string(items[1]))
		}
		if string(items[1]) != "test-value" {
			t.Errorf("key value pair is not correct. Got %v %v", string(items[0]), string(items[1]))
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// TODO: Implement fuzzing
// Append 10 char key value
func BenchmarkAppendToFile10(b *testing.B) {
	val := []byte(strings.Repeat("a", 10))
	for i := 0; i < b.N; i++ {
		appendToFile("test-key", string(val), false)
	}
}

// Append 1000 char key value

func BenchmarkAppendToFile1000(b *testing.B) {
	val := []byte(strings.Repeat("a", 1000))
	for i := 0; i < b.N; i++ {
		appendToFile("test-key", string(val), false)
	}
}

// Append 100000 char key value

func BenchmarkAppendToFile100000(b *testing.B) {
	val := []byte(strings.Repeat("a", 100000))
	for i := 0; i < b.N; i++ {
		appendToFile("test-key", string(val), false)
	}
}



// go test -run=XXX -bench=.
