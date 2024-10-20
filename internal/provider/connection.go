package provider

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
)

type ConnectionResourceModel struct {
	Host             types.String `tfsdk:"host"`
	Port             types.Int64  `tfsdk:"port"`
	Timeout          types.Int64  `tfsdk:"timeout"`
	User             types.String `tfsdk:"user"`
	Sudo             types.Bool   `tfsdk:"sudo"`
	Agent            types.Bool   `tfsdk:"agent"`
	Password         types.String `tfsdk:"password"`
	PrivateKey       types.String `tfsdk:"private_key"`
	PrivateKeyPath   types.String `tfsdk:"private_key_path"`
	PrivateKeyEnvVar types.String `tfsdk:"private_key_env_var"`
}

// applyDefaults applies defauls identical to the ones defined in
// file resource scheme, as they cannot be defined in provider
// and data source schemes.
func (m *ConnectionResourceModel) applyDefaults() {
	if m.Port.IsNull() {
		m.Port = types.Int64Value(22)
	}
	if m.Sudo.IsNull() {
		m.Sudo = types.BoolValue(false)
	}
	if m.Agent.IsNull() {
		m.Agent = types.BoolValue(false)
	}
}

func (d *ConnectionResourceModel) ConnectionHash() string {
	elements := []string{
		d.Host.ValueString(),
		d.User.ValueString(),
		strconv.Itoa(int(d.Port.ValueInt64())),
		d.Password.ValueString(),
		d.PrivateKey.ValueString(),
		d.PrivateKeyPath.ValueString(),
		strconv.FormatBool(d.Agent.ValueBool()),
	}
	return strings.Join(elements, "::")
}

func (conn *ConnectionResourceModel) Connection() (string, *ssh.ClientConfig, diag.Diagnostics) {
	dia := diag.Diagnostics{}

	host := conn.Host.ValueString()
	port := conn.Port.ValueInt64()
	user := conn.User.ValueString()

	clientConfig := ssh.ClientConfig{
		User:            user,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	if !conn.Password.IsNull() && !conn.Password.IsUnknown() {
		clientConfig.Auth = append(clientConfig.Auth, ssh.Password(conn.Password.ValueString()))
	}

	if !conn.PrivateKey.IsNull() && !conn.PrivateKey.IsUnknown() {
		signer, err := ssh.ParsePrivateKey([]byte(conn.PrivateKey.ValueString()))
		if err != nil {
			dia.AddError("Error", fmt.Sprintf("couldn't create a ssh client config from private key: %s", err.Error()))
		} else {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
		}
	}

	if !conn.PrivateKeyPath.IsNull() && !conn.PrivateKeyPath.IsUnknown() {
		content, err := os.ReadFile(conn.PrivateKeyPath.ValueString())
		if err != nil {
			dia.AddError("Error", fmt.Sprintf("couldn't read private key: %s", err.Error()))
		} else {
			signer, err := ssh.ParsePrivateKey(content)
			if err != nil {
				dia.AddError("Error", fmt.Sprintf("couldn't create a ssh client config from private key file: %s", err.Error()))
			} else {
				clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
			}
		}
	}

	if !conn.PrivateKeyEnvVar.IsNull() && !conn.PrivateKeyEnvVar.IsUnknown() {
		content := []byte(os.Getenv(conn.PrivateKeyEnvVar.ValueString()))
		signer, err := ssh.ParsePrivateKey(content)
		if err != nil {
			dia.AddError("Error", fmt.Sprintf("couldn't create a ssh client config from private key env var: %s", err.Error()))
		} else {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeys(signer))
		}
	}
	if conn.Agent.ValueBool() {
		connection, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK"))
		if err != nil {
			dia.AddError("Error", fmt.Sprintf("couldn't connect to SSH agent: %s", err.Error()))
		} else {
			clientConfig.Auth = append(clientConfig.Auth, ssh.PublicKeysCallback(agent.NewClient(connection).Signers))
		}
	}

	if !conn.Timeout.IsNull() && !conn.Timeout.IsUnknown() {
		clientConfig.Timeout = time.Duration(conn.Timeout.ValueInt64()) * time.Millisecond
	}

	return fmt.Sprintf("%s:%d", host, port), &clientConfig, dia
}
