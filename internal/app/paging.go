package app

type Page struct {
	Items  any `json:"items"`
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
	Total  int `json:"total"`
}

func normalizePage(limit, offset int) (int, int) {
	if limit <= 0 || limit > 250 {
		limit = 50
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
