package main

type Anime struct {
	Id     string `dynamodbav:"id" json:"id"`
	Title  string `dynamodbav:"title" json:"title"`
	Author string `dynamodbav:"author" json:"author"`
	Year   int    `dynamodbav:"year" json:"year"`
	Status string `dynamodbav:"status" json:"status"`
}

type CreateAnimeRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Status string `json:"status"`
}
