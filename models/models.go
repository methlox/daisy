package models

import (
	"gorm.io/gorm"
	"time"
	)

type Form struct {
	ID				uint		`gorm:"primary key;autoIncrement" json:"id"`
	Title			*string 	`json:"title"`
	Description		*string 	`json:"desc"`
	Created_at 		time.Time   `json:"created_at"`

}

type Question struct {
	ID				uint		`gorm:"primary key;autoIncrement" json:"id"`
	Form_id			uint 		`json:"form_id"`
	Question_text	*string 	`json:"ques_text"`
	Question_order	*string 	`json:"ques_order"`
	Created_at 		time.Time   `json:"created_at"`
}

type FormResponse struct {
	ID				uint		`gorm:"primary key;autoIncrement" json:"id"`
	Form_id			uint 		`json:"form_id"`
	Responded_at 	time.Time   `json:"responded_at"`
}

type Response struct {
	ID				uint		`gorm:"primary key;autoIncrement" json:"id"`
	Response_id		uint 		`json:"response_id"`
	Question_id		uint 		`json:"ques_id"`
	Answer			*string 	`json:"answer"`
}


func MigrateForm(db *gorm.DB) error {
	err := db.AutoMigrate(&Form{})
	return err
}

func MigrateQuestion(db *gorm.DB) error {
	err := db.AutoMigrate(&Question{})
	return err
}

func MigrateFormResponse(db *gorm.DB) error {
	err := db.AutoMigrate(&FormResponse{})
	return err
}

func MigrateResponse(db *gorm.DB) error {
	err := db.AutoMigrate(&Response{})
	return err
}