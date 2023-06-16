package models

import (
	. "wigglebackend/helpers"
)

type Partner struct {
	RelationUserId string `json:"RelationUserId"`
	Name           string `json:"Name"`
	Gender         string `json:"Gender"`
	AgeRange       string `json:"AgeRange"`
	PhotoUrl       string `json:"PhotoUrl"`
	CreatedAt      string `json:"CreatedAt"`
	UpdatedAt      string `json:"UpdatedAt"`
}

func InsertPartner(partner Partner) (Partner, bool) {
	sqlQuery := `
	INSERT INTO partner (
		RelationUserId, Name, Gender, AgeRange, PhotoUrl, CreatedAt, UpdatedAt
	) VALUES (?,?,?,?,?,?,?)
`

	result, err := db.Exec(sqlQuery,
		partner.RelationUserId, partner.Name, partner.Gender, partner.AgeRange,
		partner.PhotoUrl, partner.CreatedAt, partner.UpdatedAt,
	)
	if err != nil {
		CheckError(err)
		return Partner{}, false
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		CheckError(err)
		return Partner{}, false
	}

	if rowsAffected > 0 {
		return partner, true
	} else {
		return Partner{}, false
	}
}

func UpdatePartner(partner Partner) bool {
	//test func
	return true
}
