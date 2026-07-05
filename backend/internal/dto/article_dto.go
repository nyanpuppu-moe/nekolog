package dto

type ArticlePostRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type ArticlePatchRequest struct {
	Content string `json:"content"`
}
