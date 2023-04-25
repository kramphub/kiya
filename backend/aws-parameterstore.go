package backend

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
)

// AWSParameterStore implements Backend for AWS Parameter Store service.
type AWSParameterStore struct {
	client   *ssm.Client
	kmsKeyID string
}

// NewAWSParameterStore returns a new AWSParameterStore with an initialized AWS SSM client.
func NewAWSParameterStore(ctx context.Context, p *Profile) (*AWSParameterStore, error) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	return &AWSParameterStore{
		client:   ssm.NewFromConfig(cfg),
		kmsKeyID: p.CryptoKey}, nil
}

// Get returns the decrypted value for a parameter by key.
func (s *AWSParameterStore) Get(ctx context.Context, p *Profile, key string) ([]byte, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}
	output, err := s.client.GetParameter(ctx, input)
	if err != nil {
		return []byte{}, err
	}

	return []byte(*output.Parameter.Value), nil
}

// List returns all keys available.
func (s *AWSParameterStore) List(ctx context.Context, p *Profile) (list []Key, err error) {
	input := &ssm.GetParametersByPathInput{
		Path:       aws.String("/"),
		MaxResults: aws.Int32(10), // is the documented maximum
		Recursive:  aws.Bool(true),
	}
	for {
		output, err := s.client.GetParametersByPath(ctx, input)
		if err != nil {
			return []Key{}, err
		}
		for _, each := range output.Parameters {
			list = append(list, Key{
				Name:      *each.Name,
				CreatedAt: *each.LastModifiedDate,
				Info:      fmt.Sprintf("type: %s datatype: %s version: %d", *each.DataType, *each.DataType, each.Version),
				Owner:     "<Unknown>",
			})
		}
		if output.NextToken != nil {
			input.NextToken = output.NextToken
		} else {
			break
		}
	}
	return
}

// CheckExists returns true if there exists a value for a given key.
func (s *AWSParameterStore) CheckExists(ctx context.Context, p *Profile, key string) (bool, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(false), // No decryption is needed
	}
	_, err := s.client.GetParameter(ctx, input)
	// if _, ok := err.(*ssm.); ok {
	// 	return false, nil
	// }
	// other error?
	if err != nil {
		return false, err
	}
	return true, nil

}

// Put write the parameter and its value using encryption ;either the default key or the one specified in the profile.
func (s *AWSParameterStore) Put(ctx context.Context, p *Profile, key, value string, overwrite bool) error {
	input := &ssm.PutParameterInput{
		Name:      aws.String(key),
		Value:     aws.String(value),
		Overwrite: aws.Bool(overwrite),
		DataType:  aws.String("text"),
		Type:      types.ParameterTypeSecureString,
	}
	if !overwrite {
		input.Description = aws.String(fmt.Sprintf("created by %s using kiya", os.Getenv("USER")))
		input.Tags = []types.Tag{{Key: aws.String("creator"), Value: aws.String(os.Getenv("USER"))}}
	}
	// only if CryptoKey is set in the Profile then we set the KeyId
	// which overrides the default key associated with the AWS account
	if p.CryptoKey != "" {
		input.KeyId = aws.String(s.kmsKeyID)
	}
	_, err := s.client.PutParameter(ctx, input)
	if err != nil {
		return err
	}
	return nil
}

// Delete removes the parameter by its key
func (s *AWSParameterStore) Delete(ctx context.Context, p *Profile, key string) error {
	input := &ssm.DeleteParameterInput{
		Name: aws.String(key),
	}
	_, err := s.client.DeleteParameter(ctx, input)
	return err
}

// Close is not effictive for this backend
func (s *AWSParameterStore) Close() error {
	// noop
	return nil
}

func (s *AWSParameterStore) SetParameter(key string, value interface{}) {
	// noop
}
