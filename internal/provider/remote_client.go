package provider

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type RemoteClient struct {
	sshClient *ssh.Client
}

func (c *RemoteClient) WriteFile(content string, path string, permissions string) error {
	scpClient, err := c.GetSCPClient()
	if err != nil {
		return err
	}
	defer scpClient.Close()

	return scpClient.CopyFile(strings.NewReader(content), path, permissions)
}

func (c *RemoteClient) WriteFileSudo(content string, path string) error {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		return err
	}

	go func() {
		stdin.Write([]byte(content))
		stdin.Close()
	}()

	cmd := fmt.Sprintf("cat /dev/stdin | sudo tee %s", path)
	return session.Run(cmd)
}

func (c *RemoteClient) ChmodFileSudo(path string, permissions string) error {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo chmod %s %s", permissions, path)
	return session.Run(cmd)
}

func (c *RemoteClient) ChgrpFileSudo(path string, group string) error {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo chgrp %s %s", group, path)
	return session.Run(cmd)
}

func (c *RemoteClient) ChownFileSudo(path string, owner string) error {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo chown %s %s", owner, path)
	return session.Run(cmd)
}

func (c *RemoteClient) FileExistsSudo(path string) (bool, error) {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	cmd := fmt.Sprintf("test -f %s", path)
	err = session.Run(cmd)

	if err != nil {
		session2, err := sshClient.NewSession()
		if err != nil {
			return false, err
		}
		defer session2.Close()

		cmd := fmt.Sprintf("test ! -f %s", path)
		return false, session2.Run(cmd)
	}

	return true, nil
}

func (c *RemoteClient) ReadFile(path string) (string, error) {
	sftpClient, err := c.GetSFTPClient()
	if err != nil {
		return "", err
	}
	defer sftpClient.Close()

	file, err := sftpClient.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content := bytes.Buffer{}
	_, err = file.WriteTo(&content)
	if err != nil {
		return "", err
	}

	return content.String(), nil
}

func (c *RemoteClient) ReadFileSudo(path string) (string, error) {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo cat %s", path)
	content, err := session.Output(cmd)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func (c *RemoteClient) DeleteFile(path string) error {
	sftpClient, err := c.GetSFTPClient()
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	return sftpClient.Remove(path)
}

func (c *RemoteClient) DeleteFileSudo(path string) error {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo rm %s", path)
	return session.Run(cmd)
}

func NewRemoteClient(host string, clientConfig *ssh.ClientConfig) (*RemoteClient, error) {
	client, err := ssh.Dial("tcp", host, clientConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish a connection to the remote server: %s", err.Error())
	}

	return &RemoteClient{
		sshClient: client,
	}, nil
}

func (c *RemoteClient) Close() error {
	return c.sshClient.Close()
}

func (c *RemoteClient) GetSSHClient() *ssh.Client {
	return c.sshClient
}

func (c *RemoteClient) GetSCPClient() (scp.Client, error) {
	return scp.NewClientBySSH(c.sshClient)
}

func (c *RemoteClient) GetSFTPClient() (*sftp.Client, error) {
	return sftp.NewClient(c.sshClient)
}
