package response

type MetricsItem struct {
	Title string `json:"title"`
	Label string `json:"label"`
	Value string `json:"value"`
	Color string `json:"color"` // r,g,b 格式 例如： 255,255,0
	Logo  string `json:"logo"`
}
