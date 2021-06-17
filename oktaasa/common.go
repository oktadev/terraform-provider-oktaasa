package oktaasa

import (
	"bytes"
	"encoding/json"
	"log"

	"github.com/go-resty/resty/v2"
)

func checkSoftDelete(response []byte) (bool, error) {
	type sftResp struct {
		Name      string `json:"name"`
		DeletedAt string `json:"deleted_at"`
		RemovedAt string `json:"removed_at"`
	}

	var resp sftResp

	err := json.Unmarshal(response, &resp)

	if err != nil {
		return false, err
	}

	if len(resp.DeletedAt) > 0 || len(resp.RemovedAt) > 0 {
		log.Printf("[DEBUG] Object %s was deleted", resp.Name)
		return true, err
	} else {
		return false, err
	}
}

func SendGet(bearer string, path string) (*resty.Response, error) {
	composedUrl := url + path
	client := resty.New()
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Accept":       "application/json",
			"Content-Type": "Application/json"}).
		SetAuthToken(bearer).
		Get(composedUrl)

	log.Printf("[DEBUG] Get to %s. Status code: %d", composedUrl, resp.StatusCode())
	return resp, err

}

func SendPost(BearerToken, path string, body []byte) (*resty.Response, error) {
	composedUrl := url + path
	client := resty.New()
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Accept":       "application/json",
			"Content-Type": "Application/json"}).
		SetAuthToken(BearerToken).
		SetBody(bytes.NewBuffer(body)).
		Post(composedUrl)

	log.Printf("[DEBUG] POST to %s. Status code: %d", composedUrl, resp.StatusCode())
	return resp, err

}

func SendPut(BearerToken, path string, body []byte) (*resty.Response, error) {
	composedUrl := url + path
	client := resty.New()
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Accept":       "application/json",
			"Content-Type": "Application/json"}).
		SetAuthToken(BearerToken).
		SetBody(bytes.NewBuffer(body)).
		Put(composedUrl)

	log.Printf("[DEBUG] PUT to %s. Status code: %d", composedUrl, resp.StatusCode())
	return resp, err

}

func SendDelete(BearerToken, path string, body []byte) (*resty.Response, error) {
	composedUrl := url + path
	client := resty.New()
	resp, err := client.R().
		SetHeaders(map[string]string{
			"Accept":       "application/json",
			"Content-Type": "Application/json"}).
		SetAuthToken(BearerToken).
		SetBody(bytes.NewBuffer(body)).
		Delete(composedUrl)

	log.Printf("[DEBUG] DELETE to %s. Status code: %d", composedUrl, resp.StatusCode())
	return resp, err

}
