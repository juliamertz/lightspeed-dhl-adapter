package dhl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"lightspeed-dhl/config"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog/log"
)

type Client struct {
	session     *AuthSession
	credentials *config.Dhl
	Cluster     string
}

const DefaultCluster = "https://api-gw.dhlparcel.nl"

func New(credentials *config.Dhl, cluster string) Client {
	return Client{
		Cluster:     cluster,
		credentials: credentials,
		session:     nil,
	}
}

type AuthSession struct {
	AccessToken            string   `json:"accessToken"`
	AccessTokenExpiration  int      `json:"accessTokenExpiration"`
	RefreshToken           string   `json:"refreshToken"`
	RefreshTokenExpiration int      `json:"refreshTokenExpiration"`
	AccountNumbers         []string `json:"accountNumbers"`
}

func (s *AuthSession) RefreshTokenExpired() bool {
	if s == nil {
		return true
	}
	now := time.Now().Local().Unix()
	return int64(s.RefreshTokenExpiration) < now
}

func (s *AuthSession) AccessTokenExpired() bool {
	if s == nil {
		return true
	}
	now := time.Now().Local().Unix()
	return int64(s.AccessTokenExpiration) < now
}

func (c *Client) request(endpoint string, method string, body *[]byte) (*http.Response, error) {
	session := c.GetSession()
	if session == nil {

	}

	endpoint = strings.TrimPrefix(endpoint, "/")
	url := fmt.Sprintf("%s/%s", c.Cluster, endpoint)

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	if session != nil {
		req.Header.Set("Authorization", "Bearer "+c.session.AccessToken)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if body != nil {
		req.Body = io.NopCloser(bytes.NewReader(*body))
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode == 429 {
		log.Warn().Msg("Rate limit reached")
	}
	if res.StatusCode == 404 {
		log.Error().Str("Endpoint", endpoint).Interface("Response", res).Msg("404 While trying to interact with dhl api")
	}
	if res.StatusCode == 400 {
		log.Error().Str("Endpoint", endpoint).Str("Method", method).Str("Body", string(*body)).Interface("Response", res).Msg("400 Bad request from DHL api")
	}
	if res.StatusCode != 200 {
		log.Debug().Int("statuscode", res.StatusCode).Msg("Non 200 Statuscode response for dhl api request")
	}

	return res, nil
}

// TODO: Make sure session gets properly revalidated once the refreshtoken has expired
// Returns a session as long as there is a valid way to obtain it without having to re-authenticate
func (c *Client) GetSession() *AuthSession {
	if c.session != nil {
		expired := c.session.AccessTokenExpired()
		fmt.Printf("access token expired: %v\n", expired)
		if !c.session.AccessTokenExpired() {
			return c.session
		}

		expired = c.session.AccessTokenExpired()
		fmt.Printf("refresh expired: %v\n", expired)
		if !c.session.RefreshTokenExpired() {
			return nil
			// err := c.RefreshSession()
			// if err != nil {
			// 	log.Warn().Msg(fmt.Sprintf("unable to refresh session token, error: %e", err))
			// 	return nil
			// }

			return c.session
		}
	}

	return nil
}

// TODO: Test this function
func (c *Client) RefreshSession() error {
	if c.session == nil {
		return errors.New("Session has not been initialized")
	}

	if c.session.RefreshTokenExpired() {
		return errors.New("Refresh token has expired")
	}

	body, err := json.Marshal(struct {
		refreshToken string
	}{
		refreshToken: c.session.RefreshToken,
	})

	if err != nil {
		return err
	}

	res, err := c.request("authenticate/refresh-token", "POST", &body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return errors.New("Received non 200 statuscode for refresh-token request")
	}

	err = json.NewDecoder(res.Body).Decode(&c.session)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Authenticate(credentials config.Dhl) error {
	body, err := json.Marshal(credentials)
	if err != nil {
		return err
	}

	res, err := c.request("authenticate/api-key", "POST", &body)
	if err != nil {
		return err
	}

	var data AuthSession

	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		return err
	}

	c.session = &data
	return nil
}

func (c *Client) CreateDraft(draft *Draft) (error, *http.Response) {
	body, err := json.Marshal(*draft)
	if err != nil {
		return err, nil
	}

	// TODO: assert that we are authenticated at all times, also check if session is expired, in that case we re-authenticate

	res, err := c.request("drafts", "POST", &body)
	if err != nil {
		return err, nil
	}

	if res.StatusCode != 201 {
		return fmt.Errorf("Expected statuscode response 201, got %v", res.StatusCode), nil
	}

	return nil, res
}

func (c *Client) GetDrafts() ([]Draft, error) {
	res, err := c.request("drafts", "GET", nil)
	if err != nil {
		return nil, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result []Draft
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) GetLabelByReference(reference int, conf *config.Secrets) (*Label, error) {
	uri := fmt.Sprintf("labels?orderReferenceFilter=%v", reference)
	res, err := c.request(uri, "GET", nil)
	if err != nil {
		log.Err(err).Stack().Msg("Error getting label by reference")
		return nil, err
	}

	if res.StatusCode == 404 {
		log.Debug().Int("Order reference", reference).Msg("No label found")
		return nil, nil
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		log.Err(err).Stack().Msg("Error reading response body")
		return nil, err
	}

	var result []Label
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Err(err).Str("body", string(body)).Stack().Msg("Error unmarshalling response body")
		return nil, err
	}

	if len(result) == 0 {
		log.Debug().Int("Order reference", reference).Msg("No label found")
		return nil, nil
	}

	return &result[0], nil
}
