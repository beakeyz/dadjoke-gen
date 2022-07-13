package requests

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

/*
retrieves and serializes json formated data comming from our dataserver
*/
func GetJson(url string, target interface{}, add_http_prefix bool) error {
	if (strings.HasPrefix("https://", "") || strings.HasPrefix("http://", "")) && add_http_prefix {
		url = "http://" + url
	}
	var c = &http.Client{Timeout: 10 * time.Second}
	r, err := c.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}
