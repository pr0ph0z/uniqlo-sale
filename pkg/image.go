package pkg

import (
	"fmt"
	"github.com/fogleman/gg"
	"github.com/golang/freetype/truetype"
	"github.com/leekchan/accounting"
	"image"
	_ "image/jpeg"
	_ "image/png"

	"io"
	"net/http"
	"os"
	"path/filepath"
)

const FontFileName = "UniqloProBold.ttf"

type Image struct {
	Name            string
	Price           int
	DiscountedPrice int
	Path            string
	Width           int
	Height          int
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
	context := gg.NewContext(img.Width, img.Height+250)
	context.SetHexColor("#FFFFFF")
	context.Clear()

	im, err := gg.LoadImage(img.Path)
	if err != nil {
		return
	}
	context.DrawImage(im, 0, 0)

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

	context.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 50}))
	context.SetHexColor("#000000")
	context.DrawStringWrapped(img.Name, 0, float64(img.Height+10), 0, 0, float64(img.Width), float64(1), gg.AlignLeft)

	ac := accounting.Accounting{Symbol: "Rp", Precision: 0,
		Format: "%s%v", Thousand: "."}
	context.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 30}))
	drawTextWithStrikethrough(context, 0, float64(img.Height+160), ac.FormatMoney(img.Price))

	context.SetFontFace(truetype.NewFace(font, &truetype.Options{Size: 55}))
	context.SetHexColor("#FF0000")
	context.DrawStringWrapped(ac.FormatMoney(img.DiscountedPrice), 0, float64(img.Height+170), 0, 0, float64(img.Width), float64(1), gg.AlignLeft)

	path = fmt.Sprintf("%s.png", filepath.Join(os.TempDir(), RandomString(10)))
	err = context.SavePNG(path)
	if err != nil {
		return
	}

	return
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
