package entity

import (
	formatter "github.com/wisesight/spider-go-formatter"
)

type Method string

const (
	NER       Method = "ner"
	Sentiment Method = "sentiment"
	TFIDF     Method = "tfidf"
)

type EnrichementMethod struct {
	Method Method
	Post   *formatter.Spider
}
