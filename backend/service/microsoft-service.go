package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"area/repository"
	"area/schemas"
)

// Constructor

type MicrosoftService interface {
	// Service interface functions
	GetServiceActionInfo() []schemas.Action
	GetServiceReactionInfo() []schemas.Reaction
	FindActionbyName(name string) func(c chan string, option json.RawMessage, idArea uint64)
	FindReactionbyName(name string) func(option json.RawMessage, idArea uint64) string
	GetActionsName() []string
	GetReactionsName() []string
	// Service specific functions
	AuthGetServiceAccessToken(code string) (token schemas.Token, err error)
	GetUserInfo(accessToken string) (user schemas.User, err error)
	// Actions functions
	MicrosoftActionReceiveMail(
		channel chan string,
		option json.RawMessage,
		idArea uint64,
	)
	// Reactions functions
	MicrosoftReactionSendMail(
		option json.RawMessage,
		idArea uint64,
	) string
}

type microsoftService struct {
	repository        repository.MicrosoftRepository
	serviceRepository repository.ServiceRepository
	areaRepository    repository.AreaRepository
	tokenRepository   repository.TokenRepository
	actionName        []string
	reactionName      []string
	serviceInfo       schemas.Service
}

func NewMicrosoftService(
	githubTokenRepository repository.MicrosoftRepository,
	serviceRepository repository.ServiceRepository,
	areaRepository repository.AreaRepository,
	tokenRepository repository.TokenRepository,
) MicrosoftService {
	return &microsoftService{
		repository:        githubTokenRepository,
		serviceRepository: serviceRepository,
		areaRepository:    areaRepository,
		tokenRepository:   tokenRepository,
		serviceInfo: schemas.Service{
			Name:        schemas.Microsoft,
			Description: "This service is used to interact with Microsoft services",
			Oauth:       true,
			Color:       "#001DDA",
			Icon:        "https://api.iconify.design/mdi:microsoft.svg?color=%23FFFFFF",
		},
	}
}

// Service interface functions

func (service *microsoftService) GetServiceInfo() schemas.Service {
	return service.serviceInfo
}

func (service *microsoftService) GetServiceActionInfo() []schemas.Action {
	// service.actionName = append(service.actionName, string(schemas.ReceiveMicrosoftMail))
	defaultValue := struct{}{}
	option, err := json.Marshal(defaultValue)
	if err != nil {
		fmt.Println("Error marshalling default options:", err)
	}
	service.serviceInfo, err = service.serviceRepository.FindByName(
		schemas.Microsoft,
	)
	if err != nil {
		println("error find service by name: " + err.Error())
	}
	return []schemas.Action{
		{
			Name:        string(schemas.ReceiveMicrosoftMail),
			Description: "Receive a mail using Microsoft services",
			Service:     service.serviceInfo,
			Option:      option,
		},
	}
}

func (service *microsoftService) GetServiceReactionInfo() []schemas.Reaction {
	service.reactionName = append(service.reactionName, string(schemas.SendMicrosoftMail))
	defaultValue := schemas.MicrosoftReactionSendMailOptions{
		Subject:   "",
		Body:      "",
		Recipient: "",
	}
	option, err := json.Marshal(defaultValue)
	if err != nil {
		fmt.Println("Error marshalling default options:", err)
	}
	service.serviceInfo, err = service.serviceRepository.FindByName(
		schemas.Microsoft,
	)
	if err != nil {
		println("error find service by name: " + err.Error())
	}
	return []schemas.Reaction{
		{
			Name:        string(schemas.SendMicrosoftMail),
			Description: "Send a mail using Microsoft services",
			Service:     service.serviceInfo,
			Option:      option,
		},
	}
}

func (service *microsoftService) FindActionbyName(
	name string,
) func(c chan string, option json.RawMessage, idArea uint64) {
	switch name {
	case string(schemas.ReceiveMicrosoftMail):
		return service.MicrosoftActionReceiveMail
	default:
		return nil
	}
}

func (service *microsoftService) FindReactionbyName(
	name string,
) func(option json.RawMessage, idArea uint64) string {
	switch name {
	case string(schemas.SendMicrosoftMail):
		return service.MicrosoftReactionSendMail
	default:
		return nil
	}
}

func (service *microsoftService) GetActionsName() []string {
	return service.actionName
}

func (service *microsoftService) GetReactionsName() []string {
	return service.reactionName
}

func (service *microsoftService) AuthGetServiceAccessToken(
	code string,
) (token schemas.Token, err error) {
	clientID := os.Getenv("MICROSOFT_CLIENT_ID")
	if clientID == "" {
		return schemas.Token{}, schemas.ErrMicrosoftClientIdNotSet
	}

	appPort := os.Getenv("BACKEND_PORT")
	if appPort == "" {
		return schemas.Token{}, schemas.ErrBackendPortNotSet
	}

	redirectURI := "http://localhost:8081/services/microsoft"

	apiURL := "https://login.microsoftonline.com/common/oauth2/v2.0/token"

	data := url.Values{}
	data.Set("client_id", clientID)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	req, err := http.NewRequest(http.MethodPost, apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return schemas.Token{}, fmt.Errorf("unable to create request because %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return schemas.Token{}, fmt.Errorf("unable to make request because %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Println("response body: ", string(bodyBytes))

	var result schemas.MicrosoftTokenResponse
	err = json.Unmarshal(bodyBytes, &result)
	if err != nil {
		return schemas.Token{}, fmt.Errorf("unable to decode response because %w", err)
	}

	if result.AccessToken == "" || result.TokenType == "" {
		return schemas.Token{}, schemas.ErrAccessTokenNotFoundInResponse
	}

	token = schemas.Token{
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpireAt:     time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
	}
	return token, nil
}

func (service *microsoftService) GetUserInfo(
	accessToken string,
) (user schemas.User, err error) {
	ctx := context.Background()

	url := "https://graph.microsoft.com/v1.0/me"
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return schemas.User{}, fmt.Errorf("unable to create request because %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return schemas.User{}, fmt.Errorf("unable to make request because %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return schemas.User{}, fmt.Errorf("failed to fetch user info: %s", resp.Status)
	}

	var result schemas.MicrosoftUserInfo
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return schemas.User{}, fmt.Errorf("unable to decode response because %w", err)
	}

	user = schemas.User{
		Email:    result.Mail,
		Username: result.DisplayName,
	}
	if user.Email == "" {
		user.Email = result.UserPrincipalName
	}

	return user, nil
}

func (service *microsoftService) MicrosoftActionReceiveMail(
	channel chan string,
	option json.RawMessage,
	idArea uint64,
) {
	area, err := service.areaRepository.FindById(idArea)
	if err != nil {
		println("error finding area: " + err.Error())
		return
	}

	variable := schemas.MicrosoftVariableReceiveMail{}
	err = json.Unmarshal(area.StorageVariable, &variable)
	if err != nil {
		toto := struct{}{}
		err = json.Unmarshal(area.StorageVariable, &toto)
		if err != nil {
			println("error unmarshalling storage variable: " + err.Error())
			return
		} else {
			println("initializing storage variable")
			variable = schemas.MicrosoftVariableReceiveMail{
				Time: time.Now().Add(-time.Hour),
			}
			area.StorageVariable, err = json.Marshal(variable)
			if err != nil {
				println("error marshalling storage variable: " + err.Error())
				return
			}
			service.areaRepository.Update(area)
		}
	}

	if variable.Time.IsZero() {
		println("initializing storage variable")
		variable = schemas.MicrosoftVariableReceiveMail{
			Time: time.Now().Add(-time.Hour),
		}
		area.StorageVariable, err = json.Marshal(variable)
		if err != nil {
			println("error marshalling storage variable: " + err.Error())
			return
		}
		service.areaRepository.Update(area)
	}

	token, err := service.tokenRepository.FindByUserIdAndServiceId(
		area.UserId,
		area.Action.ServiceId,
	)
	if err != nil || token.Token == "" {
		println("error retrieving token or token not found")
		return
	}

	apiURL := "https://graph.microsoft.com/v1.0/me/messages?$filter=receivedDateTime+gt+" + variable.Time.Format(
		"2006-01-02T15:04:05",
	) + "Z"

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		println("error creating request: " + err.Error())
		return
	}

	// Add the authorization header
	req.Header.Set("Authorization", "Bearer "+token.Token)

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		println("error making request: " + err.Error())
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		println("error status code: " + fmt.Sprint(resp.StatusCode))
		return
	}

	// Decode the response
	var emailResponse struct {
		Value []struct {
			ID      string `json:"id"`
			Subject string `json:"subject"`
			From    struct {
				EmailAddress struct {
					Address string `json:"address"`
				} `json:"emailAddress"`
			} `json:"from"`
			ReceivedDateTime string `json:"receivedDateTime"`
		} `json:"value"`
	}

	err = json.NewDecoder(resp.Body).Decode(&emailResponse)
	if err != nil {
		println("error decoding response: " + err.Error())
		return
	}

	// Process the most recent email
	if len(emailResponse.Value) > 0 {
		println("New email received")
		latestEmail := emailResponse.Value[0]
		response := fmt.Sprintf("New email received from %s: object: %s",
			latestEmail.From.EmailAddress.Address,
			latestEmail.Subject,
		)
		println(response)
		variable.Time, err = time.Parse(time.RFC3339, latestEmail.ReceivedDateTime)
		if err != nil {
			println("error parsing time: " + err.Error())
			return
		}
		variable.Time = variable.Time.Add(time.Second)
		area.StorageVariable, err = json.Marshal(variable)
		if err != nil {
			println("error marshalling storage variable: " + err.Error())
			return
		}
		service.areaRepository.Update(area)
		println("response sent to channel")
		channel <- response
	} else {
		println("No new emails")
	}
	time.Sleep(time.Second * 10)
}

func (service *microsoftService) MicrosoftReactionSendMail(
	option json.RawMessage,
	idArea uint64,
) string {
	// Parse the options
	options := schemas.MicrosoftReactionSendMailOptions{}
	err := json.Unmarshal(option, &options)
	if err != nil {
		fmt.Println("Error unmarshalling options:", err)
		return "Error unmarshalling options: " + err.Error()
	}

	// Retrieve the area
	area, err := service.areaRepository.FindById(idArea)
	if err != nil {
		fmt.Println("Error finding area:", err)
		return "Error finding area: " + err.Error()
	}

	// Retrieve the token
	token, err := service.tokenRepository.FindByUserIdAndServiceId(
		area.UserId,
		area.Reaction.ServiceId,
	)
	if err != nil {
		fmt.Println("Error finding token:", err)
		return "Error finding token: " + err.Error()
	}
	if token.Token == "" {
		fmt.Println("Error: Token not found")
		return "Error: Token not found"
	}

	// Microsoft Graph API URL for sending mail
	apiURL := "https://graph.microsoft.com/v1.0/me/sendMail"

	// Construct the email payload
	payload := map[string]interface{}{
		"message": map[string]interface{}{
			"subject": options.Subject,
			"body": map[string]string{
				"contentType": "Text",
				"content":     options.Body,
			},
			"toRecipients": []map[string]map[string]string{
				{
					"emailAddress": {
						"address": options.Recipient,
					},
				},
			},
		},
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling email payload:", err)
		return "Error marshalling email payload: " + err.Error()
	}

	// Create the HTTP request
	req, err := http.NewRequest(http.MethodPost, apiURL, bytes.NewBuffer(payloadBytes))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return "Error creating HTTP request: " + err.Error()
	}

	// Add headers
	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending email request:", err)
		return "Error sending email request: " + err.Error()
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		fmt.Println("Error sending email:", string(bodyBytes))
		return "Error sending email: " + string(bodyBytes)
	}

	return "Email sent successfully!"
}
