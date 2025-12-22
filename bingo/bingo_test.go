package bingo

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestGenerateIndex(t *testing.T) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	max := 150

	indexes := generateIndex(r, max)

	if len(indexes) != 25 {
		t.Errorf("number of indexes is not correct")
	}

	seen := make(map[int]bool)
	for _, num := range indexes {
		if seen[num] {
			t.Errorf("duplicated value: %d", num)
		}
		seen[num] = true
	}

	fmt.Println("Generated indexes:", indexes)

}
