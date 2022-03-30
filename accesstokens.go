package pulumiapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type AccessToken struct {
	id    string
	value string
}

func (c *Client) CreateAccessToken(description string) (*AccessToken, error) {
	path := "user/tokens"
	endpt := baseURL.ResolveReference(&url.URL{Path: path})

	values := map[string]string{"description": description}
	data, err := json.Marshal(values)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", endpt.String(), bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Accept", "application/vnd.pulumi+8")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+c.token)
	res, err := c.c.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	var token AccessToken

	switch res.StatusCode {
	case 200, 201:
		err = json.NewDecoder(res.Body).Decode(&token)
		if err != nil {
			return nil, err
		}

		return &token, nil
	case 400, 401, 403, 404, 500:
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			panic(err)
		}

		if errRes.StatusCode == 0 {
			errRes.StatusCode = res.StatusCode
		}
		return nil, &errRes

	default:
		return nil, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

}

func (c *Client) DeleteToken(tokenId string) error {
	if len(tokenId) == 0 {
		return errors.New("tokenid length must be greater than zero")
	}

	path := fmt.Sprintf("user/tokens/%s", tokenId)
	endpt := baseURL.ResolveReference(&url.URL{Path: path})

	req, err := http.NewRequest("DELETE", endpt.String(), nil)
	if err != nil {
		return err
	}

	req.Header.Add("Accept", "application/vnd.pulumi+8")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+c.token)

	res, err := c.c.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 204:
		return nil
	case 400, 401, 403, 404, 405, 500:
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			panic(err)
		}

		if errRes.StatusCode == 0 {
			errRes.StatusCode = res.StatusCode
		}
		return &errRes
	default:
		return fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
}
