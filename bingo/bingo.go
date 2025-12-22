package bingo

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"log"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

type BingoData struct {
	Goals []BingoGoal `json:"goals"`
}

type BingoGoal struct {
	Name       string `json:"name"`
	Difficulty int    `json:"difficulty"`
}

type BingoCard struct {
	Goals [25]BingoGoal `json:"goals"`
	Seed  string        `json:"seed"`
}

var (
	data    *BingoData
	once    sync.Once
	loadErr error
)

func InitData(filepath string) error {
	once.Do(func() {
		file, err := os.ReadFile(filepath)
		if err != nil {
			loadErr = fmt.Errorf("failed to read file: %w", err)
			return
		}

		var bd BingoData
		if err := json.Unmarshal(file, &bd); err != nil {
			loadErr = fmt.Errorf("failed to unmarshal json: %w", err)
			return
		}

		if len(bd.Goals) < 25 {
			loadErr = fmt.Errorf("insufficient goals: needs at least 25, got %d", len(bd.Goals))
			return
		}

		data = &bd
		log.Printf("Bingo data successfully loaded.")
	})
	return loadErr
}

func GetBingoData() *BingoData {
	return data
}

func (bc BingoCard) GetSeed() string {
	return bc.Seed
}

func CreateBingoCard(seed string) (BingoCard, error) {
	var bc BingoCard
	bc.Seed = seed

	// generate bingo card
	if err := generate(&bc); err != nil {
		return bc, fmt.Errorf("failed to create bingo card: %w", err)
	}

	return bc, nil
}

func getRNG(bc *BingoCard) *rand.Rand {
	var seedInt64 int64

	if len(bc.Seed) != 0 {
		seedInt64 = strToInt64(bc.Seed)
	} else {
		seedInt64 = time.Now().UnixNano()
		bc.Seed = strconv.FormatInt(seedInt64, 10)
	}

	return rand.New(rand.NewSource(seedInt64))
}

func generateIndex(r *rand.Rand, max int) [25]int {
	indexes := make([]int, 0, 25)
	indexSet := make(map[int]bool)

	for len(indexes) < 25 {
		index := r.Intn(max)

		if !indexSet[index] {
			indexSet[index] = true
			indexes = append(indexes, index)
		}
	}

	return [25]int(indexes)
}

func generate(bc *BingoCard) error {
	masterData := GetBingoData()

	if masterData == nil {
		return fmt.Errorf("bingo data not initialized")
	}

	r := getRNG(bc)
	indexes := generateIndex(r, len(masterData.Goals))

	for i, index := range indexes {
		goal := masterData.Goals[index]
		bc.Goals[i] = goal
	}

	log.Printf("Bingo card created. Seed: %v", bc.Seed)
	return nil
}

func strToInt64(s string) int64 {
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}
	h := fnv.New64a()
	h.Write([]byte(s))
	return int64(h.Sum64())
}
