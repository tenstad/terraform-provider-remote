package provider

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

func writeFileToHost(host string, filename string, content string) {
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
	session.Run(fmt.Sprintf("cat /dev/stdin | tee %s", filename))
}
