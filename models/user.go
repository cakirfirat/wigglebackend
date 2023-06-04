package models

import (
	"database/sql"
	"fmt"
	"strings"
	. "wigglebackend/helpers"
)

type User struct {
	Id         string `json:"Id"`
	Username   string `json:"UserName"`
	Mail       string `json:"Mail"`
	Phone      string `json:"Phone"`
	Password   string `json:"Password"`
	Gender     string `json:"Gender"`
	BirthDate  string `json:"BirthDate"`
	VerifyCode string `json:"VerifyCode"`
	IsVerify   int    `json:"IsVerify"`
	PhotoUrl   string `json:"PhotoUrl"`
	CreatedAt  string `json:"CreatedAt"`
	UpdatedAt  string `json:"UpdatedAt"`
}

func InsertUser(user User) (bool, int64) {
	sqlQuery := `
	INSERT INTO user (
		Username, Mail, Phone, Password, Gender, BirthDate,
		VerifyCode, IsVerify, PhotoUrl, CreatedAt, UpdatedAt
	) VALUES (?,?,?,?,?,?,?,?,?,?,?)
`

	result, err := db.Exec(sqlQuery,
		user.Username, user.Mail, user.Phone, user.Password, user.Gender,
		user.BirthDate, user.VerifyCode, user.IsVerify,
		user.PhotoUrl, user.CreatedAt, user.UpdatedAt,
	)
	if err != nil {
		CheckError(err)
		return false, 0
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		CheckError(err)
		return false, 0
	}

	lastInsertID, err := result.LastInsertId()
	if err != nil {
		CheckError(err)
		return false, 0
	}

	if rowsAffected > 0 {
		return true, lastInsertID
	} else {
		return false, 0
	}
}

func CheckPhoneNumber(phoneNumber string) bool {
	var count int
	fmt.Println(phoneNumber)
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE Phone LIKE CONCAT('%', ?, '%')", phoneNumber).Scan(&count)
	if err != nil {
		CheckError(err)
	}
	if count > 0 {
		return true
	} else {
		return false
	}

}

func CheckVerifyCode(phoneNumber, verifyCode string) bool {
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM user WHERE Phone AND VerifyCode", phoneNumber, verifyCode).Scan(&count)
	if err != nil {
		CheckError(err)
	}
	if count > 0 {
		return false
	} else {
		return true
	}

}

func UpdateUserFromId(Id string, updateFields map[string]interface{}) bool {
	sqlQuery := "UPDATE user SET "
	values := make([]interface{}, 0)

	for key, val := range updateFields {
		sqlQuery += key + "=?,"
		values = append(values, val)
	}
	sqlQuery = strings.TrimSuffix(sqlQuery, ",") // remove the last comma
	sqlQuery += " WHERE Id = ?"
	values = append(values, Id)

	result, err := db.Exec(sqlQuery, values...)
	if err != nil {
		CheckError(err)
		return false
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func GetIdFromPhone(phoneNumber string) string {
	var id string
	err := db.QueryRow("SELECT Id FROM user WHERE Phone = ?", phoneNumber).Scan(&id)
	if err != nil {
		CheckError(err)
	}
	return id
}

func UpdateUserFromPhone(phoneNumber string, updateFields map[string]interface{}) bool {
	sqlQuery := "UPDATE user SET "
	values := make([]interface{}, 0)

	for key, val := range updateFields {
		sqlQuery += key + "=?,"
		values = append(values, val)
	}
	sqlQuery = strings.TrimSuffix(sqlQuery, ",") // remove the last comma
	sqlQuery += " WHERE Phone = ?"
	values = append(values, phoneNumber)

	result, err := db.Exec(sqlQuery, values...)
	if err != nil {
		CheckError(err)
		return false
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected > 0 {
		return true
	} else {
		return false
	}
}

func UpdateVerifyCodeFromPhone(phoneNumber, verifyCode string, isVerify int) bool {
	var count int
	err := db.QueryRow("UPDATE user SET VerifyCode=?, IsVerify=? WHERE Phone=?", verifyCode, isVerify, phoneNumber).Scan(&count)
	if err != nil {
		CheckError(err)
	}
	if count > 0 {
		return false
	} else {
		return true
	}
}

func CheckPhoneNumberAndPassword(phoneNumber, password string) (*User, bool) {
	var user User
	err := db.QueryRow("SELECT * FROM user WHERE Phone = ? AND Password = ?", phoneNumber, password).Scan(
		&user.Id,
		&user.Username,
		&user.Mail,
		&user.Phone,
		&user.Password,
		&user.Gender,
		&user.BirthDate,
		&user.VerifyCode,
		&user.IsVerify,
		&user.PhotoUrl,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, false // Kullanıcı bulunamadı
		}
		CheckError(err)
		return nil, false // Diğer hata durumları
	}
	return &user, true // Kullanıcı bulundu
}
