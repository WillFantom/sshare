package github

import (
	"context"
	"fmt"

	gh "github.com/google/go-github/github"
	"golang.org/x/oauth2"

	"github.com/willfantom/sshare/keys"
)

// GitHubAgent holds a GitHub API client that has been given a user token.
type GitHubAgent struct {
	client *gh.Client
}

// NewAgent creates a new SSH agent like entity, providing access to the
// authenticated user's public keys and titles of them.
func NewAgent(token string) *GitHubAgent {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &GitHubAgent{
		client: gh.NewClient(tc),
	}
}

// GetKeys returns the SSH keys of the authenticated GitHub user as
// authorized_key style items. If this failed for any reason, such as not being
// able to contact GitHub or the token not being valid, an error is returned.
func (gha *GitHubAgent) GetKeys() ([]*keys.Key, error) {
	ghKeys, _, err := gha.client.Users.ListKeys(context.Background(), "", &gh.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to obtain users ssh keys from github: %w", err)
	}
	allKeys := make([]*keys.Key, len(ghKeys))
	for idx, ghk := range ghKeys {
		k, err := keys.NewKey(*ghk.Key, *ghk.Title)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key from github: %w", err)
		}
		allKeys[idx] = k
	}
	return allKeys, nil
}
