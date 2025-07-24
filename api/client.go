package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const backendURL = "http://localhost:8081"

type ErrMsg error

var httpClient = &http.Client{
	Timeout: 5 * time.Second,
}

type PasteCreatedMsg struct {
	TempID string
	CreatePasteResponse
}

type PasteListFetchedMsg struct {
	List []Paste
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

func CreatePaste(reqBody PasteRequest, token, tempID string) tea.Cmd {
	return func() tea.Msg {
		jsonBody, err := json.Marshal(reqBody)
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to marshal request: %w", err))
		}

		req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/pastes", backendURL), bytes.NewBuffer(jsonBody))
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to generate create paste request"))
		}

		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")

		resp, err := httpClient.Do(req)
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to make create paste request"))
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ErrMsg(fmt.Errorf("create paste request failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
		}

		var pasteResponse CreatePasteResponse
		if err := json.NewDecoder(resp.Body).Decode(&pasteResponse); err != nil {
			return ErrMsg(fmt.Errorf("failed to decode create paste response: %w", err))
		}

		return PasteCreatedMsg{
			TempID:              tempID,
			CreatePasteResponse: pasteResponse,
		}
	}
}

func GetPastes(publicKey string) tea.Cmd {
	return func() tea.Msg {
		req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/pastes?public_key=%s", backendURL, url.QueryEscape(publicKey)), nil)
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to create request: %w", err))
		}

		resp, err := httpClient.Do(req)
		if err != nil {
			return ErrMsg(fmt.Errorf("failed to make request: %w", err))
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return ErrMsg(fmt.Errorf("get pastes failed with status %d: %s", resp.StatusCode, string(bodyBytes)))
		}

		var pastes []Paste
		if err := json.NewDecoder(resp.Body).Decode(&pastes); err != nil {
			return ErrMsg(fmt.Errorf("failed to decode response: %w", err))
		}

		return PasteListFetchedMsg{
			List: pastes,
		}
	}
}
