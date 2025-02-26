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

type SpotifyService interface {
	// Service interface functions
	FindActionByName(name string) func(c chan string, option json.RawMessage, area schemas.Area)
	FindReactionByName(name string) func(option json.RawMessage, area schemas.Area) string
	GetServiceActionInfo() []schemas.Action
	GetServiceReactionInfo() []schemas.Reaction
	// Service specific functions
	AuthGetServiceAccessToken(code string) (token schemas.Token, err error)
	GetUserInfo(accessToken string) (user schemas.User, err error)
	// Actions functions
	SpotifyActionMusicPlayed(c chan string, option json.RawMessage, area schemas.Area)
	// Reactions functions
	SpotifyReactionSkipNextMusic(option json.RawMessage, area schemas.Area) string
	SpotifyReactionSkipPreviousMusic(option json.RawMessage, area schemas.Area) string
}

type spotifyService struct {
	repository        repository.SpotifyRepository // This is a repository for the Spotify service
	serviceRepository repository.ServiceRepository // This is a repository for the service
	areaRepository    repository.AreaRepository    // This is a repository for the area
	tokenRepository   repository.TokenRepository   // This is a repository for the token
	serviceInfo       schemas.Service              // This is the service information
}

// NewSpotifyService creates a new instance of SpotifyService with the provided repositories.
// It initializes the service with Spotify-specific information such as name, description, OAuth support, color, and icon.
//
// Parameters:
//   - githubTokenRepository: repository.SpotifyRepository - Repository for handling Spotify tokens.
//   - serviceRepository: repository.ServiceRepository - Repository for handling service-related operations.
//   - areaRepository: repository.AreaRepository - Repository for handling area-related operations.
//   - tokenRepository: repository.TokenRepository - Repository for handling general token operations.
//
// Returns:
//   - SpotifyService: A new instance of SpotifyService.
func NewSpotifyService(
	githubTokenRepository repository.SpotifyRepository,
	serviceRepository repository.ServiceRepository,
	areaRepository repository.AreaRepository,
	tokenRepository repository.TokenRepository,
) SpotifyService {
	return &spotifyService{
		repository:        githubTokenRepository,
		serviceRepository: serviceRepository,
		areaRepository:    areaRepository,
		tokenRepository:   tokenRepository,
		serviceInfo: schemas.Service{
			Name:        schemas.Spotify,
			Description: "This service is a music service",
			Oauth:       true,
			Color:       "#1DC000",
			Icon:        "https://api.iconify.design/mdi:spotify.svg?color=%23FFFFFF",
		},
	}
}

// Service interface functions

// GetServiceInfo returns the service information for the Spotify service.
// It retrieves the service information from the service's internal state.
//
// Returns:
//
//	schemas.Service: The service information for the Spotify service.
func (service *spotifyService) GetServiceInfo() schemas.Service {
	return service.serviceInfo
}

// FindActionByName returns a function that matches the given action name.
// The returned function takes a channel, a JSON raw message, and an area schema as parameters.
// If the action name matches a predefined action, the corresponding function is returned.
// If the action name does not match any predefined actions, nil is returned.
//
// Parameters:
// - name: The name of the action to find.
//
// Returns:
// - A function that matches the given action name, or nil if no match is found.
func (service *spotifyService) FindActionByName(
	name string,
) func(c chan string, option json.RawMessage, area schemas.Area) {
	switch name {
	case string(schemas.MusicPlayed):
		return service.SpotifyActionMusicPlayed
	default:
		return nil
	}
}

// FindReactionByName returns a function that performs a specific Spotify reaction
// based on the provided name. The returned function takes a JSON raw message and
// an area schema as parameters and returns a string.
//
// Parameters:
//   - name: The name of the reaction to find.
//
// Returns:
//   - A function that takes a JSON raw message and an area schema as parameters
//     and returns a string. If the name does not match any known reactions, it
//     returns nil.
func (service *spotifyService) FindReactionByName(
	name string,
) func(option json.RawMessage, area schemas.Area) string {
	switch name {
	case string(schemas.SkipNextMusic):
		return service.SpotifyReactionSkipNextMusic
	case string(schemas.SkipPreviousMusic):
		return service.SpotifyReactionSkipPreviousMusic
	default:
		return nil
	}
}

// GetServiceActionInfo retrieves information about the Spotify service action.
// It marshals a default SpotifyActionMusicPlayedOption to JSON and updates the serviceInfo
// by finding the service by name using the service repository. If any errors occur during
// these operations, they are printed to the console. The function returns a slice of
// schemas.Action containing details about the "MusicPlayed" action.
//
// Returns:
//
//	[]schemas.Action: A slice containing the action information for the Spotify service.
func (service *spotifyService) GetServiceActionInfo() []schemas.Action {
	defaultValue := schemas.SpotifyActionMusicPlayedOption{
		Name: "Believer",
	}
	option, err := json.Marshal(defaultValue)
	if err != nil {
		println("error marshal timer option: " + err.Error())
	}
	service.serviceInfo, err = service.serviceRepository.FindByName(
		schemas.Spotify,
	) // must update the serviceInfo
	if err != nil {
		println("error find service by name: " + err.Error())
	}
	return []schemas.Action{
		{
			Name:               string(schemas.MusicPlayed),
			Description:        "This action check if a music is played",
			Service:            service.serviceInfo,
			Option:             option,
			MinimumRefreshRate: 10,
		},
	}
}

// GetServiceReactionInfo retrieves the reaction information for the Spotify service.
// It marshals a default option value to JSON and updates the service information
// by finding the service by name. If any errors occur during these operations,
// they are printed to the console. The function returns a slice of Reaction
// structs containing the name, description, service information, and option for
// each reaction.
func (service *spotifyService) GetServiceReactionInfo() []schemas.Reaction {
	defaultValue := struct{}{}
	option, err := json.Marshal(defaultValue)
	if err != nil {
		println("error marshal timer option: " + err.Error())
	}
	service.serviceInfo, err = service.serviceRepository.FindByName(
		schemas.Spotify,
	) // must update the serviceInfo
	if err != nil {
		println("error find service by name: " + err.Error())
	}
	return []schemas.Reaction{
		{
			Name:        string(schemas.SkipNextMusic),
			Description: "This reaction will skip to the next music",
			Service:     service.serviceInfo,
			Option:      option,
		},
		{
			Name:        string(schemas.SkipPreviousMusic),
			Description: "This reaction will skip to the previous music",
			Service:     service.serviceInfo,
			Option:      option,
		},
	}
}

// Service specific functions

// AuthGetServiceAccessToken exchanges an authorization code for a Spotify access token.
// It retrieves the client ID and secret from environment variables, constructs the
// necessary request to the Spotify API, and returns the access token along with any
// error encountered during the process.
//
// Parameters:
//   - code: The authorization code received from Spotify's authorization endpoint.
//
// Returns:
//   - token: The access token and related information.
//   - err: An error if the token exchange fails or any other issue occurs.
func (service *spotifyService) AuthGetServiceAccessToken(
	code string,
) (token schemas.Token, err error) {
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	if clientID == "" {
		return schemas.Token{}, schemas.ErrSpotifyClientIdNotSet
	}

	clientSecret := os.Getenv("SPOTIFY_SECRET")
	if clientSecret == "" {
		return schemas.Token{}, schemas.ErrSpotifySecretNotSet
	}

	redirectURI, err := getRedirectURI(service.serviceInfo.Name)
	if err != nil {
		return schemas.Token{}, fmt.Errorf("unable to get redirect URI because %w", err)
	}

	apiURL := "https://accounts.spotify.com/api/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, apiURL, nil)
	if err != nil {
		return schemas.Token{}, fmt.Errorf(
			"unable to create request because %w",
			err,
		)
	}

	req.URL.RawQuery = data.Encode()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(clientID, clientSecret)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return schemas.Token{}, fmt.Errorf("unable to make request because %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		println("Status code", resp.StatusCode)
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("body: %+v\n", body)
		return schemas.Token{}, fmt.Errorf(
			"unable to get token because %v",
			resp.Status,
		)
	}

	var result schemas.SpotifyTokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return schemas.Token{}, fmt.Errorf(
			"unable to decode response because %w",
			err,
		)
	}

	if result.AccessToken == "" {
		fmt.Printf("Token exchange failed. Response body: %v\n", resp.Body)
		return schemas.Token{}, schemas.ErrAccessTokenNotFoundInResponse
	}

	resp.Body.Close()

	token = schemas.Token{
		Token:        result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpireAt:     time.Now().Add(time.Duration(result.ExpiresIn) * time.Second),
	}

	return token, nil
}

// GetUserInfo retrieves the Spotify user information using the provided access token.
// It sends a GET request to the Spotify API endpoint "https://api.spotify.com/v1/me".
// The access token is included in the Authorization header of the request.
//
// Parameters:
//   - accessToken: A string containing the Spotify access token.
//
// Returns:
//   - user: A schemas.User struct containing the user's information (username and email).
//   - err: An error if the request fails or the response cannot be decoded.
//
// Possible errors:
//   - If the HTTP request cannot be created or executed.
//   - If the response status code is not 200 OK.
//   - If the response body cannot be decoded into the expected struct.
func (service *spotifyService) GetUserInfo(accessToken string) (user schemas.User, err error) {
	ctx := context.Background()
	// Create a new HTTP request
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.spotify.com/v1/me",
		nil,
	)
	if err != nil {
		return schemas.User{}, fmt.Errorf("unable to create request because %w", err)
	}

	// Add the Authorization header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	println("accessToken", accessToken)

	// Make the request using the default HTTP client
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return schemas.User{}, fmt.Errorf("unable to make request because %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		errorResponse := schemas.SpotifyErrorResponse{}
		err = json.NewDecoder(resp.Body).Decode(&errorResponse)
		if err != nil {
			return schemas.User{}, fmt.Errorf(
				"unable to decode error response because %w",
				err,
			)
		}

		resp.Body.Close()
		return schemas.User{}, fmt.Errorf(
			"unable to get user info because %v %v",
			errorResponse.Error.Status,
			errorResponse.Error.Message,
		)
	}

	result := schemas.SpotifyUserInfo{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return schemas.User{}, fmt.Errorf("unable to decode response because %w", err)
	}

	resp.Body.Close()

	user = schemas.User{
		Username: result.DisplayName,
		Email:    result.Email,
	}

	return user, nil
}

// getSpotifyPlaybackResponse retrieves the current playback state from the Spotify API.
// It takes a schemas.Token as an argument, which contains the access token required for authentication.
// The function returns a schemas.SpotifyPlaybackResponse containing the playback state and an error if any occurred during the request.
//
// The function performs the following steps:
// 1. Creates a new HTTP GET request to the Spotify API endpoint for the current playback state.
// 2. Sets the Authorization header with the provided access token.
// 3. Sends the request using an HTTP client.
// 4. Checks the response status code to ensure it is 200 OK.
// 5. Decodes the JSON response body into a schemas.SpotifyPlaybackResponse struct.
// 6. Returns the decoded playback response and any error encountered during the process.
//
// Parameters:
// - token: schemas.Token containing the access token for Spotify API authentication.
//
// Returns:
// - schemas.SpotifyPlaybackResponse: The current playback state from the Spotify API.
// - error: An error if any occurred during the request or response processing.
func getSpotifyPlaybackResponse(token schemas.Token) (schemas.SpotifyPlaybackResponse, error) {
	apiURL := "https://api.spotify.com/v1/me/player"

	ctx := context.Background()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, apiURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return schemas.SpotifyPlaybackResponse{}, err
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return schemas.SpotifyPlaybackResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error: Status code %d\n", resp.StatusCode)
		return schemas.SpotifyPlaybackResponse{}, err
	}

	var playbackResponse schemas.SpotifyPlaybackResponse
	err = json.NewDecoder(resp.Body).Decode(&playbackResponse)
	if err != nil {
		fmt.Println("Error decoding response:", err)
		return schemas.SpotifyPlaybackResponse{}, err
	}

	return playbackResponse, nil
}

// InitializedSpotifyStorageVariable initializes the Spotify storage variable for a given area.
// It attempts to unmarshal the storage variable from the area. If unmarshaling fails, it initializes
// the storage variable to a default false value and updates the area in the repository.
//
// Parameters:
//   - area: The area containing the storage variable to be initialized.
//
// Returns:
//   - variable: The initialized Spotify storage variable.
//   - err: An error if any occurred during the process.
func (service *spotifyService) InitializedSpotifyStorageVariable(
	area schemas.Area,
) (variable schemas.SpotifyStorageVariable, err error) {
	variable = schemas.SpotifyStorageVariableInit
	err = json.Unmarshal(area.StorageVariable, &variable)
	if err != nil {
		toto := struct{}{}
		err = json.Unmarshal(area.StorageVariable, &toto)
		if err != nil {
			println("error unmarshaling storage variable: " + err.Error())
			return variable, err
		} else {
			println("initializing storage variable")
			variable = schemas.SpotifyStorageVariableFalse
			area.StorageVariable, err = json.Marshal(variable)
			if err != nil {
				println("error marshalling storage variable: " + err.Error())
				return variable, err
			}
			err = service.areaRepository.Update(area)
			if err != nil {
				println("error updating area: " + err.Error())
				return variable, err
			}
		}
	}

	if variable == schemas.SpotifyStorageVariableInit {
		variable = schemas.SpotifyStorageVariableFalse
		area.StorageVariable, err = json.Marshal(variable)
		if err != nil {
			println("error marshalling storage variable: " + err.Error())
			return variable, err
		}
		err = service.areaRepository.Update(area)
		if err != nil {
			println("error updating area: " + err.Error())
			return variable, err
		}
	}
	return variable, nil
}

// Actions functions

// SpotifyActionMusicPlayed handles the action when music is played on Spotify.
// It checks the current playback status and updates the area storage variable accordingly.
// If the currently playing track matches the expected track, it sends a message to the provided channel.
//
// Parameters:
//   - c: A channel to send messages.
//   - option: A JSON raw message containing the options for the action.
//   - area: The area schema containing user and action details.
//
// The function performs the following steps:
//  1. Unmarshals the option JSON into a SpotifyActionMusicPlayedOption struct.
//  2. Initializes the Spotify storage variable for the given area.
//  3. Retrieves the token for the user and service.
//  4. Gets the current playback response from Spotify.
//  5. Checks if music is currently playing and if the track matches the expected track.
//  6. Updates the area storage variable and sends a message if the track matches.
//  7. Updates the area storage variable if the track does not match or no music is playing.
//  8. Sleeps for the minimum refresh rate or the action refresh rate, whichever is greater.
func (service *spotifyService) SpotifyActionMusicPlayed(
	c chan string,
	option json.RawMessage,
	area schemas.Area,
) {
	optionJSON := schemas.SpotifyActionMusicPlayedOption{}
	err := json.Unmarshal(option, &optionJSON)
	if err != nil {
		fmt.Println("Error unmarshalling option:", err)
		return
	}

	variableDatabaseStorage, err := service.InitializedSpotifyStorageVariable(area)
	if err != nil {
		println("error initializing storage variable: " + err.Error())
	}

	token, err := service.tokenRepository.FindByUserIdAndServiceId(
		area.UserId,
		area.Action.ServiceId,
	)
	if err != nil || token.Token == "" {
		fmt.Println("Error finding token or token not found")
		return
	}

	playbackResponse, err := getSpotifyPlaybackResponse(token)
	if err != nil {
		fmt.Println("Error getting playback response:", err)
		return
	}

	if playbackResponse.IsPlaying {
		artistNames := []string{}
		for _, artist := range playbackResponse.Item.Artists {
			artistNames = append(artistNames, artist.Name)
		}
		if strings.EqualFold(playbackResponse.Item.Name, optionJSON.Name) {
			if variableDatabaseStorage == schemas.SpotifyStorageVariableFalse {
				message := fmt.Sprintf("Currently playing: %s by %s",
					playbackResponse.Item.Name,
					strings.Join(artistNames, ", "),
				)
				variableDatabaseStorage = schemas.SpotifyStorageVariableTrue
				area.StorageVariable, err = json.Marshal(variableDatabaseStorage)
				if err != nil {
					println("error marshalling storage variable: " + err.Error())
					return
				}
				err = service.areaRepository.Update(area)
				if err != nil {
					println("error updating area: " + err.Error())
					return
				}
				fmt.Println(message)
				c <- message
			}
		} else {
			if variableDatabaseStorage == schemas.SpotifyStorageVariableTrue {
				variableDatabaseStorage = schemas.SpotifyStorageVariableFalse
				area.StorageVariable, err = json.Marshal(variableDatabaseStorage)
				if err != nil {
					println("error marshalling storage variable: " + err.Error())
					return
				}
				err = service.areaRepository.Update(area)
				if err != nil {
					println("error updating area: " + err.Error())
					return
				}
			}
			message := fmt.Sprintf("Currently playing: %s by %s, but expected: %s",
				playbackResponse.Item.Name,
				strings.Join(artistNames, ", "),
				optionJSON.Name,
			)
			fmt.Println(message)
		}
	} else {
		if variableDatabaseStorage == schemas.SpotifyStorageVariableTrue {
			variableDatabaseStorage = schemas.SpotifyStorageVariableFalse
			area.StorageVariable, err = json.Marshal(variableDatabaseStorage)
			if err != nil {
				println("error marshalling storage variable: " + err.Error())
				return
			}
			err = service.areaRepository.Update(area)
			if err != nil {
				println("error updating area: " + err.Error())
				return
			}
		}
		fmt.Println("No music is currently playing.")
	}

	if (area.Action.MinimumRefreshRate) > area.ActionRefreshRate {
		time.Sleep(time.Second * time.Duration(area.Action.MinimumRefreshRate))
	} else {
		time.Sleep(time.Second * time.Duration(area.ActionRefreshRate))
	}
}

// Reactions functions

// SpotifyReactionSkipNextMusic skips to the next track in the user's Spotify player.
// It takes a JSON raw message option and an Area schema as parameters.
// The function retrieves the user's Spotify token from the token repository using the user ID and service ID.
// If the token is found, it sends a POST request to the Spotify API to skip to the next track.
// The function returns a string indicating the result of the operation.
//
// Parameters:
//   - option: json.RawMessage - The raw JSON message containing options for the reaction.
//   - area: schemas.Area - The area schema containing user and reaction information.
//
// Returns:
//   - string: A message indicating the result of the operation.
func (service *spotifyService) SpotifyReactionSkipNextMusic(
	option json.RawMessage,
	area schemas.Area,
) string {
	token, err := service.tokenRepository.FindByUserIdAndServiceId(
		area.UserId,
		area.Reaction.ServiceId,
	)
	if err != nil {
		fmt.Println("Error finding token:", err)
		return "Error finding token:" + err.Error()
	}
	if token.Token == "" {
		fmt.Println("Error: Token not found")
		return "Error: Token not found"
	}
	apiURL := "https://api.spotify.com/v1/me/player/next"

	ctx := context.Background()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiURL,
		bytes.NewBuffer([]byte("{}")),
	)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "Error creating request:" + err.Error()
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "Error making request:" + err.Error()
	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	return "Spotify skip next music"
}

// SpotifyReactionSkipPreviousMusic skips to the previous track in the user's Spotify player.
// It takes a JSON raw message option and an Area schema as parameters, and returns a string message.
//
// Parameters:
//   - option: json.RawMessage containing additional options for the reaction.
//   - area: schemas.Area containing user and reaction information.
//
// Returns:
//   - A string message indicating the result of the operation.
//
// The function retrieves the user's Spotify token from the token repository using the user ID and service ID.
// If the token is found and valid, it sends a POST request to the Spotify API to skip to the previous track.
// The function handles errors related to token retrieval, request creation, and request execution, and logs appropriate messages.
func (service *spotifyService) SpotifyReactionSkipPreviousMusic(
	option json.RawMessage,
	area schemas.Area,
) string {
	token, err := service.tokenRepository.FindByUserIdAndServiceId(
		area.UserId,
		area.Reaction.ServiceId,
	)
	if err != nil {
		fmt.Println("Error finding token:", err)
		return "Error finding token:" + err.Error()
	}
	if token.Token == "" {
		fmt.Println("Error: Token not found")
		return "Error: Token not found"
	}
	apiURL := "https://api.spotify.com/v1/me/player/previous"

	ctx := context.Background()

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		apiURL,
		bytes.NewBuffer([]byte("{}")),
	)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return "Error creating request:" + err.Error()
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return "Error making request:" + err.Error()
	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	return "SpotifyR skip to previous music"
}
