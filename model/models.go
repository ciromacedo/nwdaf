package model

import "time"

type RegistrationAccept struct {
	Date time.Time `json:"date"`
	Amf Amf
	Ue	Ue
	PlmnId PlmnId

}

type Amf struct {
	Id string  `json:"id"`
	Locale string  `json:"locale"`
}

type Ue struct {
	Suci string  `json:"suci"`
	Supi string  `json:"supi"`
}

type PlmnId struct {
	Mcc string  `json:"mcc"`
	Mnc string  `json:"mnc"`
}

/* CONFIG */
type Config struct{
	Port int
	MongoURI string
	DBName string
}

type CollectionInfo struct {
	DocumentName string  `json:"Name"`
	NumberOfRecords int64
}

type Article struct {
	Title string `json:"Title"`
	Desc string `json:"desc"`
	Content string `json:"content"`
}