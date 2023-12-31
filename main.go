package ldap

import (
	"context"
	"crypto/tls"
	"fmt"
	"strings"

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

func filter_error(err error) string {
	return strings.Replace(err.Error(), "0x00", "", -1)
}

func Run(ctx context.Context, target string, command string, expectedOutput string, username string, password string, options map[string]interface{}) (bool, string) {
	optionsStruct := optionsStruct{}
	optionsStruct.Marshal(options)

	if !strings.Contains(target, ":") {
		if optionsStruct.LDAPS {
			target = fmt.Sprintf("%s:636", target)
		} else {
			target = fmt.Sprintf("%s:389", target)
		}
	}

	errChan := make(chan error)

	go func() {
		defer close(errChan)

		client, err := ldap.Dial("tcp", target)
		if err != nil {
			errChan <- fmt.Errorf("failed to connect to LDAP server: %s", filter_error(err))
			return
		}
		defer client.Close()

		if optionsStruct.LDAPS {
			err = client.StartTLS(&tls.Config{InsecureSkipVerify: true})
			if err != nil {
				errChan <- fmt.Errorf("failed to start TLS: %s", filter_error(err))
				return
			}
		}

		err = client.Bind(username, password)
		if err != nil {
			errChan <- fmt.Errorf("failed to bind to LDAP server: %s", filter_error(err))
			return
		}

		errChan <- nil

	}()

	select {
	case <-ctx.Done():
		return false, fmt.Sprintf("check timed out: %s", filter_error(ctx.Err()))
	case err := <-errChan:
		if err == nil {
			return true, ""
		}
		return false, fmt.Sprintf("check encountered an error: %s", filter_error(err))
	}
}
