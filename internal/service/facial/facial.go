package facial

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

func assembleWSAuthURL(requestURL, method, apiKey, apiSecret string) (string, string, string, error) {
	u, err := url.Parse(requestURL)
	if err != nil {
		return "", "", "", err
	}

	host := u.Host
	path := u.Path
	date := time.Now().UTC().Format(http.TimeFormat)
	signatureOrigin := fmt.Sprintf("host: %s\ndate: %s\n%s %s HTTP/1.1", host, date, method, path)
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(signatureOrigin))
	signatureSha := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	authorizationOrigin := fmt.Sprintf(
		`api_key="%s", algorithm="hmac-sha256", headers="host date request-line", signature="%s"`,
		apiKey, signatureSha,
	)
	authorization := base64.StdEncoding.EncodeToString([]byte(authorizationOrigin))
	values := url.Values{}
	values.Set("host", host)
	values.Set("date", date)
	values.Set("authorization", authorization)

	return requestURL + "?" + values.Encode(), date, host, nil
}

func genBody(appid, imgData, imgFormat, serverID string) ([]byte, error) {
	body := map[string]interface{}{
		"header": map[string]interface{}{
			"app_id": appid,
			"status": 3,
		},
		"parameter": map[string]interface{}{
			serverID: map[string]interface{}{
				"service_kind":    "face_detect",
				"detect_points":   "1",
				"detect_property": "1",
				"face_detect_result": map[string]interface{}{
					"encoding": "utf8",
					"compress": "raw",
					"format":   "json",
				},
			},
		},
		"payload": map[string]interface{}{
			"input1": map[string]interface{}{
				"encoding": imgFormat,
				"status":   3,
				"image":    imgData,
			},
		},
	}

	return json.Marshal(body)
}

func Detect(appId, apiKey, apiSecret, serverID, imageData, imageFormat string) (interface{}, error) {
	apiURL := fmt.Sprintf("http://api.xf-yun.com/v1/private/%s", serverID)

	requestURL, date, host, err := assembleWSAuthURL(apiURL, "POST", apiKey, apiSecret)
	if err != nil {
		return nil, err
	}

	body, err := genBody(appId, imageData, imageFormat, serverID)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 10 * time.Second}

	httpReq, err := http.NewRequest("POST", requestURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("host", host)
	httpReq.Header.Set("app_id", appId)
	httpReq.Header.Set("date", date)

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]interface{}
	if err := json.Unmarshal(respData, &data); err != nil {
		return nil, err
	}

	payload, ok := data["payload"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid response payload")
	}

	faceResult, ok := payload["face_detect_result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid face_detect_result")
	}

	text, ok := faceResult["text"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid text")
	}

	decoded, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return nil, err
	}

	var result interface{}
	if err = json.Unmarshal(decoded, &result); err != nil {
		return nil, err
	}

	return result, nil
}
