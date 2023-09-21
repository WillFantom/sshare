package transfer

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// TransferSh holds configuration for a transfer.sh action.
type TransferSh struct {
	fileName     string
	maxDownloads int
	maxDays      int
	password     string
}

const (
	// TransferBaseURL is the base URL of transfer.sh using https.
	DefaultTransferBaseURL string = "https://transfer.sh"
)

var (
	transferBaseURL string = DefaultTransferBaseURL
)

func SetTransferShURL(tshURL string) error {
	if u, err := url.Parse(tshURL); err != nil {
		return fmt.Errorf("invalid transfer.sh url: %w", err)
	} else if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("transfer.sh url not http or https scheme")
	} else {
		transferBaseURL = u.String()
	}
	return nil
}

// NewTransferSh returns an instance of a TransferSh config with appropriate
// default value for the context. Specifically, the file name is set to
// authorized_keys with a max of 10 downloads and a life of 2 days.
func NewTransferSh() *TransferSh {
	return &TransferSh{
		fileName:     "authorized_keys",
		maxDownloads: 10,
		maxDays:      2,
		password:     "",
	}
}

// WithFilename can be used to modify the filename of a Transfer.sh upload.
func (tsh TransferSh) WithFilename(filename string) TransferSh {
	tsh.fileName = filename
	return tsh
}

// WithMaxDownloads can be used to modify the maximum downloads for a
// transfer.sh upload.
func (tsh TransferSh) WithMaxDownloads(maxDownloads int) TransferSh {
	tsh.maxDownloads = maxDownloads
	return tsh
}

// WithMaxDays can be used to modify the maximum days a transfer.sh upload is
// visible for.
func (tsh TransferSh) WithMaxDays(maxDays int) TransferSh {
	tsh.maxDays = maxDays
	return tsh
}

// WithPassword can be used to set the serverside encryption password for a
// transfer.sh file upload.
func (tsh TransferSh) WithPassword(password string) TransferSh {
	tsh.password = password
	return tsh
}

// Upload creates a new file with the contents of data and uploads it to
// transfer.sh. Returned are the URLs to both download/curl the file and to
// delete the file form transfer.sh. If the upload fails, an error is returned.
func (tsh TransferSh) Upload(data string) (*TransferShFile, error) {
	uploadURL, err := url.JoinPath(transferBaseURL, tsh.fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to build upload url for transfer.sh: %w", err)
	}
	dataBuffer := bytes.NewReader([]byte(data))
	req, err := http.NewRequest("PUT", uploadURL, dataBuffer)
	if err != nil {
		return nil, fmt.Errorf("failed to build upload request: %w", err)
	}
	req.Header.Set("Content-Type", "text/plain")
	req.Header.Set("Max-Downloads", strconv.Itoa(tsh.maxDownloads))
	req.Header.Set("Max-Days", strconv.Itoa(tsh.maxDays))
	if tsh.password != "" {
		req.Header.Set("X-Encrypt-Password", tsh.password)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform the upload request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to upload data to transfer.sh: %w", err)
	}
	downloadURLBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response from transfer.sh: %w", err)
	}
	downloadURL := string(downloadURLBytes)
	deleteURL := resp.Header.Get("X-Url-Delete")
	if deleteURL == "" || downloadURL == "" {
		return nil, fmt.Errorf("failed to obtain all urls for the transfer.sh upload")
	}
	return NewFile(downloadURL, path.Base(deleteURL)), nil
}
