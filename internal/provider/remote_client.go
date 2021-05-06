package provider

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/bramvdbogaerde/go-scp"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type RemoteClient struct {
	sshClient *ssh.Client
}

func (c *RemoteClient) WriteFile(d *schema.ResourceData) error {
	scpClient, err := c.GetSCPClient()
	if err != nil {
		return err
	}
	defer scpClient.Close()

	return scpClient.CopyFile(strings.NewReader(d.Get("content").(string)), d.Get("path").(string), d.Get("permissions").(string))
}

func (c *RemoteClient) WriteFileSudo(d *schema.ResourceData) error {
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

	content := d.Get("content").(string)
	go func() {
		stdin.Write([]byte(content))
		stdin.Close()
	}()

	cmd := fmt.Sprintf("cat /dev/stdin | sudo tee %s", d.Get("path").(string))
	return session.Run(cmd)
}

func (c *RemoteClient) ChmodFileSudo(d *schema.ResourceData) error {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo chmod %s %s", d.Get("permissions").(string), d.Get("path").(string))
	return session.Run(cmd)
}

func (c *RemoteClient) FileExistsSudo(d *schema.ResourceData) (bool, error) {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return false, err
	}
	defer session.Close()

	path := d.Get("path").(string)
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

func (c *RemoteClient) ReadFile(d *schema.ResourceData) error {
	sftpClient, err := c.GetSFTPClient()
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	file, err := sftpClient.Open(d.Get("path").(string))
	if err != nil {
		return err
	}
	defer file.Close()

	content := bytes.Buffer{}
	_, err = file.WriteTo(&content)

	if err != nil {
		return err
	}

	d.Set("content", string(content.String()))
	return nil
}

func (c *RemoteClient) ReadFileSudo(d *schema.ResourceData) error {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo cat %s", d.Get("path").(string))
	content, err := session.Output(cmd)
	if err != nil {
		return err
	}

	d.Set("content", string(content))
	return nil
}

func (c *RemoteClient) DeleteFile(d *schema.ResourceData) error {
	sftpClient, err := c.GetSFTPClient()
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	return sftpClient.Remove(d.Get("path").(string))
}

func (c *RemoteClient) DeleteFileSudo(d *schema.ResourceData) error {
	sshClient := c.GetSSHClient()

	session, err := sshClient.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	cmd := fmt.Sprintf("sudo rm %s", d.Get("path").(string))
	return session.Run(cmd)
}

func NewRemoteClient(host string, clientConfig ssh.ClientConfig) (*RemoteClient, error) {
	client, err := ssh.Dial("tcp", host, &clientConfig)
	if err != nil {
		return nil, fmt.Errorf("couldn't establish a connection to the remote server: %s", err.Error())
	}

	return &RemoteClient{
		sshClient: client,
	}, nil
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
