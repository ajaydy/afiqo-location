package maps

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func Get(ctx context.Context, ExtraURL string, parameters map[string]string, data interface{}) error {

	url := fmt.Sprintf("%s%s", url, ExtraURL)

	client := &http.Client{}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()

	q.Add("key", apiKey)

	for key, parameter := range parameters {
		q.Add(key, parameter)
	}

	req.URL.RawQuery = q.Encode()

	fmt.Println(req.URL)

	resp, err := client.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	rBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(rBody, data); err != nil {
		log.Printf("error on unmarshal %v", err)
		var errorResponse ErrorModel
		err := json.Unmarshal(rBody, &errorResponse)
		if err != nil {
			return err
		}

		return errors.New(fmt.Sprintf(`Invalid params %s`, errorResponse.Error.Message))
	}

	return nil

}
