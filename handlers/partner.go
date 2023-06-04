package handlers

import (
	"encoding/json"
	"net/http"
	"time"
	. "wigglebackend/helpers"
	. "wigglebackend/models"
)

func AddPartnerHandler(w http.ResponseWriter, r *http.Request) {
	var partner Partner
	var jsonData map[string]interface{}
	errorDecoder := json.NewDecoder(r.Body).Decode(&jsonData)
	locale := r.Header.Get("Accept-Language")
	token := r.Header.Get("Authorization")

	CheckError(errorDecoder)

	relationUserId, err := ExtractUserId(token)
	name, err := GetJSONField(jsonData, "name")
	gender, err := GetJSONField(jsonData, "gender")
	ageRange, err := GetJSONField(jsonData, "ageRange")
	photoUrl, err := GetJSONField(jsonData, "photoUrl")

	partner.RelationUserId = relationUserId
	partner.Name = name
	partner.Gender = gender
	partner.AgeRange = ageRange
	partner.PhotoUrl = photoUrl
	partner.CreatedAt = time.Now().String()
	partner.UpdatedAt = time.Now().String()

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Partner add error 0001")))
		return
	}
	createdAtFormatted := FormatTime(partner.CreatedAt)
	updatedAtFormatted := FormatTime(partner.UpdatedAt)

	if partner, success := InsertPartner(partner); success {
		response := map[string]interface{}{
			"relationUserId": partner.RelationUserId,
			"name":           partner.Name,
			"gender":         partner.Gender,
			"ageRange":       partner.AgeRange,
			"photoUrl":       partner.PhotoUrl,
			"createdAt":      createdAtFormatted,
			"updatedAt":      updatedAtFormatted,
		}
		jsonResponseData, err := json.Marshal(response)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(Localizate(locale, "Partner add error 0002")))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponseData)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Partner add error 0003")))
		return
	}

}
