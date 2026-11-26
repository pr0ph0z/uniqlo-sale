package pkg

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/leekchan/accounting"

	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/pr0ph0z/uniqlo-sale/shared"
)

const FontFileName = "UniqloProBold.ttf"

type Image struct {
	Name                 string
	SaleColors           []shared.SaleColor
	SaleSizes            []string
	Price                int
	DiscountedPrice      int
	DiscountedPercentage int
	Path                 string
	Width                int
	Height               int
}

func Download(url string) (img Image, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	out, err := os.CreateTemp("", "*.jpg")
	if err != nil {
		return
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return
	}
	img.Path = out.Name()

	_, err = out.Seek(0, 0)
	if err != nil {
		return
	}

	im, _, err := image.DecodeConfig(out)
	if err != nil {
		return
	}
	img.Width = im.Width
	img.Height = im.Height

	return
}

func (img Image) PutPrice() (path string, err error) {
	im, err := gg.LoadImage(img.Path)
	if err != nil {
		return
	}

	dir, err := os.Getwd()
	if err != nil {
		return
	}

	fontBytes, err := os.ReadFile(filepath.Join(dir, "assets", FontFileName))
	if err != nil {
		return
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return
	}

	const padding = 16
	titleFace := truetype.NewFace(font, &truetype.Options{Size: 50})
	detailFace := truetype.NewFace(font, &truetype.Options{Size: 30})
	measure := gg.NewContext(1, 1)
	width := float64(img.Width - 2*padding)
	measure.SetFontFace(titleFace)
	title := measure.WordWrap(img.Name, width)
	measure.SetFontFace(detailFace)
	var sizes []string
	if len(img.SaleSizes) > 0 {
		sizes = measure.WordWrap("Sale sizes: "+strings.Join(img.SaleSizes, ", "), width)
	}

	colorHeight := 0
	if len(img.SaleColors) > 0 {
		colorHeight = 40 + 44*len(img.SaleColors)
	}
	height := img.Height + 2*padding + 58*len(title) + colorHeight + 38*len(sizes) + 105
	context := gg.NewContext(img.Width, height)
	context.SetHexColor("#FFFFFF")
	context.Clear()
	context.DrawImage(im, 0, 0)
	context.SetHexColor("#000000")
	y := img.Height + padding
	context.SetFontFace(titleFace)
	for _, line := range title {
		context.DrawString(line, padding, float64(y+45))
		y += 58
	}
	context.SetFontFace(detailFace)
	if len(img.SaleColors) > 0 {
		context.DrawString("Sale colors:", padding, float64(y+30))
		y += 40
		for _, color := range img.SaleColors {
			x := float64(padding)
			if swatch := loadSwatch(color.SwatchURL); swatch != nil {
				context.Push()
				context.Translate(x, float64(y))
				context.Scale(36/float64(swatch.Bounds().Dx()), 36/float64(swatch.Bounds().Dy()))
				context.DrawImage(swatch, 0, 0)
				context.Pop()
				x += 46
			}
			context.DrawString(color.Name, x, float64(y+29))
			y += 44
		}
	}
	for _, line := range sizes {
		context.DrawString(line, padding, float64(y+30))
		y += 38
	}

	ac := accounting.Accounting{Symbol: "Rp", Precision: 0,
		Format: "%s%v", Thousand: "."}
	drawTextWithStrikethrough(context, padding, float64(y+35), ac.FormatMoney(img.Price))

	context.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 55}))
	context.SetHexColor("#FF0000")
	context.DrawString(fmt.Sprintf("%s (%d%%)", ac.FormatMoney(img.DiscountedPrice), img.DiscountedPercentage), padding, float64(y+92))

	path = fmt.Sprintf("%s.png", filepath.Join(os.TempDir(), RandomString(10)))
	err = context.SavePNG(path)
	if err != nil {
		return
	}

	return
}

func loadSwatch(rawURL string) image.Image {
	u, err := url.Parse(rawURL)
	if err != nil || u.Scheme != "https" || u.Host != "image.uniqlo.com" {
		return nil
	}
	query := u.Query()
	query.Set("width", "36")
	u.RawQuery = query.Encode()
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(u.String())
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil
	}
	swatch, _, err := image.Decode(resp.Body)
	if err != nil || swatch.Bounds().Empty() {
		return nil
	}
	return swatch
}

func drawTextWithStrikethrough(dc *gg.Context, x, y float64, text string) {
	// Draw the original text
	dc.DrawString(text, x, y)

	// Get the width and height of the text
	w, h := dc.MeasureString(text)

	// Calculate the y-position for the strikethrough line
	// (typically around 40% from the top of the text)
	strikethroughY := y - h*0.4

	// Draw the strikethrough line
	dc.DrawLine(x, strikethroughY, x+w+5, strikethroughY)
	dc.Stroke()
}
