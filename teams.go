package pulumiapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

type Teams struct {
	Teams []Team
}

type Team struct {
	Type        string `json:"kind"`
	Name        string
	DisplayName string
	Description string
	Members     []TeamMember
}

type TeamMember struct {
	Name        string
	GithubLogin string
	AvatarUrl   string
	Role        string
}

func (c *Client) ListTeams(orgName string) ([]Team, error) {
	var teams []Team
	if len(orgName) == 0 {
		return teams, errors.New("empty orgName")
	}

	path := fmt.Sprintf("orgs/%s/teams", orgName)
	endpt := baseURL.ResolveReference(&url.URL{Path: path})

	req, err := http.NewRequest("GET", endpt.String(), nil)
	if err != nil {
		return teams, err
	}

	req.Header.Add("Accept", "application/vnd.pulumi+8")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+c.token)
	res, err := c.c.Do(req)
	if err != nil {
		return teams, err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		var teamArray Teams
		err = json.NewDecoder(res.Body).Decode(&teamArray)
		if err != nil {
			return teams, err
		}

		return teamArray.Teams, nil
	case 400, 401, 403, 404, 500:
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			panic(err)
		}

		if errRes.StatusCode == 0 {
			errRes.StatusCode = res.StatusCode
		}
		return teams, &errRes

	default:
		return teams, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
}

func (c *Client) GetTeam(orgName string, teamName string) (Team, error) {
	var team Team
	if len(orgName) == 0 {
		return team, errors.New("empty orgName")
	}

	if len(teamName) == 0 {
		return team, errors.New("empty orgName")
	}

	path := fmt.Sprintf("orgs/%s/teams/%s", orgName, teamName)
	endpt := baseURL.ResolveReference(&url.URL{Path: path})

	req, err := http.NewRequest("GET", endpt.String(), nil)
	if err != nil {
		return team, err
	}

	req.Header.Add("Accept", "application/vnd.pulumi+8")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+c.token)
	res, err := c.c.Do(req)
	if err != nil {
		return team, err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 200:
		err = json.NewDecoder(res.Body).Decode(&team)
		if err != nil {
			return team, err
		}

		return team, nil
	case 400, 401, 403, 404, 500:
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			panic(err)
		}

		if errRes.StatusCode == 0 {
			errRes.StatusCode = res.StatusCode
		}
		return team, &errRes

	default:
		return team, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
}

func (c *Client) CreateTeam(orgName string, teamName string, teamType string, displayName string, description string) (Team, error) {
	var team Team

	if len(orgName) == 0 {
		return team, errors.New("orgname must not be empty")
	}

	if len(teamName) == 0 {
		return team, errors.New("teamname must not be empty")
	}

	if len(teamType) == 0 {
		return team, errors.New("teamtype must not be empty")
	}

	teamtypeList := []string{"github", "pulumi"}
	if !Contains(teamtypeList, teamType) {
		return team, errors.New("teamtype must be either `pulumi` or `github`")
	}

	path := fmt.Sprintf("orgs/%s/teams/%s", orgName, teamType)
	endpt := baseURL.ResolveReference(&url.URL{Path: path})

	values := map[string]string{"organization": orgName, "teamType": teamType, "name": teamName, "displayName": displayName, "description": description}
	data, err := json.Marshal(values)
	if err != nil {
		return team, err
	}

	req, err := http.NewRequest("POST", endpt.String(), bytes.NewBuffer(data))
	if err != nil {
		return team, err
	}

	req.Header.Add("Accept", "application/vnd.pulumi+8")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+c.token)
	res, err := c.c.Do(req)
	if err != nil {
		return team, err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 200, 201:
		err = json.NewDecoder(res.Body).Decode(&team)
		if err != nil {
			return team, err
		}

		return team, nil
	case 400, 401, 403, 404, 500:
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			panic(err)
		}

		if errRes.StatusCode == 0 {
			errRes.StatusCode = res.StatusCode
		}
		return team, &errRes

	default:
		return team, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
}

func (c *Client) UpdateTeam(orgName string, teamName string, displayName string, description string) (Team, error) {
	var team Team

	if len(orgName) == 0 {
		return team, errors.New("orgname must not be empty")
	}

	if len(teamName) == 0 {
		return team, errors.New("teamname must not be empty")
	}

	path := fmt.Sprintf("orgs/%s/teams/%s", orgName, teamName)
	endpt := baseURL.ResolveReference(&url.URL{Path: path})

	values := map[string]string{
		"newDisplayName": displayName,
		"newDescription": description,
	}
	data, err := json.Marshal(values)
	if err != nil {
		return team, err
	}

	req, err := http.NewRequest("PATCH", endpt.String(), bytes.NewBuffer(data))
	if err != nil {
		return team, err
	}

	req.Header.Add("Accept", "application/vnd.pulumi+8")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", "token "+c.token)

	res, err := c.c.Do(req)
	if err != nil {
		return team, err
	}

	defer res.Body.Close()

	switch res.StatusCode {
	case 204:
		team.Description = description
		team.DisplayName = displayName
		return team, nil
	case 400, 401, 403, 404, 405, 500:
		var errRes ErrorResponse
		err = json.NewDecoder(res.Body).Decode(&errRes)
		if err != nil {
			panic(err)
		}

		if errRes.StatusCode == 0 {
			errRes.StatusCode = res.StatusCode
		}
		return team, &errRes
	default:
		return team, fmt.Errorf("unexpected status code %d", res.StatusCode)
	}
}

func (c *Client) DeleteTeam(orgName string, teamName string) error {

	if len(orgName) == 0 {
		return errors.New("orgname must not be empty")
	}

	if len(teamName) == 0 {
		return errors.New("teamname must not be empty")
	}

	path := fmt.Sprintf("orgs/%s/teams/%s", orgName, teamName)
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

func (c *Client) updateTeamMembership(orgName string, teamName string, userName string, addOrRemove string) error {
	if len(orgName) == 0 {
		return errors.New("orgname must not be empty")
	}

	if len(teamName) == 0 {
		return errors.New("teamname must not be empty")
	}

	if len(userName) == 0 {
		return errors.New("username must not be empty")
	}

	addOrRemoveValues := []string{"add", "remove"}
	if !Contains(addOrRemoveValues, addOrRemove) {
		return errors.New("value must be `add` or `remove`")
	}

	path := fmt.Sprintf("orgs/%s/teams/%s", orgName, teamName)
	endpt := baseURL.ResolveReference(&url.URL{Path: path})

	values := map[string]string{"memberAction": addOrRemove, "member": userName}
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("PATCH", endpt.String(), bytes.NewBuffer(data))
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
	case 200, 204:
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

func (c *Client) AddMemberToTeam(orgName string, teamName string, userName string) error {
	err := c.updateTeamMembership(orgName, teamName, userName, "add")
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (c *Client) DeleteMemberFromTeam(orgName string, teamName string, userName string) error {
	err := c.updateTeamMembership(orgName, teamName, userName, "remove")
	if err != nil {
		return err
	} else {
		return nil
	}
}
