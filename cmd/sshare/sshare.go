package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/willfantom/sshare/internal/transfer"
	"github.com/willfantom/sshare/internal/ui"
	"github.com/willfantom/sshare/keys"
)

var (
	transferDownloads int      = 0
	transferDays      int      = 0
	sshAgentPath      string   = os.Getenv("SSH_AUTH_SOCK")
	sshAgentPass      string   = ""
	keyFilepaths      []string = make([]string, 0)
)

var (
	rootCmd = &cobra.Command{
		Use:   "sshare",
		Short: "Easily share links to your SSH public keys",
		Long:  `Share your public SSH keys found in your agent via curl-able transfer.sh links.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if transferDays <= 0 {
				ui.Errorln("Uploaded content must be visible on transfer.sh for at least 1 day", true)
			}
			if transferDownloads <= 0 {
				ui.Errorln("Uploaded content must be downloadable via transfer.sh for at least 1 time", true)
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			chosenKeys := make([]*keys.Key, 0)
			for _, fp := range keyFilepaths {
				k, err := keys.NewKeyFromFile(fp)
				if err != nil {
					ui.Errorln(err.Error(), true)
				}
				chosenKeys = append(chosenKeys, k)
			}
			if sshAgentPath != "" {
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
)

func main() {
	if err := rootCmd.Execute(); err != nil {
		ui.Errorln(err.Error(), true)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&sshAgentPath, "agent", "a", sshAgentPath, "path to the target ssh Agent socket")
	rootCmd.PersistentFlags().StringVarP(&sshAgentPass, "passphrase", "p", sshAgentPass, "passphrase for the ssh agent")
	rootCmd.PersistentFlags().StringArrayVarP(&keyFilepaths, "key-file", "k", keyFilepaths, "additional key file(s) to include in the generated authorized_keys")
	rootCmd.PersistentFlags().IntVarP(&transferDownloads, "max-downloads", "m", 10, "maximum number of times any content shared can be downloaded")
	rootCmd.PersistentFlags().IntVarP(&transferDays, "max-days", "d", 2, "number of days that the content will remain available via transfer.sh")
}
