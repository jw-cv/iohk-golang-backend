package model

import (
	"iohk-golang-backend/ent/customer"
	"time"
)

type Gender string

const (
	GenderMale   Gender = "MALE"
	GenderFemale Gender = "FEMALE"
)

type Customer struct {
	ID         int
	Name       string
	Surname    string
	Number     int
	Gender     Gender
	Country    string
	Dependants int
	BirthDate  time.Time
}

func (g Gender) ToEntGender() customer.Gender {
	switch g {
	case GenderMale:
		return customer.GenderMale
	case GenderFemale:
		return customer.GenderFemale
	default:
		return customer.Gender("")
	}
}

func (g Gender) ToDatabaseValue() string {
	switch g {
	case GenderMale:
		return "Male"
	case GenderFemale:
		return "Female"
	default:
		return ""
	}
}
