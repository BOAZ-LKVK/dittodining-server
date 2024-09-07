package model

import "time"

type RestaurantReviewKakaoStatistics struct {
	AverageScore float32 `json:"averageScore"`
	Count        int64   `json:"count"`
}

type RestaurantReviewNaverStatistics struct {
	AverageScore float32 `json:"averageScore"`
	Count        int64   `json:"count"`
}

type RestaurantReviewStatistics struct {
	Kakao *RestaurantReviewKakaoStatistics `json:"kakao"`
	Naver *RestaurantReviewNaverStatistics `json:"naver"`
}

type RestaurantReviewItem struct {
	RestaurantReviewID int64     `json:"restaurantReviewId"`
	WriterName         string    `json:"writerName"`
	Score              float32   `json:"score"`
	Content            string    `json:"content"`
	WrittenAt          time.Time `json:"writtenAt"`
}

type RestaurantReview struct {
	Statistics *RestaurantReviewStatistics `json:"statistics"`
	Reviews    []RestaurantReviewItem      `json:"reviews"`
}
