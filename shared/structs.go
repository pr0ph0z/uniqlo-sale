package shared

type ProductResponse struct {
	Result Result `json:"result"`
}

type Result struct {
	Items      []Item     `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Item struct {
	Images                         Images   `json:"images"`
	Name                           string   `json:"name"`
	Prices                         Prices   `json:"prices"`
	ProductID                      string   `json:"productId"`
	PriceGroup                     string   `json:"priceGroup"`
	RepresentativeColorDisplayCode string   `json:"representativeColorDisplayCode"`
	Colors                         []Option `json:"colors"`
	Sizes                          []Option `json:"sizes"`
}

type Option struct {
	Name        string `json:"name"`
	DisplayCode string `json:"displayCode"`
}

type Images struct {
	Images []Image           `json:"sub"`
	Chips  map[string]string `json:"chip"`
}

type Image struct {
	URL string `json:"image"`
}

type Prices struct {
	Base  Price `json:"base"`
	Promo Price `json:"promo"`
}

type Price struct {
	Value int `json:"value"`
}

type Pagination struct {
	Count  int `json:"count"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
}

type Product struct {
	ProductID       string      `json:"product_id"`
	ProductURL      string      `json:"product_url"`
	Name            string      `json:"name"`
	ImageURL        string      `json:"image_url"`
	BasePrice       string      `json:"base_price"`
	DiscountedPrice string      `json:"discounted_price"`
	SaleColors      []SaleColor `json:"sale_colors,omitempty"`
	SaleSizes       []string    `json:"sale_sizes,omitempty"`
}

type SaleColor struct {
	Name      string `json:"name"`
	SwatchURL string `json:"swatch_url,omitempty"`
}
