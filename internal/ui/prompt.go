package ui

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/willfantom/sshare/keys"
)

// SelectKey runs a terminal selection prompt that allows a user to select keys
// from a given slice. Returned is a slice of selected keys or an error is the
// prompt exits.
func SelectKey(authKeys []*keys.Key) ([]*keys.Key, error) {
	keyOptions := make([]string, len(authKeys))
	for idx, k := range authKeys {
		keyOptions[idx] = fmt.Sprintf("%d. %s (%s)", idx+1, k.Name(), k.Fingerprint())
	}
	ms := survey.MultiSelect{
		Message:  "Select the keys that you would like to share",
		Options:  keyOptions,
		PageSize: 5,
	}
	var chosenKeyOptions []string
	if err := survey.AskOne(&ms, &chosenKeyOptions); err != nil {
		return nil, fmt.Errorf("prompt exited")
	}
	if len(chosenKeyOptions) <= 0 {
		return nil, fmt.Errorf("no keys were selected")
	}
	chosenKeys := make([]*keys.Key, len(chosenKeyOptions))
	for idx, ck := range chosenKeyOptions {
		keyIndex, err := strconv.Atoi(strings.SplitN(ck, ".", 2)[0])
		if err != nil {
			return nil, fmt.Errorf("failed to obtain key index for a selected key")
		}
		chosenKeys[idx] = authKeys[keyIndex-1]
	}
	return chosenKeys, nil
}
