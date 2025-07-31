package main

import (
	"reflect"
	"sort"
	"testing"
	"time"
)

func makeChan(vals []int, delay time.Duration) <-chan int {
	ch := make(chan int)
	go func() {
		defer close(ch)
		for _, v := range vals {
			time.Sleep(delay)
			ch <- v
		}
	}()
	return ch
}

func TestMerge(t *testing.T) {
	ch1 := makeChan([]int{1, 2}, 5*time.Millisecond)
	ch2 := makeChan([]int{3, 4}, 5*time.Millisecond)

	out := Merge(ch1, ch2)
	var result []int
	for v := range out {
		result = append(result, v)
	}

	sort.Ints(result)
	expected := []int{1, 2, 3, 4}

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("expected %v, got %v", expected, result)
	}
}
