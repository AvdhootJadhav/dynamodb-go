package main

type Anime struct {
	Id     string `dynamodbav:"id"`
	Title  string `dynamodbav:"title"`
	Author string `dynamodbav:"author"`
	Year   int    `dynamodbav:"year"`
	Status string `dynamodbav:"status"`
}

type CreateAnimeRequest struct {
	Title  string `json:"title"`
	Author string `json:"author"`
	Year   int    `json:"year"`
	Status string `json:"status"`
}
