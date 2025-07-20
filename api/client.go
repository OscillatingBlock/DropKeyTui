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

const backendURL = "http://localhost:8081"

type ErrMsg error

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

func RegisterUser(pubKeyB64 string) tea.Cmd {
	return func() tea.Msg {
		reqBody := RegisterUserRequest{PublicKey: pubKeyB64}
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to marshal request: %w", err))
		}

		resp, err := httpClient.Post(fmt.Sprintf("%s/api/users", backendURL), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to make register User request"))
		}

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ErrMsg(fmt.Errorf("register user request failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
		}

		var registerResp RegisterUserResponse
		if err := json.NewDecoder(resp.Body).Decode(&registerResp); err != nil {
			return ErrMsg(fmt.Errorf("failed to decode register response: %w", err))
		}
		return registerResp
	}
}

func AuthenticateUser(reqBody AuthRequest) tea.Cmd {
	return func() tea.Msg {
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to marshal request: %w", err))
		}

		resp, err := httpClient.Post(fmt.Sprintf("%s/api/users/auth", backendURL), "application/json", bytes.NewBuffer(jsonBody))
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to authenticate user userID : %v", reqBody.ID))
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ErrMsg(fmt.Errorf("auth request failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
		}

		var authResponse AuthResponse
		if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
			return ErrMsg(fmt.Errorf("failed to decode auth response: %w", err))
		}

		return authResponse
	}
}
