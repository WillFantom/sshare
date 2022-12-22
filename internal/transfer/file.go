package transfer

// TransferShFile is information pertaining to a file uploaded to transfer.sh,
// including the download URL and the deletion token.
type TransferShFile struct {
	downloadURL string
	deleteToken string
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
