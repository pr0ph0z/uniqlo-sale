package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/pr0ph0z/uniqlo-sale/pkg"
	"github.com/pr0ph0z/uniqlo-sale/shared"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

type Media struct {
	Type    string `json:"type"`
	Media   string `json:"media"`
	File    string `json:"file"`
	Caption string `json:"caption"`
}

type Request struct {
	ChatID string  `json:"chat_id"`
	Media  []Media `json:"media"`
}

func SendMediaGroup(medias []Media) (err error) {
	botToken := os.Getenv("BOT_TOKEN")
	chatID := os.Getenv("CHAT_ID")
	baseURL := fmt.Sprintf("https://api.telegram.org/bot%s/sendMediaGroup", botToken)

	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)
	multipartWriter.WriteField("chat_id", chatID)

	mediaJSON, err := json.Marshal(medias)
	if err != nil {
		return
	}

	multipartWriter.WriteField("media", string(mediaJSON))

	for i, media := range medias {
		var (
			file *os.File
			part io.Writer
		)
		file, err = os.Open(media.File)
		if err != nil {
			return err
		}
		defer file.Close()

		part, err = multipartWriter.CreateFormFile(fmt.Sprintf("file%d", i), media.File)
		if err != nil {
			return err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			return err
		}
	}

	multipartWriter.Close()

	req, err := http.NewRequest("POST", baseURL, &requestBody)
	if err != nil {
		return
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			return
		}
		bodyString := string(bodyBytes)
		err = errors.New(bodyString)
		return
	}

	return
}

func GetProducts() (products []shared.Product, err error) {
	const ProductLimit = 200
	baseURL := "https://www.uniqlo.com/id/api/commerce/v3/en/products?path=15119&flagCodes=discount&offset=0&isV2Review=true&limit=" + strconv.Itoa(ProductLimit)

	req, err := http.NewRequest("GET", baseURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	var productResponse shared.ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&productResponse)
	if err != nil {
		return
	}

	for _, product := range productResponse.Result.Items {
		products = append(products, shared.Product{
			ProductID:       product.ProductID,
			Name:            product.Name,
			ImageURL:        product.Images.Images[0].URL,
			BasePrice:       product.Prices.Base.Value,
			DiscountedPrice: product.Prices.Promo.Value,
		})
	}

	return
}

func Process(products []shared.Product) (err error) {
	var medias []Media
	for i, product := range products {
		var (
			img  pkg.Image
			path string
		)
		img, err = pkg.Download(fmt.Sprintf("%s?width=800", product.ImageURL))
		if err != nil {
			return
		}
		img.Name = product.Name
		img.Price, err = pkg.RemoveZeroes(product.BasePrice)
		if err != nil {
			return
		}
		img.DiscountedPrice, err = pkg.RemoveZeroes(product.DiscountedPrice)
		if err != nil {
			return
		}

		path, err = img.PutPrice()
		if err != nil {
			return
		}
		medias = append(medias, Media{
			Type:    "photo",
			Media:   "attach://file" + strconv.Itoa(i),
			File:    path,
			Caption: "https://www.uniqlo.com/id/id/products/" + product.ProductID,
		})
	}

	err = SendMediaGroup(medias)
	if err != nil {
		return
	}

	return
}
