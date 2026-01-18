package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strconv"
	"time"

	"github.com/pr0ph0z/uniqlo-sale/pkg"
	"github.com/pr0ph0z/uniqlo-sale/shared"
	"github.com/rs/zerolog/log"
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
		log.Err(err).Msgf("error marshalling media")
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
			log.Err(err).Msgf("error opening file")
			return err
		}
		defer file.Close()

		part, err = multipartWriter.CreateFormFile(fmt.Sprintf("file%d", i), media.File)
		if err != nil {
			log.Err(err).Msgf("error creating form file")
			return err
		}

		_, err = io.Copy(part, file)
		if err != nil {
			log.Err(err).Msgf("error copying file")
			return err
		}
	}

	multipartWriter.Close()

	req, err := http.NewRequest("POST", baseURL, &requestBody)
	if err != nil {
		log.Err(err).Msgf("error creating request")
		return
	}
	req.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msgf("error sending request")
		return
	}

	if resp.StatusCode != http.StatusOK {
		var bodyBytes []byte
		bodyBytes, err = io.ReadAll(resp.Body)
		if err != nil {
			log.Err(err).Msgf("error reading response body")
			return
		}
		bodyString := string(bodyBytes)
		err = errors.New(bodyString)
		return
	}

	return
}

func GetProducts() (products []shared.Product, err error) {
	const limit = 36
	baseURLTemplate := "https://www.uniqlo.com/id/api/commerce/v5/id/products?path=15119%%2C%%2C%%2C&flagCodes=discount&genderId=15119&offset=%d&limit=36&imageRatio=3x4&httpFailure=true"

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Err(err).Msgf("error creating cookie jar")
		return
	}
	client := &http.Client{
		Jar: jar,
	}

	firstURL := fmt.Sprintf(baseURLTemplate, 0)
	req, err := http.NewRequest("GET", firstURL, nil)
	if err != nil {
		log.Err(err).Msgf("error creating request")
		return
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,id-ID;q=0.8,id;q=0.7")
	req.Header.Set("Referer", "https://www.uniqlo.com/id/id/feature/sale/men")
	req.Header.Set("x-fr-clientid", "uq.id.web-spa")
	req.Header.Set("sec-ch-ua", `"Google Chrome";v="141", "Not?A_Brand";v="8", "Chromium";v="141"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("DNT", "1")

	resp, err := client.Do(req)
	if err != nil {
		log.Err(err).Msgf("error sending request")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Err(err).Msgf("unexpected status code %d", resp.StatusCode)
		return
	}

	var productResponse shared.ProductResponse
	err = json.NewDecoder(resp.Body).Decode(&productResponse)
	if err != nil {
		log.Err(err).Msgf("error decoding response")
		return
	}

	for _, product := range productResponse.Result.Items {
		imageURL := ""
		for _, img := range product.Images.Images {
			if img.URL != "" {
				imageURL = img.URL
				break
			}
		}
		if imageURL != "" {
			products = append(products, shared.Product{
				ProductID:       product.ProductID,
				Name:            product.Name,
				ImageURL:        imageURL,
				BasePrice:       strconv.Itoa(product.Prices.Base.Value),
				DiscountedPrice: strconv.Itoa(product.Prices.Promo.Value),
			})
		}
	}

	log.Info().Msgf("total products: %d", productResponse.Result.Pagination.Total)
	total := productResponse.Result.Pagination.Total
	totalLoops := total / limit

	for i := 1; i <= totalLoops; i++ {
		time.Sleep(1 * time.Second)

		offset := i * limit
		url := fmt.Sprintf(baseURLTemplate, offset)

		req, err = http.NewRequest("GET", url, nil)
		if err != nil {
			log.Err(err).Msgf("error creating request")
			return
		}
		req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/141.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "*/*")
		req.Header.Set("Accept-Language", "en-US,en;q=0.9,id-ID;q=0.8,id;q=0.7")
		req.Header.Set("Referer", "https://www.uniqlo.com/id/id/feature/sale/men")
		req.Header.Set("x-fr-clientid", "uq.id.web-spa")
		req.Header.Set("sec-ch-ua", `"Google Chrome";v="141", "Not?A_Brand";v="8", "Chromium";v="141"`)
		req.Header.Set("sec-ch-ua-mobile", "?0")
		req.Header.Set("sec-ch-ua-platform", `"macOS"`)
		req.Header.Set("sec-fetch-dest", "empty")
		req.Header.Set("sec-fetch-mode", "cors")
		req.Header.Set("sec-fetch-site", "same-origin")
		req.Header.Set("DNT", "1")

		resp, err = client.Do(req)
		if err != nil {
			log.Err(err).Msgf("error creating request")
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			log.Err(err).Msgf("unexpected status code %d at offset %d", resp.StatusCode, offset)
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&productResponse)
		if err != nil {
			log.Err(err).Msgf("error decoding response")
			return
		}
		for _, product := range productResponse.Result.Items {
			imageURL := ""
			for _, img := range product.Images.Images {
				if img.URL != "" {
					imageURL = img.URL
					break
				}
			}
			if imageURL != "" {
				products = append(products, shared.Product{
					ProductID:       product.ProductID,
					Name:            product.Name,
					ImageURL:        imageURL,
					BasePrice:       strconv.Itoa(product.Prices.Base.Value),
					DiscountedPrice: strconv.Itoa(product.Prices.Promo.Value),
				})
			} else {
				log.Warn().Msgf("no image url found for product %s", product.ProductID)
				continue
			}
		}
	}

	return
}

func Process(products []shared.Product) (err error) {
	chunkedProducts := chunkBy(products, len(products))
	if len(products) > 10 {
		chunkedProducts = chunkBy(products, 10)
	}
	for _, products = range chunkedProducts {
		var medias []Media
		for i, product := range products {
			var (
				img  pkg.Image
				path string
			)
			img, err = pkg.Download(fmt.Sprintf("%s?width=800", product.ImageURL))
			if err != nil {
				log.Err(err).Msgf("error downloading image")
				return
			}
			img.Name = product.Name
			img.Price, err = pkg.RemoveZeroes(product.BasePrice)
			if err != nil {
				log.Err(err).Msgf("error removing zeroes from base price")
				return
			}
			img.DiscountedPrice, err = pkg.RemoveZeroes(product.DiscountedPrice)
			if err != nil {
				log.Err(err).Msgf("error removing zeroes from discounted price")
				return
			}
			img.DiscountedPercentage = int(math.Round(float64(img.Price-img.DiscountedPrice) / float64(img.Price) * 100))

			path, err = img.PutPrice()
			if err != nil {
				log.Err(err).Msgf("error putting price")
				return
			}
			medias = append(medias, Media{
				Type:    "photo",
				Media:   "attach://file" + strconv.Itoa(i),
				File:    path,
				Caption: fmt.Sprintf("%s\nhttps://www.uniqlo.com/id/id/products/%s", product.Name, product.ProductID),
			})
		}

		err = SendMediaGroup(medias)
		if err != nil {
			return
		}
	}

	return
}

func chunkBy[T any](items []T, chunkSize int) (chunks [][]T) {
	for chunkSize < len(items) {
		items, chunks = items[chunkSize:], append(chunks, items[0:chunkSize:chunkSize])
	}
	return append(chunks, items)
}
