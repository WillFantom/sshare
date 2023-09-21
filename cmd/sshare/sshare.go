package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willfantom/sshare/internal/github"
	"github.com/willfantom/sshare/internal/transfer"
	"github.com/willfantom/sshare/internal/ui"
	"github.com/willfantom/sshare/keys"
)

var (
	tshInstanceURL    string   = transfer.DefaultTransferBaseURL
	transferDownloads int      = 0
	transferDays      int      = 0
	password          string   = ""
	sshAgentPath      string   = os.Getenv("SSH_AUTH_SOCK")
	sshAgentPass      string   = ""
	githubToken       string   = ""
	keyFilepaths      []string = make([]string, 0)
	rawKeys           []string = make([]string, 0)
)

var (
	rootCmd = &cobra.Command{
		Use:   "sshare",
		Short: "Easily share links to your SSH public keys",
		Long:  `Share your public SSH keys found in your agent via curl-able transfer.sh links.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if transferDays <= 0 {
				ui.Errorln("Uploaded content must be visible on a transfer.sh instance for at least 1 day", true)
			}
			if transferDownloads <= 0 {
				ui.Errorln("Uploaded content must be downloadable via a transfer.sh instance at least 1 time", true)
			}
			if err := transfer.SetTransferShURL(tshInstanceURL); err != nil {
				ui.Errorln("Transfer.sh URL is not valid: "+err.Error(), true)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			chosenKeys := make([]*keys.Key, 0)

			if len(keyFilepaths) > 0 {
				ui.Infoln("Adding keys from files...")
			}
			for _, fp := range keyFilepaths {
				k, err := keys.NewKeyFromFile(fp)
				if err != nil {
					ui.Errorln(err.Error(), true)
				}
				chosenKeys = append(chosenKeys, k)
			}

			if len(rawKeys) > 0 {
				ui.Infoln("Adding keys from raw values...")
			}
			for _, rk := range rawKeys {
				k, err := keys.NewKey(rk, "")
				if err != nil {
					ui.Errorln(err.Error(), true)
				}
				chosenKeys = append(chosenKeys, k)
			}

			if githubToken != "" {
				ui.Infoln("Adding keys from GitHub...")
				sshAgent := github.NewAgent(githubToken)
				authorizedKeys, err := sshAgent.GetKeys()
				if err != nil {
					ui.Errorln(err.Error(), true)
				}
				selectedKeys, err := ui.SelectKey(authorizedKeys)
				if err != nil {
					ui.Errorln(err.Error(), true)
				}
				chosenKeys = append(chosenKeys, selectedKeys...)
			}

			if sshAgentPath != "" {
				ui.Infoln("Adding keys from SSH agent...")
				agentOpts := make([]keys.AgentOpt, 0)
				if sshAgentPass != "" {
					agentOpts = append(agentOpts, keys.AgentPassphraseOpt(sshAgentPass))
				}
				sshAgent, err := keys.NewSSHAgent(sshAgentPath, agentOpts...)
				if err != nil {
					ui.Errorln(err.Error(), true)
				}
				authorizedKeys, err := sshAgent.GetKeys()
				if err != nil {
					ui.Errorln(err.Error(), true)
				}
				selectedKeys, err := ui.SelectKey(authorizedKeys)
				if err != nil {
					ui.Errorln(err.Error(), true)
				}
				chosenKeys = append(chosenKeys, selectedKeys...)
			} else {
				ui.Warnln("No SSH agent path has been provided. Skipping...")
			}
			if len(chosenKeys) == 0 {
				ui.Errorln("No keys were selected", true)
			}
			keyText := keys.CreateAuthorizedKeys(chosenKeys)
			ui.Infoln(fmt.Sprintf("Generated authorized_keys:\n%s", keyText))
			uploadedFile, err := transfer.NewTransferSh().
				WithMaxDays(transferDays).
				WithMaxDownloads(transferDownloads).
				Upload(keyText)
			if err != nil {
				ui.Errorln(err.Error(), true)
			}
			ui.Successln(fmt.Sprintf("File Download URL: %s", uploadedFile.DownloadURL()))
			ui.Infoln(fmt.Sprintf("File Delete Token: %s", uploadedFile.DeleteToken()))
		},
	}
	deleteCmd = &cobra.Command{
		Use:   "delete [file_url] [delete_token]",
		Short: "Delete an uploaded authorized_keys file",
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			tshFile := transfer.NewFile(args[0], args[1])
			if err := tshFile.Delete(); err != nil {
				ui.Errorln("Failed to delete file", true)
			}
			ui.Successln("File Deleted")
		},
	}
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		ui.Errorln(err.Error(), true)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&sshAgentPath, "agent", "a", sshAgentPath, "path to the target ssh agent socket ($SSH_AUTH_SOCK)")
	rootCmd.Flags().StringVarP(&sshAgentPass, "passphrase", "p", sshAgentPass, "passphrase for the ssh agent")
	rootCmd.Flags().StringVarP(&githubToken, "github-token", "g", githubToken, "github token with permission to read ssh keys")
	rootCmd.Flags().StringSliceVarP(&keyFilepaths, "key-file", "f", keyFilepaths, "additional key file(s) to include in the generated authorized_keys")
	rootCmd.Flags().StringSliceVarP(&rawKeys, "key", "k", rawKeys, "additional keys to include in the generated authorized_keys")
	rootCmd.PersistentFlags().StringVar(&tshInstanceURL, "url", transfer.DefaultTransferBaseURL, "url of the target transfer.sh instance")
	rootCmd.Flags().IntVarP(&transferDownloads, "max-downloads", "m", 10, "maximum number of times any content shared can be downloaded")
	rootCmd.Flags().IntVarP(&transferDays, "max-days", "d", 2, "number of days that the content will remain available via transfer.sh")
	rootCmd.Flags().StringVarP(&password, "encrypt", "e", password, "password for transfer.sh server-side encryption")
	rootCmd.AddCommand(deleteCmd)
}
