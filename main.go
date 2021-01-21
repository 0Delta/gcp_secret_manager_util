package gcp_secret_manager_util

import (
	"context"
	"log"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

type SecretData struct {
	Secret string `json:"secret"`
	data   []byte
}

var (
	ctx    context.Context
	client *secretmanager.Client
)

func setup() (err error) {
	if ctx != nil && client != nil {
		return
	}
	// Create the client.
	ctx = context.Background()
	client, err = secretmanager.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to setup client: %v", err)
		return err
	}
	return nil
}

func (s *SecretData) Get() []byte {
	if s.data == nil {
		_ = s.decrypt()
	}
	return s.data
}

func (s *SecretData) String() string {
	return string(s.Get())
}

func (s *SecretData) decrypt() (err error) {
	setup()
	// Build the request.
	accessRequest := &secretmanagerpb.AccessSecretVersionRequest{
		Name: s.Secret,
	}

	// Call the API.
	result, err := client.AccessSecretVersion(ctx, accessRequest)
	if err != nil {
		log.Fatalf("failed to access secret version: %v", err)
		return err
	}
	s.data = make([]byte, len(result.Payload.Data))
	_ = copy(s.data, result.Payload.Data)
	return nil
}
