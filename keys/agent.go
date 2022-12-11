package keys

import (
	"fmt"
	"net"

	"golang.org/x/crypto/ssh/agent"
)

// Agent is a connection to an SSH agent, along with some supplementray data for
// manageing the connection.
type Agent struct {
	conn       net.Conn
	sshAgent   agent.ExtendedAgent
	path       string
	passphrase string
}

// AgentOpt is a configuration option for an SSH agent connection.
type AgentOpt func(*Agent)

// NewSSHAgent establishes a connection to a unix socket at the given path
// representing an SSH agent. Options can be given to configure the connection
// and how the agent is used. Options are executed in the order provided.
// Returned is an SSH agent connection or an error if unsuccessful.
func NewSSHAgent(path string, opts ...AgentOpt) (*Agent, error) {
	conn, err := net.Dial("unix", path)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to agent socket: %w", err)
	}
	client := agent.NewClient(conn)
	if client == nil {
		return nil, fmt.Errorf("socket could not be as ssh agent: %w", err)
	}
	sshAgent := Agent{
		conn:       conn,
		sshAgent:   client,
		path:       path,
		passphrase: "",
	}
	for _, opt := range opts {
		opt(&sshAgent)
	}
	return &sshAgent, nil
}

// AgentPassphraseOpt should be provided when creating a new agent if the SSH
// agent is locked by a passphrase.
func AgentPassphraseOpt(passphrase string) AgentOpt {
	return func(a *Agent) {
		a.passphrase = passphrase
	}
}

// Close closes the underlying socket connection used to communicate with the
// SSH agent.
func (a *Agent) Close() error {
	return a.conn.Close()
}

// GetKeys uses the SSH agent connection to retreive all the keys present in the
// agent in the authorized key format. If the agent is locked and the correct
// passphrase option has been provided, the agent will be unlocked and re-locked
// as required. Returned are the keys present in the agent or an error if
// unsuccessful.
func (a *Agent) GetKeys() ([]*Key, error) {
	if a.passphrase != "" {
		if err := a.unlock(); err != nil {
			return nil, err
		}
		//TODO: Check for errors when re-locking the agent
		defer a.lock()
	}
	agentKeys, err := a.sshAgent.List()
	if err != nil {
		return nil, fmt.Errorf("failed to get keys from agent: %w", err)
	}
	keys := make([]*Key, len(agentKeys))
	for idx, k := range agentKeys {
		key, err := NewKey(k.String(), k.Comment)
		if err != nil {
			return nil, fmt.Errorf("failed to parse key from agent: %w", err)
		}
		keys[idx] = key
	}
	return keys, nil
}

func (a *Agent) unlock() error {
	if err := a.sshAgent.Unlock([]byte(a.passphrase)); err != nil {
		return fmt.Errorf("failed to unlock ssh agent: %w", err)
	}
	return nil
}

func (a *Agent) lock() error {
	if err := a.sshAgent.Lock([]byte(a.passphrase)); err != nil {
		return fmt.Errorf("failed to lock ssh agent: %w", err)
	}
	return nil
}
