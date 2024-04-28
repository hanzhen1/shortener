// Code generated by goctl. DO NOT EDIT.
package types

type ConvertRequest struct {
	LongURL string `json:"longUrl" validate:"required"`
}

type ConvertResponse struct {
	ShortURL string `json:"shortUrl"`
}

type ShowRequest struct {
	ShortURL string `path:"shortUrl" validate:"required"`
}

type ShowResponse struct {
	LongURL string `json:"longUrl"`
}
