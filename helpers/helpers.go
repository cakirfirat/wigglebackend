package helpers

import (
	"bufio"
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/xeonx/timeago"
	"golang.org/x/text/language"
)

func CheckError(err error) {
	if err != nil {
		log.Println(err.Error())
	}
}

func Md5Hash(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return hex.EncodeToString(hasher.Sum(nil))
}

func SendSms(phoneno, message string) {
	_ = godotenv.Load("../.env")

	url := "https://api.netgsm.com.tr/sms/send/get"
	method := "POST"

	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("usercode", "8503051043")
	_ = writer.WriteField("password", "Oo_110308020")
	_ = writer.WriteField("gsmno", phoneno)
	_ = writer.WriteField("message", message)
	_ = writer.WriteField("msgheader", "8503051043")
	err := writer.Close()
	CheckError(err)

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	req, err := http.NewRequest(method, url, payload)

	CheckError(err)

	req.Header.Set("Content-Type", writer.FormDataContentType())

	res, err := client.Do(req)
	CheckError(err)
	defer res.Body.Close()

	// read response body
	scanner := bufio.NewScanner(res.Body)
	var response []byte
	for scanner.Scan() {
		response = append(response, scanner.Bytes()...)
	}

	// print response body
	fmt.Println(string(response))

}

var SECRET = []byte("super-secret-auth")

func CreateJwt(userId string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["exp"] = time.Now().Add(time.Hour * 730).Unix()
	claims["userId"] = userId
	claims["time"] = time.Now().Unix()
	tokenStr, err := token.SignedString(SECRET)
	CheckError(err)
	return tokenStr, nil
}

func ValidateJwt(next func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			bearerToken := strings.Split(authHeader, " ")
			if len(bearerToken) == 2 && strings.EqualFold(bearerToken[0], "Bearer") {
				token, err := jwt.Parse(bearerToken[1], func(t *jwt.Token) (interface{}, error) {
					_, ok := t.Method.(*jwt.SigningMethodHMAC)
					if !ok {
						w.WriteHeader(http.StatusUnauthorized)
						w.Write([]byte("Not authorized"))
					}
					return SECRET, nil
				})
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
					w.Write([]byte("Not authorized"))
				}
				if token.Valid {
					next(w, r)
				}
			} else {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Not authorized"))
			}
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Not authorized"))
		}
	})
}

func ExtractUserId(authHeader string) (string, error) {

	if authHeader == "" {
		return "", errors.New("Authorization header not present")
	}

	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
		return "", errors.New("Invalid Authorization header format")
	}

	tokenStr := tokenParts[1]
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return SECRET, nil
	})
	if err != nil {
		return "", err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userId := claims["userId"].(string)
		return userId, nil
	} else {
		return "", errors.New("Invalid token")
	}
}

func CreateOtp() string {
	rand.Seed(time.Now().UnixNano())
	return strconv.Itoa(rand.Intn(999999))
}

func Localizate(lang, text string) string {
	// Get the current working directory
	cwd, err := os.Getwd()

	CheckError(err)

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	langPath := filepath.Join(cwd, "../helpers", "lang")

	switch lang {
	case "tr-tr":
		bundle.LoadMessageFile(filepath.Join(langPath, "tr-TR.json"))
	case "en-en":
		bundle.LoadMessageFile(filepath.Join(langPath, "en-EN.json"))
	default:
		bundle.LoadMessageFile(filepath.Join(langPath, "en-EN.json"))
	}

	localizer := i18n.NewLocalizer(bundle, lang)

	return localizer.MustLocalize(&i18n.LocalizeConfig{DefaultMessage: &i18n.Message{ID: text}})
}

func GenerateUUID() (string, error) {
	newUUID, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	return newUUID.String(), nil
}

func GetJSONField(jsonData map[string]interface{}, fieldName string) (string, error) {
	if value, ok := jsonData[fieldName].(string); ok {
		return value, nil
	}
	return "", fmt.Errorf("%s eksik", fieldName)
}

func FormatTime(timeString string) string {
	parts := strings.Split(timeString, " m=")

	if len(parts) < 2 {
		fmt.Println("Invalid time format")
		return ""
	}

	// Custom time layout to match your time format
	layout := "2006-01-02 15:04:05.999999 -0700 MST"

	t, err := time.Parse(layout, parts[0])
	if err != nil {
		fmt.Println("Error parsing time:", err)
		return ""
	}

	return timeago.Turkish.Format(t)
}
