package main

import (
	"github.com/cespare/xxhash"
	zlogger "github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"strings"
	"time"
	"uniqlo-sale/internal"
	"uniqlo-sale/pkg"
)

func main() {
	zlogger.ErrorStackMarshaler = pkgerrors.MarshalStack
	log := zlogger.New(os.Stdout).With().Caller().Logger()

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

	hash := xxhash.Sum64String(strings.Join(lastFetchedItems.ProductIDs, ""))
	if lastFetchedItems.TotalProducts == products.Result.Pagination.Total {
		if lastFetchedItems.Hash == hash {
			log.Warn().Msg("no updates on the products")
			return
		}
	}
	lastFetchedItems.TotalProducts = products.Result.Pagination.Total
	lastFetchedItems.Hash = hash

	//err = internal.Process(products)
	//if err != nil {
	//	log.Err(err).Send()
	//	return
	//}

	var productIDs []string
	for _, item := range products.Result.Items {
		productIDs = append(productIDs, item.ProductID)
	}

	lastFetchedItems.ProductIDs = productIDs
	lastFetchedItems.FetchedAt = time.Now()

	err = pkg.WriteLastFetchedItems(lastFetchedItems)
	if err != nil {
		log.Err(err).Send()
		return
	}

	log.Info().Msg("finished")
}
