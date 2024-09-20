package pkg

import (
	"encoding/json"
	"github.com/pr0ph0z/uniqlo-sale/shared"
	"github.com/rs/zerolog/log"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
var LastFetchedFileName = "last-fetch.json"

type LastFetch struct {
	TotalProducts int              `json:"total_products"`
	Hash          uint64           `json:"hash"`
	Products      []shared.Product `json:"products"`
	FetchedAt     time.Time        `json:"fetched_at"`
}

func RandomString(length int) (text string) {
	b := make([]rune, length)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

func RemoveZeroes(value string) (ret int, err error) {
	intPart := strings.Split(value, ".")[0]
	ret, err = strconv.Atoi(intPart)
	return
}

func LastFetchedItems() (lastFetch LastFetch, err error) {
	jsonFile, err := os.Open("last-fetch.json")
	if err != nil {
		return
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)
	err = json.Unmarshal(byteValue, &lastFetch)
	if err != nil {
		return
	}

	return
}

func WriteLastFetchedItems(lastFetchedItems LastFetch) (err error) {
	newFetchedItems, err := json.Marshal(lastFetchedItems)
	if err != nil {
		log.Err(err).Send()
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		return
	}
	err = os.WriteFile(filepath.Join(wd, LastFetchedFileName), newFetchedItems, 0644)
	if err != nil {
		return
	}

	return
}
