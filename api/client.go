package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const backendURL = "http://localhost:8080"

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func RegisterUser(pubKeyB64 string) tea.Cmd {
	return func() tea.Msg {
		reqBody := RegisterUserRequest{PublicKey: pubKeyB64}
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return error(fmt.Errorf("failed to marshal request: %w", err))
		}

		resp, err := httpClient.Post(fmt.Sprintf("%s/users", backendURL), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			return error(fmt.Errorf("failed to make register User request"))
		}

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return error(fmt.Errorf("register user request failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
		}

		var registerResp RegisterUserResponse
		if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
			return errMsg{err: fmt.Errorf("failed to decode register response: %w", err)}
		}
		user := User{
			ID:        registerResp.ID,
			PublicKey: pubKeyB64,
		}
		return user
	}
}

func AuthenticateUser(reqBody AuthRequest) tea.Cmd {
	return func() tea.Msg {
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return error(fmt.Errorf("failed to marshal request: %w", err))
		}

		resp, err := httpClient.Post(fmt.Sprintf("%s/users/auth", backendURL), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			return error(fmt.Errorf("failed to authenticate user userID : %v", reqBody.ID))
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return error(fmt.Errorf("auth request failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
		}

		var authResponse AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
			return error(fmt.Errorf("failed to decode auth response: %w", err))
		}

		return authResponse
	}
}
