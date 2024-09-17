package shared

type ProductResponse struct {
	Result Result `json:"result"`
}

type Result struct {
	Items      []Item     `json:"items"`
	Pagination Pagination `json:"pagination"`
}

type Item struct {
	Images    Images `json:"images"`
	Name      string `json:"name"`
	Prices    Prices `json:"prices"`
	ProductID string `json:"productId"`
}

type Images struct {
	Images []Image `json:"main"`
}

type Image struct {
	URL string `json:"url"`
}

type Prices struct {
	Base  Price `json:"base"`
	Promo Price `json:"promo"`
}

type Price struct {
	Value string `json:"value"`
}

type Pagination struct {
	Count  int `json:"count"`
	Total  int `json:"total"`
	Offset int `json:"offset"`
}

type Product struct {
	ProductID       string `json:"product_id"`
	Name            string `json:"name"`
	ImageURL        string `json:"image_url"`
	BasePrice       string `json:"base_price"`
	DiscountedPrice string `json:"discounted_price"`
}
