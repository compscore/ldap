package ldap

import (
	"context"
	"crypto/tls"
	"fmt"

	"github.com/go-ldap/ldap/v3"
)

type optionsStruct struct {
	// Use LDAPS instead of LDAP
	LDAPS bool `compscore:"ldaps"`
}

func (o *optionsStruct) Marshal(options map[string]interface{}) {
	ldapsInterface, ok := options["ldaps"]
	if ok {
		ldaps, ok := ldapsInterface.(bool)
		if ok {
			o.LDAPS = ldaps
		}
	}
}

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string, options map[string]interface{}) (bool, string) {
	optionsStruct := optionsStruct{}
	optionsStruct.Marshal(options)

	errChan := make(chan error)

	go func() {
		defer close(errChan)

		client, err := ldap.Dial("tcp", target)
		if err != nil {
			errChan <- fmt.Errorf("failed to connect to LDAP server: %s", err.Error())
			return
		}
		defer client.Close()

		if optionsStruct.LDAPS {
			err = client.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				errChan <- fmt.Errorf("failed to start TLS: %s", err.Error())
				return
			}
		}

		err = client.Bind(username, password)
		if err != nil {
			errChan <- fmt.Errorf("failed to bind to LDAP server: %s", err.Error())
			return
		}

		errChan <- nil

	}()

	select {
	case <-ctx.Done():
		return false, fmt.Sprintf("check timed out: %s", ctx.Err().Error())
	case err := <-errChan:
		if err == nil {
			return true, ""
		}
		return false, fmt.Sprintf("check encountered an error: %s", err.Error())
	}
}
