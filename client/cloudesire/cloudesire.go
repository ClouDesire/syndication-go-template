package cloudesire

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Subscription struct {
	ID               int    `json:"id"`
	DeploymentStatus string `json:"deploymentStatus"`
	Paid             bool   `json:"paid"`
}

func GetSubscription(id int) Subscription {
	client := &http.Client{}

	req := subscriptionRequest(http.MethodGet, id, nil)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var subscription Subscription
	if err := json.Unmarshal(body, &subscription); err != nil {
		panic(err)
	}
	return subscription
}

func UpdateSubscription(id int, status string) {
	log.Printf("Setting subscription %s to %s", strconv.Itoa(id), status)

	if _, ro := os.LookupEnv("CMW_READ_ONLY"); ro {
		return
	}

	client := &http.Client{}

	buf := []byte(fmt.Sprintf(`{"deploymentStatus": "%s"}`, status))
	req := subscriptionRequest(http.MethodPatch, id, bytes.NewBuffer(buf))
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
}

func subscriptionRequest(method string, id int, body io.Reader) *http.Request {
	url := strings.Join([]string{requiredEnv("CMW_BASE_URL"), "subscription", strconv.Itoa(id)}, "/")
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("CMW-Auth-Token", requiredEnv("CMW_AUTH_TOKEN"))
	return req
}

func requiredEnv(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	log.Panic(key, " env variable not found")
	return ""
}
