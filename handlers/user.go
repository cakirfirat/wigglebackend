package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	. "wigglebackend/helpers"
	. "wigglebackend/models"
)

func GenesisHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var jsonData map[string]interface{}
	errorDecoder := json.NewDecoder(r.Body).Decode(&jsonData)
	CheckError(errorDecoder)

	phoneNumber, err := GetJSONField(jsonData, "phoneNumber")
	locale := r.Header.Get("Accept-Language")
	user.CreatedAt = time.Now().String()
	user.UpdatedAt = time.Now().String()
	user.Phone = phoneNumber

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Genesis error 0001")))
		return
	}
	if CheckPhoneNumber(phoneNumber) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Genesis error 0002")))
		return
	}
	otp := CreateOtp()
	user.VerifyCode = otp
	if success, userID := InsertUser(user); success {
		SendSms(phoneNumber, Localizate(locale, "OTP Message")+otp)

		response := map[string]interface{}{
			"userId": userID,
		}
		jsonResponseData, err := json.Marshal(response)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(Localizate(locale, "Genesis error 0003")))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(jsonResponseData)
		return
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Genesis error 0004")))
		return
	}

}

func VerifyCodeHandler(w http.ResponseWriter, r *http.Request) {
	var jsonData map[string]interface{}
	errorDecoder := json.NewDecoder(r.Body).Decode(&jsonData)
	CheckError(errorDecoder)

	phoneNumber, err := GetJSONField(jsonData, "phoneNumber")
	verifyCode, err := GetJSONField(jsonData, "verifyCode")
	userId, errorUserId := GetJSONField(jsonData, "userId")
	CheckError(err)
	if errorUserId != nil {
		userId = GetIdFromPhone(phoneNumber)
	}
	timeNow := time.Now().String()
	updateFields := map[string]interface{}{
		"IsVerify":  1,
		"UpdatedAt": timeNow,
	}
	locale := r.Header.Get("Accept-Language")
	if !CheckVerifyCode(phoneNumber, verifyCode) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "VerifyCode error 0001")))
		return
	}
	if !UpdateUserFromId(userId, updateFields) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "VerifyCode error 0002")))
		return
	}
	token, err := CreateJwt(userId)

	response := map[string]interface{}{
		"accessToken": token,
	}
	responseJson, err := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(responseJson))
	return

}

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var jsonData map[string]interface{}
	errorDecoder := json.NewDecoder(r.Body).Decode(&jsonData)
	locale := r.Header.Get("Accept-Language")
	token := r.Header.Get("Authorization")

	CheckError(errorDecoder)

	userId, err := ExtractUserId(token)

	fmt.Println(userId)

	gender, err := GetJSONField(jsonData, "gender")
	birthDate, err := GetJSONField(jsonData, "birthDate")
	password, err := GetJSONField(jsonData, "password")
	hashPassword := Md5Hash(password)
	username, err := GetJSONField(jsonData, "username")

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Register error 0001")))
		return
	}

	updateFields := map[string]interface{}{
		"Gender":    gender,
		"BirthDate": birthDate,
		"Password":  hashPassword,
		"Username":  username,
	}

	if !UpdateUserFromId(userId, updateFields) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Register error 0001")))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Localizate(locale, "Register success 0002")))
	return

}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var jsonData map[string]interface{}
	errorDecoder := json.NewDecoder(r.Body).Decode(&jsonData)
	locale := r.Header.Get("Accept-Language")

	CheckError(errorDecoder)

	phoneNumber, err := GetJSONField(jsonData, "phoneNumber")
	password, err := GetJSONField(jsonData, "password")
	hashPassword := Md5Hash(password)

	user, isAuthenticated := CheckPhoneNumberAndPassword(phoneNumber, hashPassword)
	if !isAuthenticated {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		response := struct {
			Message string `json:"message"`
		}{
			Message: Localizate(locale, "Login error 0001"),
		}
		jsonResponse, err := json.Marshal(response)
		if err != nil {
			CheckError(err)
			// Diğer hata durumlarıyla ilgili işlemler
		}
		w.Write(jsonResponse)
		return
	}

	CheckError(err)

	token, err := CreateJwt(user.Id)

	response := map[string]interface{}{
		"userId":      user.Id,
		"accessToken": token,
		"username":    user.Username,
		"photoUrl":    user.PhotoUrl,
		"birthDate":   user.BirthDate,
		"gender":      user.Gender,
		"phone":       user.Phone,
	}

	responseJson, err := json.Marshal(response)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(responseJson)
	return

}

func ForgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	var jsonData map[string]interface{}
	errorDecoder := json.NewDecoder(r.Body).Decode(&jsonData)
	CheckError(errorDecoder)

	phoneNumber, err := GetJSONField(jsonData, "phoneNumber")
	locale := r.Header.Get("Accept-Language")
	user.Phone = phoneNumber

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Verify Code error 0001")))
		return
	}
	if !CheckPhoneNumber(phoneNumber) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Verify Code error 0002")))
		return
	}
	otp := CreateOtp()
	user.VerifyCode = otp

	SendSms(phoneNumber, Localizate(locale, "OTP Message")+otp)

	if !UpdateVerifyCodeFromPhone(phoneNumber, otp, 2) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(Localizate(locale, "Verify Code error 0003")))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(Localizate(locale, "Verify Code success 0001")))
	return

}
