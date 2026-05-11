package models

type Product struct {
	ID, Description string
}

type ProductFeatures struct {
	ID       string   `json:"id"`
	Name     string   `json:"name"`
	Brand    string   `json:"brand"`
	Category string   `json:"category"`
	Price    float64  `json:"price_rub"`
	Currency string   `json:"currency"`
	KeySpecs []string `json:"key_specs"`
}

type PipelineResult struct {
	Model      string            `json:"model"`
	TotalItems int               `json:"total_items"`
	Failed     int               `json:"failed"`
	Products   []ProductFeatures `json:"products"`
}
