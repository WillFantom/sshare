package transfer

import (
	"net/http"
	"net/url"

	"fmt"
)

// TransferShFile is information pertaining to a file uploaded to transfer.sh,
// including the download URL and the deletion token.
type TransferShFile struct {
	downloadURL string
	deleteToken string
}

func NewFile(downloadURL, deleteToken string) *TransferShFile {
	return &TransferShFile{
		downloadURL: downloadURL,
		deleteToken: deleteToken,
	}
}

// DownloadURL is the full URL that can be used to download the contents of a
// file uploaded to the https://transfer.sh file store.
func (tshFile TransferShFile) DownloadURL() string {
	return tshFile.downloadURL
}

// DeleteToken is the token required to delete the upload file prior to its
// normal expiration.
func (tshFile TransferShFile) DeleteToken() string {
	return tshFile.deleteToken
}

// Delete attempts to delete the file from transfer.sh using the file's download
// link and delete token. If an error is returned this may indicate that the
// file has already been deleted or that the request to delete was simply not
// successful.
func (tshFile TransferShFile) Delete() error {
	deleteURL, err := url.JoinPath(tshFile.downloadURL, tshFile.deleteToken)
	if err != nil {
		return fmt.Errorf("failed to generate valid url for transfer.sh file delete")
	}
	req, err := http.NewRequest(http.MethodDelete, deleteURL, nil)
	if err != nil {
		return fmt.Errorf("failed to generate file delete request: %w", err)
	}
	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform the delete request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to delete file from transfer.sh: %w", err)
	}
	return nil
}
