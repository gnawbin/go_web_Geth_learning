package response

// @Author: morris
type MetricsItem struct {
	Title string `json:"title"`
	Label string `json:"label"`
	Value string `json:"value"`
	// r,g,b 格式 例如：255,255,0
	Color string `json:"color"`
	Logo  string `json:"logo"`
}
