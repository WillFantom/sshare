package pastebin

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// Pastebin contains the details to connect to and use the Pastebin API.
type Pastebin struct {
	url     string
	token   string
	userKey string
	folder  string
}

const (
	// DefaultPastebinURL is the URL for the public Pastebin instance.
	DefaultPastebinURL string = "https://pastebin.com"
)

// NewPastebin creates a new Pastebin with the given developer token. If no
// developer token is provided, an error is returned as this is required for all
// API operations.
func NewPastebin(developerToken string, opts ...PastebinOpt) (*Pastebin, error) {
	if developerToken == "" {
		return nil, fmt.Errorf("must provide a pastebin developer token")
	}
	pb := Pastebin{
		url:     DefaultPastebinURL,
		token:   developerToken,
		userKey: "",
		folder:  "",
	}
	for _, opt := range opts {
		if err := opt(&pb); err != nil {
			return nil, err
		}
	}
	return &pb, nil
}

// GenerateUserKey uses a Pastebin username and password to generate a userkey
// for the Pastebin API that never expires. If the login fails an error is
// returned.
func (pb Pastebin) GenerateUserKey(username, password string) (string, error) {
	path, err := url.JoinPath(pb.url, "api", "api_login.php")
	if err != nil {
		return "", fmt.Errorf("failed to create url for login form: %w", err)
	}
	formData := url.Values{}
	formData.Add("api_dev_key", pb.token)
	formData.Add("api_user_name", username)
	formData.Add("api_user_password", password)
	resp, err := http.PostForm(path, formData)
	if err != nil {
		return "", fmt.Errorf("login request failed: %w", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("pastebin login failed: status %d", resp.StatusCode)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login request response: %w", err)
	}
	return string(bodyBytes), nil
}

// Post creates a new paste on Pastebin with the given options. If successful
// the pastecode is returned, otherwise an error is returned.
func (pb Pastebin) Post(paste string, expiration ExpiryTime, visibility Visibility) (string, error) {
	path, err := url.JoinPath(pb.url, "api", "api_post.php")
	if err != nil {
		return "", fmt.Errorf("failed to create url for post form: %w", err)
	}
	formData := url.Values{}
	formData.Add("api_dev_key", pb.token)
	formData.Add("api_option", "paste")
	formData.Add("api_paste_format", "sshconfig")
	formData.Add("api_paste_private", visibility.String())
	formData.Add("api_paste_expire_date", expiration.String())
	if pb.userKey != "" {
		formData.Add("api_user_key", pb.userKey)
	}
	if pb.folder != "" {
		formData.Add("api_folder_key", pb.folder)
	}
	formData.Add("api_paste_code", paste)
	resp, err := http.PostForm(path, formData)
	if err != nil {
		return "", fmt.Errorf("new paste request failed: %w", err)
	}
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read new paste request response: %w", err)
	}
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("pastebin new paste failed: status %d (%s)", resp.StatusCode, string(bodyBytes))
	}
	pasteURL, err := url.Parse(string(bodyBytes))
	if err != nil {
		return "", fmt.Errorf("no valid paste url was returned")
	}
	return strings.TrimPrefix(pasteURL.Path, "/"), nil
}

func (pb Pastebin) GeneratePasteURL(pastecode string) (string, error) {
	if pasteURL, err := url.JoinPath(pb.url, pastecode); err != nil {
		return "", fmt.Errorf("failed to generate paste url: %w", err)
	} else {
		return pasteURL, nil
	}
}

func (pb Pastebin) GenerateRawPasteURL(pastecode string) (string, error) {
	if pasteURL, err := url.JoinPath(pb.url, "raw", pastecode); err != nil {
		return "", fmt.Errorf("failed to generate raw paste url: %w", err)
	} else {
		return pasteURL, nil
	}
}
