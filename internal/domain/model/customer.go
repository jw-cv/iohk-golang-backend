package model

import (
	"iohk-golang-backend-preprod/ent/customer"
	"time"
)

type Gender string

const (
	GenderMale   Gender = "Male"
	GenderFemale Gender = "Female"
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
		return customer.Gender("") // or handle this case as you see fit
	}
}
