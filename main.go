package main

import (
	"github.com/cespare/xxhash"
	"github.com/pr0ph0z/uniqlo-sale/internal"
	"github.com/pr0ph0z/uniqlo-sale/pkg"
	"github.com/pr0ph0z/uniqlo-sale/shared"
	zlogger "github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"
	"time"
)

func main() {
	zlogger.ErrorStackMarshaler = pkgerrors.MarshalStack
	log := zlogger.New(os.Stdout).With().Caller().Logger()
	log.Info().Msg("starting the program")

	lastFetchedItems, err := pkg.LastFetchedItems()
	if err != nil {
		log.Err(err).Send()
		return
	}

	products, err := internal.GetProducts()
	if err != nil {
		log.Err(err).Send()
		return
	}

	lastFetchedItemsSet := make(map[string]struct{}, len(lastFetchedItems.Products))
	for _, v := range lastFetchedItems.Products {
		lastFetchedItemsSet[v.ProductID] = struct{}{}
	}

	var (
		productIDs       []string
		filteredProducts []shared.Product
	)
	for _, product := range products {
		productIDs = append(productIDs, product.ProductID)
		if _, exists := lastFetchedItemsSet[product.ProductID]; !exists {
			filteredProducts = append(filteredProducts, product)
		}
	}

	hash := xxhash.Sum64String(strings.Join(productIDs, ""))
	if lastFetchedItems.TotalProducts > len(products) || lastFetchedItems.Hash == hash {
		log.Warn().Msg("no updates on the products")
		return
	}
	lastFetchedItems.TotalProducts = len(products)
	lastFetchedItems.Hash = hash

	log.Info().Msg("processing")
	err = internal.Process(filteredProducts)
	if err != nil {
		log.Err(err).Send()
		return
	}

	lastFetchedItems.Products = products
	lastFetchedItems.FetchedAt = time.Now()

	err = pkg.WriteLastFetchedItems(lastFetchedItems)
	if err != nil {
		log.Err(err).Send()
		return
	}

	log.Info().Msgf("finished with %d new products", len(filteredProducts))
}
