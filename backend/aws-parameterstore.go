package backend

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// AWSParameterStore implements Backend
type AWSParameterStore struct {
	client   *ssm.SSM
	kmsKeyID string
}

func NewAWSParameterStore(ctx context.Context, p *Profile) (*AWSParameterStore, error) {
	sess := session.New(&aws.Config{
		Region:     aws.String(p.Location),
		MaxRetries: aws.Int(2),
	})
	return &AWSParameterStore{client: ssm.New(sess), kmsKeyID: p.CryptoKey}, nil
}

func (s *AWSParameterStore) Get(ctx context.Context, p *Profile, key string) ([]byte, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String(key),
	}
	output, err := s.client.GetParameter(input)
	if err != nil {
		return []byte{}, err
	}
	return []byte(*output.Parameter.Value), nil
}
func (s *AWSParameterStore) List(ctx context.Context, p *Profile) (list []Key, err error) {
	input := &ssm.GetParametersByPathInput{
		Path:       aws.String("/"),
		MaxResults: aws.Int64(10),
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
				Owner:     fmt.Sprintf("type: %v datatype: %v version: %v", *each.Type, *each.DataType, *each.Version),
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
func (s *AWSParameterStore) CheckExists(ctx context.Context, p *Profile, key string) (bool, error) {
	input := &ssm.GetParameterInput{
		Name: aws.String(key),
	}
	_, err := s.client.GetParameter(input)
	if _, ok := err.(*ssm.ParameterNotFound); ok {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil

}
func (s *AWSParameterStore) Put(ctx context.Context, p *Profile, key, value string) error {
	input := &ssm.PutParameterInput{
		Name:        aws.String(key),
		Value:       aws.String(value),
		Overwrite:   aws.Bool(false),
		DataType:    aws.String("text"),
		Description: aws.String(fmt.Sprintf("created by %s using kiya", os.Getenv("USER"))),
		Tags:        []*ssm.Tag{{Key: aws.String("creator"), Value: aws.String(os.Getenv("USER"))}},
		KeyId:       aws.String(s.kmsKeyID),
		Type:        aws.String("SecureString"),
	}
	_, err := s.client.PutParameter(input)
	if err != nil {
		return err
	}
	return nil
}
func (s *AWSParameterStore) Delete(ctx context.Context, p *Profile, key string) error {
	input := &ssm.DeleteParameterInput{
		Name: aws.String(key),
	}
	_, err := s.client.DeleteParameter(input)
	return err
}

func (s *AWSParameterStore) Close() error {
	// noop
	return nil
}
