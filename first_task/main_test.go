package main

import (
	"bytes"
	"io"
	"os"
	"testing"
	"time"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	_ = w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

func TestPrintSorted_OneChannelEmpty(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	// ch2 сразу закрыт
	go func() {
		close(ch2)
	}()

	go func() {
		defer close(ch1)
		ch1 <- 1
		ch1 <- 2
		ch1 <- 3
	}()

	output := captureOutput(func() {
		PrintSorted(ch1, ch2)
	})

	expected := "1\n2\n3\n"
	if output != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
	}
}

func TestPrintSorted_ChannelClosesFirst(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		ch1 <- 1
	}()

	go func() {
		defer close(ch2)
		time.Sleep(20 * time.Millisecond)
		ch2 <- 2
		ch2 <- 3
	}()

	output := captureOutput(func() {
		PrintSorted(ch1, ch2)
	})

	expected := "1\n2\n3\n"
	if output != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
	}
}

func TestPrintSorted_Interleaved(t *testing.T) {
	ch1 := make(chan int)
	ch2 := make(chan int)

	go func() {
		defer close(ch1)
		ch1 <- 1
		time.Sleep(10 * time.Millisecond)
		ch1 <- 4
	}()

	go func() {
		defer close(ch2)
		ch2 <- 2
		time.Sleep(10 * time.Millisecond)
		ch2 <- 3
	}()

	output := captureOutput(func() {
		PrintSorted(ch1, ch2)
	})

	expected := "1\n2\n3\n4\n"
	if output != expected {
		t.Errorf("expected:\n%s\ngot:\n%s", expected, output)
	}
}
