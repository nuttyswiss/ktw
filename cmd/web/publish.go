package main

import (
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	slashpath "path"
	"path/filepath"

	"github.com/pkg/sftp"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"golang.org/x/crypto/ssh/knownhosts"
)

func init() {
	var cmd = &cobra.Command{
		Use:   "publish",
		Short: "Publish the website",
		RunE:  publish,
	}
	cli.AddCommand(cmd)
}

func idFile(file string) (ssh.AuthMethod, error) {
	buf, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	key, err := ssh.ParsePrivateKey(buf)
	if err != nil {
		return nil, err
	}
	return ssh.PublicKeys(key), nil
}

func homedir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func sshdir() string {
	if h := os.Getenv("SSH_DIR"); h != "" {
		return h
	}
	return filepath.Join(homedir(), ".ssh")
}

type endpoint struct {
	client  *sftp.Client
	root    string
	closers []io.Closer
}

// authMethods returns a slice of all available authentication methods.
func (e *endpoint) authMethods() ([]ssh.AuthMethod, error) {
	methods := []ssh.AuthMethod{}

	// 1. Identity files
	for _, keyFile := range []string{"id_rsa", "id_dsa", "id_ecdsa", "id_ed25519"} {
		filePath := filepath.Join(sshdir(), keyFile)
		if _, err := os.Stat(filePath); err == nil { // Check if file exists
			method, err := idFile(filePath)
			if err != nil {
				return nil, fmt.Errorf("failed to load key file %s: %w", filePath, err)
			}
			methods = append(methods, method)
		}
	}

	// 2. SSH agent
	if agentConn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		methods = append(methods, ssh.PublicKeysCallback(agent.NewClient(agentConn).Signers))
		e.closers = append(e.closers, agentConn)
	}

	return methods, nil
}

// Close all open connections.
func (e *endpoint) Close() error {
	for _, c := range e.closers {
		if err := c.Close(); err != nil {
			return err
		}
	}
	return nil
}

// Client returns the SFTP client.
func (e *endpoint) Client() *sftp.Client {
	return e.client
}

func (e *endpoint) FullPath(path string) string {
	return slashpath.Join(e.root, path)
}

// NewEndpoint creates a new endpoint.
func NewEndpoint(dest string) (*endpoint, error) {
	e := &endpoint{}

	baseURL, err := url.Parse(dest)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	if baseURL.Scheme != "sftp" {
		return nil, fmt.Errorf("unsupported URL scheme: %s", baseURL.Scheme)
	}
	e.root = baseURL.Path

	user := os.Getenv("USER")
	if baseURL.User.Username() != "" {
		user = baseURL.User.Username()
	}

	auth, err := e.authMethods()
	if err != nil {
		e.Close()
		return nil, fmt.Errorf("failed to get auth methods: %w", err)
	}

	hostcheck, err := knownhosts.New(filepath.Join(sshdir(), "known_hosts"))
	if err != nil {
		return nil, fmt.Errorf("failed to create known hosts callback: %w", err)
	}

	// Create SSH config using agent for authentication
	sshConfig := &ssh.ClientConfig{
		User:            user,
		Auth:            auth,
		HostKeyCallback: hostcheck,
	}

	// Connect to the remote server
	host := baseURL.Hostname()
	port := baseURL.Port()
	if port == "" {
		port = "22"
	}
	client, err := ssh.Dial("tcp", net.JoinHostPort(host, port), sshConfig)
	if err != nil {
		e.Close()
		return nil, fmt.Errorf("failed to dial: %w", err)
	}
	e.closers = append(e.closers, client)

	// Create a new SFTP client
	sftpClient, err := sftp.NewClient(client)
	if err != nil {
		e.Close()
		return nil, fmt.Errorf("failed to create SFTP client: %w", err)
	}
	e.closers = append(e.closers, sftpClient)
	e.client = sftpClient

	return e, nil
}

func publish(cmd *cobra.Command, args []string) error {
	dest := viper.GetString("root")
	if dest == "" {
		return fmt.Errorf("config is missing 'root' key")
	}
	endpoint, err := NewEndpoint(dest)
	if err != nil {
		return err
	}
	defer endpoint.Close()
	client := endpoint.Client()

	// Specify the source directory to upload
	root := viper.GetString("dir")
	if root == "" {
		return fmt.Errorf("config is missing 'dir' key")
	}
	fmt.Printf("Generating from %s\n", root)

	cwd, err := client.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current working directory: %w", err)
	}
	_ = cwd

	// Walk the source directory and upload each file
	err = filepath.Walk(root, func(src string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		srcpath, err := filepath.Rel(root, src)
		if err != nil {
			return err
		}

		// Open the local file
		localFile, err := os.Open(src)
		if err != nil {
			return fmt.Errorf("failed to open local file: %w", err)
		}
		defer localFile.Close()

		// Create the remote file and directory paths
		remotePath := endpoint.FullPath(srcpath)
		remoteDir := slashpath.Dir(remotePath)

		// Create the remote directory if it doesn't exist
		if err := client.MkdirAll(remoteDir); err != nil {
			return fmt.Errorf("failed to create remote directory: %w", err)
		}

		// Create the remote file
		remoteFile, err := client.Create(remotePath)
		if err != nil {
			return fmt.Errorf("failed to create remote file: %w", err)
		}
		defer remoteFile.Close()

		// Copy the local file contents to the remote file
		if _, err := io.Copy(remoteFile, localFile); err != nil {
			return fmt.Errorf("failed to copy file contents: %w", err)
		}

		fmt.Printf("Uploaded: %s\n", srcpath)
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to walk directory: %w", err)
	}

	fmt.Println("Website published successfully!")
	return nil
}
