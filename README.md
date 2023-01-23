# SSHare    ![GitHub release (latest SemVer)](https://img.shields.io/github/v/tag/willfantom/sshare?display_name=tag&label=%20&sort=semver)  ![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/willfantom/sshare/release.yml?label=%20&logo=github)

Quickly share SSH keys from your agent as `curl`able links! Use the `sshare` TUI to easily select keys from an SSH agent and automatically generate a [`transfer.sh`](https://transfer.sh) upload that can easily be added to an authorized key file via `curl`. 

---

## Install 

 - **Linux**
   - Download the appropriate version for your system from the [GitHub release page](https://github.com/WillFantom/sshare/releases)
   - Move the `sshare` binary to a directory in your `PATH`
 - **MacOS**
   - Download and install via Homebrew: 
      ```
      brew install willfantom/tap/sshare
      ```

## Usage

Simply run: `sshare`

Sensible defaults are used, such as connecting to the SSH agent as defined by the `SSH_AUTH_SOCK` environment variable.

- To include keys that are not found in your agent, you can specify a file containing the public key with the `-f` flag (and can be specified multiple times):

  ```
  sshare -f ~/.ssh/id_rsa.pub
  ```

- An SSH agent socket path other than the one found in `SSH_AUTH_SOCK` can be provided using the `-a` flag:

  ```
  sshare -a /tmp/ssh-XXXXXXanCbmG/agent.8
  ```

- Add keys from a generated link to an `authorized_keys` file:

  ```
  curl --fail-with-body https://transfer.sh/<unique_code>/authorized_keys >> ~/.ssh/authorized_keys
  ```

- Uploads can be deleted by either:
  - Opening the downloaded link in a browser and deleting the key using the given deletion code
  - Or by running:
    ```
    curl -X DELETE <downloadURL>/<deletionCode>
    ```

- Also provide your GitHub SSH keys as options for keys to share:

  ```
  sshare -g "<a github access token with read ssh keys permission>"
  ```

---

```
Share your public SSH keys found in your agent via curl-able transfer.sh links.

Usage:
  sshare [flags]

Flags:
  -a, --agent string           path to the target ssh agent socket ($SSH_AUTH_SOCK) (default "~/.1password/agent.sock")
  -g, --github-token string    github token with permission to read ssh keys
  -h, --help                   help for sshare
  -k, --key stringArray        additional keys to include in the generated authorized_keys
  -f, --key-file stringArray   additional key file(s) to include in the generated authorized_keys
  -d, --max-days int           number of days that the content will remain available via transfer.sh (default 2)
  -m, --max-downloads int      maximum number of times any content shared can be downloaded (default 10)
  -p, --passphrase string      passphrase for the ssh agent
```
