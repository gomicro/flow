package session

import (
	"fmt"
	"net/http"

	aws_session "github.com/aws/aws-sdk-go/aws/session"
)

func NewSession() (*aws_session.Session, error) {
	sess, err := aws_session.NewSession()
	if err != nil {
		return nil, fmt.Errorf("Error creating session: %w", err)
	}

	sess.Config.HTTPClient = &http.Client{}

	return sess, nil
}
