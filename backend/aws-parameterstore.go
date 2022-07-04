package backend

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// AWSParameterStore implements Backend for AWS Parameter Store service.
type AWSParameterStore struct {
	client   *ssm.SSM
	kmsKeyID string
}

// NewAWSParameterStore returns a new AWSParameterStore with an initialized AWS SSM client.
func NewAWSParameterStore(ctx context.Context, p *Profile) (*AWSParameterStore, error) {
	sess := session.New(&aws.Config{
		Region:     aws.String(p.Location),
		MaxRetries: aws.Int(2),
	})
	return &AWSParameterStore{client: ssm.New(sess), kmsKeyID: p.CryptoKey}, nil
}

// Get returns the decrypted value for a parameter by key.
func (s *AWSParameterStore) Get(ctx context.Context, p *Profile, key string) ([]byte, error) {
	input := &ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(true),
	}
	output, err := s.client.GetParameter(input)
	if err != nil {
		return []byte{}, err
	}
	return []byte(*output.Parameter.Value), nil
}

// List returns all keys available.
func (s *AWSParameterStore) List(ctx context.Context, p *Profile) (list []Key, err error) {
	input := &ssm.GetParametersByPathInput{
		Path:       aws.String("/"),
		MaxResults: aws.Int64(10), // is the documented maximum
		Recursive:  aws.Bool(true),
	}
	for {
		output, err := s.client.GetParametersByPath(input)
		if err != nil {
			return []Key{}, err
		}
		for _, each := range output.Parameters {
			list = append(list, Key{
				Name:      *each.Name,
				CreatedAt: *each.LastModifiedDate,
				Info:      fmt.Sprintf("type: %v datatype: %v version: %v", *each.Type, *each.DataType, *each.Version),
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
	_, err := s.client.GetParameter(input)
	if _, ok := err.(*ssm.ParameterNotFound); ok {
		return false, nil
	}
	// other error?
	if err != nil {
		return false, err
	}
	return true, nil

}

// Put write the parameter and its value using encryption ;either the default key or the one specified in the profile.
func (s *AWSParameterStore) Put(ctx context.Context, p *Profile, key, value string) error {
	input := &ssm.PutParameterInput{
		Name:        aws.String(key),
		Value:       aws.String(value),
		Overwrite:   aws.Bool(false),
		DataType:    aws.String("text"),
		Description: aws.String(fmt.Sprintf("created by %s using kiya", os.Getenv("USER"))),
		Tags:        []*ssm.Tag{{Key: aws.String("creator"), Value: aws.String(os.Getenv("USER"))}},
		Type:        aws.String("SecureString"),
	}
	// only if CryptoKey is set in the Profile then we set the KeyId
	// which overrides the default key associated with the AWS account
	if p.CryptoKey != "" {
		input.KeyId = aws.String(s.kmsKeyID)
	}
	_, err := s.client.PutParameter(input)
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
	_, err := s.client.DeleteParameter(input)
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
