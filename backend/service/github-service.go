package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"area/repository"
	"area/schemas"
)

// Constructor

type GithubService interface {
	// Service interface functions
	GetServiceActionInfo() []schemas.Action
	GetServiceReactionInfo() []schemas.Reaction
	FindActionbyName(name string) func(c chan string, option json.RawMessage, idArea uint64)
	FindReactionbyName(name string) func(option json.RawMessage, idArea uint64) string
	// Service specific functions
	AuthGetServiceAccessToken(code string) (token schemas.Token, err error)
	GetUserInfo(accessToken string) (user schemas.User, err error)
	// Actions functions
	// Reactions functions
}

type githubService struct {
	repository        repository.GithubRepository
	serviceRepository repository.ServiceRepository
	areaRepository    repository.AreaRepository
	tokenRepository   repository.TokenRepository
	serviceInfo       schemas.Service
}

func NewGithubService(
	repository repository.GithubRepository,
	serviceRepository repository.ServiceRepository,
	areaRepository repository.AreaRepository,
	tokenRepository repository.TokenRepository,
) GithubService {
	return &githubService{
		repository:        repository,
		serviceRepository: serviceRepository,
		areaRepository:    areaRepository,
		tokenRepository:   tokenRepository,
		serviceInfo: schemas.Service{
			Name:        schemas.Github,
			Description: "This service is a code repository service",
			Oauth:       true,
			Color:       "#000000",
			Icon:        "https://api.iconify.design/mdi:github.svg?color=%23FFFFFF",
		},
	}
}

// Service interface functions

func (service *githubService) GetServiceInfo() schemas.Service {
	return service.serviceInfo
}

func (service *githubService) GetServiceActionInfo() []schemas.Action {
	return []schemas.Action{}
}

func (service *githubService) GetServiceReactionInfo() []schemas.Reaction {
	return []schemas.Reaction{}
}

func (service *githubService) FindActionbyName(
	name string,
) func(c chan string, option json.RawMessage, idArea uint64) {
	switch name {
	default:
		return nil
	}
}

func (service *githubService) FindReactionbyName(
	name string,
) func(option json.RawMessage, idArea uint64) string {
	switch name {
	default:
		return nil
	}
}

// Service specific functions

func (service *githubService) AuthGetServiceAccessToken(
	code string,
) (token schemas.Token, err error) {
	clientID := os.Getenv("GITHUB_CLIENT_ID")
	if clientID == "" {
		return schemas.Token{}, schemas.ErrGithubClientIdNotSet
	}

	clientSecret := os.Getenv("GITHUB_SECRET")
	if clientSecret == "" {
		return schemas.Token{}, schemas.ErrGithubSecretNotSet
	}

	appPort := os.Getenv("BACKEND_PORT")
	if appPort == "" {
		return schemas.Token{}, schemas.ErrBackendPortNotSet
	}

	redirectURI := "http://localhost:8081/services/github"

	apiURL := "https://github.com/login/oauth/access_token"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("client_secret", clientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, nil)
	if err != nil {
		return schemas.Token{}, fmt.Errorf("unable to create request because %w", err)
	}

	req.URL.RawQuery = data.Encode()
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return schemas.Token{}, fmt.Errorf("unable to make request because %w", err)
	}

	var result schemas.GitHubTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return schemas.Token{}, fmt.Errorf(
			"unable to decode response because %w",
			err,
		)
	}

	resp.Body.Close()

	token = schemas.Token{
		Token: result.AccessToken,
		// RefreshToken:  result.RefreshToken,
		// ExpireAt: result.ExpiresIn,
	}
	return token, nil
}

func (service *githubService) GetUserEmail(accessToken string) (email string, err error) {
	ctx := context.Background()

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.github.com/user/emails",
		nil,
	)
	if err != nil {
		return email, fmt.Errorf("unable to create request because %w", err)
	}

	// Add the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("X-GitHub-Api-Version", "2022-11-28")

	// Make the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return email, fmt.Errorf("unable to make request because %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		// Read and log the error response for debugging
		errorBody, _ := io.ReadAll(resp.Body)
		return email, fmt.Errorf(
			"unexpected status code: %d, response: %s",
			resp.StatusCode,
			string(errorBody),
		)
	}

	result := []schemas.GithubUserEmail{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return email, fmt.Errorf("unable to decode response because %w", err)
	}

	resp.Body.Close()

	for _, email := range result {
		if email.Primary {
			return email.Email, nil
		}
	}

	return email, fmt.Errorf("unable to find primary email")
}

func (service *githubService) GetUserInfoAccount(
	accessToken string,
) (user schemas.User, err error) {
	ctx := context.Background()

	// Create a new HTTP request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.github.com/user", nil)
	if err != nil {
		return user, fmt.Errorf("unable to create request because %w", err)
	}

	// Add the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Make the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return schemas.User{}, fmt.Errorf("unable to make request because %w", err)
	}

	result := schemas.GithubUserInfo{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return schemas.User{}, fmt.Errorf("unable to decode response because %w", err)
	}

	resp.Body.Close()

	user = schemas.User{
		Username: result.Login,
		Email:    result.Email,
	}
	return user, nil
}

func (service *githubService) GetUserInfo(accessToken string) (user schemas.User, err error) {
	user, err = service.GetUserInfoAccount(accessToken)
	if err != nil {
		return user, err
	}

	email, err := service.GetUserEmail(accessToken)
	if err != nil {
		return user, err
	}

	user = schemas.User{
		Username: user.Username,
		Email:    email,
	}

	fmt.Printf("user %+v\n", user)
	return user, nil
}

// Actions functions
// Reactions functions
