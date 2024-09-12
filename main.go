package main

import (
	zlogger "github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"uniqlo-sale/internal"
)

func main() {
	zlogger.ErrorStackMarshaler = pkgerrors.MarshalStack
	log := zlogger.New(os.Stdout).With().Caller().Logger()

	products, err := internal.GetProducts()
	if err != nil {
		log.Err(err).Send()
	}

	err = internal.Process(products)
	if err != nil {
		log.Err(err).Send()
	}
}
