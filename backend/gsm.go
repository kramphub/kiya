package backend

import (
	"context"
	"fmt"
	"strings"

	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"google.golang.org/api/iterator"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GSM struct {
	client *secretmanager.Client
}

func NewGSM(client *secretmanager.Client) *GSM {
	return &GSM{client: client}
}

func (b *GSM) Get(ctx context.Context, p *Profile, key string) ([]byte, error) {
	if p == nil {
		return nil, fmt.Errorf("provided profile cannot be nil")
	}

	result, err := b.client.AccessSecretVersion(ctx, &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf(
			"projects/%s/secrets/%s/versions/%s",
			p.ProjectID,
			key,
			"latest",
		),
	})
	if err != nil {
		return nil, err
	}

	if result.Payload == nil || result.Payload.Data == nil {
		return nil, fmt.Errorf("failed to get secret from GSM, a nil result was returned")
	}

	return result.Payload.Data, nil
}

func (b *GSM) List(ctx context.Context, p *Profile) ([]Key, error) {
	it := b.client.ListSecrets(ctx, &secretmanagerpb.ListSecretsRequest{
		Parent: fmt.Sprintf("projects/%s", p.ProjectID),
	})

	var keys []Key
	for {
		secret, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to list secrets from GSM, %w", err)
		}

		keys = append(keys, Key{
			Name:      b.fullNameToName(secret.Name),
			CreatedAt: secret.CreateTime.AsTime(),
			Info:      "creator: <Unknown>", // no owner
			Owner:     "<Unknown>",
		})
	}

	return keys, nil
}

func (b *GSM) CheckExists(ctx context.Context, p *Profile, key string) (bool, error) {
	_, err := b.Get(ctx, p, key)
	return err == nil, err
}

func (b *GSM) Put(ctx context.Context, p *Profile, key, value string) error {
	_, err := b.client.CreateSecret(ctx, &secretmanagerpb.CreateSecretRequest{
		Parent:   fmt.Sprintf("projects/%s", p.ProjectID),
		SecretId: key,
		Secret: &secretmanagerpb.Secret{
			Replication: &secretmanagerpb.Replication{
				Replication: &secretmanagerpb.Replication_Automatic_{},
			},
		},
	})
	if err != nil {
		statusErr, ok := status.FromError(err)
		if !ok || statusErr.Code() != codes.AlreadyExists {
			return fmt.Errorf("failed to create secret in GSM, %w", err)
		}
	}

	_, err = b.client.AddSecretVersion(ctx, &secretmanagerpb.AddSecretVersionRequest{
		Parent: fmt.Sprintf("projects/%s/secrets/%s", p.ProjectID, key),
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte(value),
		},
	})
	if err != nil {
		return fmt.Errorf("failed to add version after creating the secret in GSM, %w", err)
	}

	return nil
}

func (b *GSM) Delete(ctx context.Context, p *Profile, key string) error {
	err := b.client.DeleteSecret(ctx, &secretmanagerpb.DeleteSecretRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s", p.ProjectID, key),
	})
	if err != nil {
		return fmt.Errorf("failed to delete secret from GSM, %w", err)
	}

	return nil
}

func (b *GSM) Close() error {
	return b.client.Close()
}

// SetMasterPassword is not relevant for this backend
func (b *GSM) SetMasterPassword(_ []byte) {
	// noop
}

///

func (b *GSM) fullNameToName(fullName string) string {
	return fullName[strings.LastIndex(fullName, "/")+1:]
}
