package deathbycaptcha

import (
	"net/http"
	"time"
	"log"
	"bytes"
	"strings"
	"os"
	"mime/multipart"
	"path/filepath"
	"io"
	"errors"
)

type DeathByCaptcha struct {}

func New() (*DeathByCaptcha) {
	return &DeathByCaptcha{}
}

func (DeathByCaptcha) Solve(username, password, imagePath string) (captchaText string, err error) {
	extraParams := map[string]string {
		"username": username,
		"password": password,
	}

	request, err := newfileUploadRequest("http://api.dbcapi.me/api/captcha", extraParams, "captchafile", imagePath)
	if err != nil {
		return
	}
	client := &http.Client{}
	startTime := time.Now()
	log.Printf("Start time: %v\n", startTime)
	resp, err := client.Do(request)
	
	if err != nil {
		return captchaText, err
	} else {
		body := &bytes.Buffer{}
		_, err := body.ReadFrom(resp.Body)
		if err != nil {
			return captchaText, err
		}
		resp.Body.Close()

		resultURL := resp.Request.URL.String()
		result := strings.Split(body.String(), "&")
		if len(result) <= 1 {
			log.Println("Result not solved")
			log.Panicln(body.String())
			
			return captchaText, errors.New("Result not solved")
		}

		for {
			if strings.HasPrefix(result[2], "text=") {
				captchaText = strings.TrimPrefix(result[2], "text=")

				if len(captchaText) > 0 {
					finishTime := time.Now()
					log.Printf("Complete time: %v\n", finishTime)
					log.Printf("Elapsed time: %v\n", finishTime.Sub(startTime))

					break
				}
			}

			
			resp, err = client.Get(resultURL)
			if err != nil {
				return captchaText, err
			}
			body := &bytes.Buffer{}
			_, err := body.ReadFrom(resp.Body)
			if err != nil {
				return captchaText, err
			}
			resp.Body.Close()
			result = strings.Split(body.String(), "&")
		}

		return captchaText, err
	}
}

func newfileUploadRequest(uri string, params map[string]string, paramName, path string) (*http.Request, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", uri, body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req, err
}