package domain

type Product struct {
	ID       uint64  `gorm:"column:id;primary_key" json:"id"`
	ItemName string  `json:"description"`
	Lat      float64 `json:"lat"`
	Lng      float64 `json:"lng"`
	ImageURL string  `json:"img_URL"`
	URL      string  `json:"product_URL"`
}

func (p *Product) TableName() string {
	return "items"
}

type Query struct {
	Term   string  `json:"term"`
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Radius float64 `json:"radius"`
}

type Point struct {
	X float64
	Y float64
}
