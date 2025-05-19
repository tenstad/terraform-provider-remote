package provider

import (
	"errors"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"golang.org/x/crypto/ssh"
)

var (
	errTypecast = errors.New("typecast error")
	errGet      = errors.New("get error")
)

func Get[T any](d *schema.ResourceData, key string) (T, error) {
	value, ok, err := GetOk[T](d, key)
	if !ok {
		var t T
		return t, fmt.Errorf("%w: %s is undefined", errGet, key)
	}
	return value, err
}

func GetOk[T any](d *schema.ResourceData, key string) (T, bool, error) {
	raw, ok := d.GetOk(key)
	if !ok {
		var t T
		return t, false, nil
	}

	if value, ok := raw.(T); ok {
		return value, true, nil
	}

	var t T
	return t, true, fmt.Errorf("%w: %s to %T: %v", errTypecast, key, t, raw)
}

func parsePrivateKey(d *schema.ResourceData, privateKey string) (ssh.Signer, error) {
	privateKeyPass, ok, err := GetOk[string](d, "conn.0.private_key_pass")
	if ok {
		if err != nil {
			return nil, err
		}
		return ssh.ParsePrivateKeyWithPassphrase([]byte(privateKey), []byte(privateKeyPass))
	}
	signer, err := ssh.ParsePrivateKey([]byte(privateKey))
	if _, ok := err.(*ssh.PassphraseMissingError); ok {
		return nil, fmt.Errorf("private key is encrypted, password required but not supplied")
	}
	return signer, err
}

func writeFileToHost(host string, filename string, content string, group string, user string) {
	sshClient, err := ssh.Dial("tcp", host, &ssh.ClientConfig{
		User:            "root",
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Auth:            []ssh.AuthMethod{ssh.Password("password")},
	})
	if err != nil {
		panic(err)
	}

	session, err := sshClient.NewSession()
	if err != nil {
		panic(err)
	}
	defer session.Close()

	stdin, err := session.StdinPipe()
	if err != nil {
		panic(err)
	}

	go func() {
		stdin.Write([]byte(content))
		stdin.Close()
	}()
	session.Run(fmt.Sprintf("cat /dev/stdin | tee %s && chgrp %s %s && chown %s %s", filename, group, filename, user, filename))
}
