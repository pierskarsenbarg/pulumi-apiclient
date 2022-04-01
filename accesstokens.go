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
	Id    string
	Value string `json:"tokenValue"`
}

func (c *Client) CreateAccessToken(description string) (AccessToken, error) {
	var accessToken AccessToken

	path := "user/tokens"
	endpt := c.baseurl.ResolveReference(&url.URL{Path: path})

	values := map[string]string{"description": description}
	data, err := json.Marshal(values)
	if err != nil {
		return accessToken, err
	}

	req, err := http.NewRequest("POST", endpt.String(), bytes.NewBuffer(data))
	if err != nil {
		return accessToken, err
	}

	req.Header.Add("Accept", "application/vnd.pulumi+8")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+c.token)
	res, err := c.c.Do(req)
	if err != nil {
		return accessToken, err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 200, 201:
		err = json.NewDecoder(res.Body).Decode(&accessToken)
		if err != nil {
			return accessToken, err
		}

		return accessToken, nil
	case 400, 401, 403, 404, 500:
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			panic(err)
		}

		if errRes.StatusCode == 0 {
			errRes.StatusCode = res.StatusCode
		}
		return accessToken, &errRes

	default:
		return accessToken, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}

}

func (c *Client) DeleteAccessToken(tokenId string) error {
	if len(tokenId) == 0 {
		return errors.New("tokenid length must be greater than zero")
	}

	path := fmt.Sprintf("user/tokens/%s", tokenId)
	endpt := c.baseurl.ResolveReference(&url.URL{Path: path})

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
