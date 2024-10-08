package model

import (
	"github.com/shopspring/decimal"
	"time"
)

type RestaurantReviewKakaoStatistics struct {
	AverageScore decimal.Decimal `json:"averageScore"`
	Count        int64           `json:"count"`
}

type RestaurantReviewNaverStatistics struct {
	AverageScore decimal.Decimal `json:"averageScore"`
	Count        int64           `json:"count"`
}

type RestaurantReviewStatistics struct {
	Kakao *RestaurantReviewKakaoStatistics `json:"kakao"`
	Naver *RestaurantReviewNaverStatistics `json:"naver"`
}

type RestaurantReviewItem struct {
	RestaurantReviewID int64               `json:"restaurantReviewId"`
	WriterName         string              `json:"writerName"`
	Score              decimal.NullDecimal `json:"score"`
	Content            *string             `json:"content"`
	WroteAt            time.Time           `json:"wroteAt"`
}

type RestaurantReview struct {
	Statistics *RestaurantReviewStatistics `json:"statistics"`
	Reviews    []*RestaurantReviewItem     `json:"reviews"`
	TotalCount int64                       `json:"totalCount"`
}
