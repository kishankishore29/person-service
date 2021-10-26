package models

type Person struct {
	Id      string `gorm:"default:uuid_generate_v4()"` //ID of the person-record
	Name    string // Name of the person
	Age     int32  // Age of the Person
	Email   string // Email address of the person
	Country string // Country of the person
}
