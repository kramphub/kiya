package backend

import (
	"context"

	"github.com/Azure/azure-sdk-for-go/sdk/keyvault/azsecrets"
)

const latestKeyVersion = ""

type AKV struct {
	client *azsecrets.Client
}

func NewAKV(client *azsecrets.Client) *AKV {
	return &AKV{client}
}

func (b *AKV) Get(ctx context.Context, _ *Profile, key string) ([]byte, error) {
	resp, err := b.client.GetSecret(ctx, key, latestKeyVersion, nil)
	if err != nil {
		return nil, err
	}
	return []byte(*resp.Value), nil
}

func (b *AKV) List(ctx context.Context, _ *Profile) ([]Key, error) {
	pager := b.client.NewListSecretsPager(nil)

	var keys []Key
	for pager.More() {
		page, err := pager.NextPage(ctx)
		if err != nil {
			return nil, err
		}

		for _, v := range page.Value {
			keys = append(keys, Key{
				Name:      v.ID.Name(),
				CreatedAt: *v.Attributes.Created,
				Info:      "creator: <Unknown>", // no owner
				Owner:     "<Unknown>",
			})
		}
	}
	return keys, nil
}

func (b *AKV) CheckExists(ctx context.Context, _ *Profile, key string) (bool, error) {
	_, err := b.client.GetSecret(ctx, key, latestKeyVersion, nil)
	return err == nil, err
}

func (b *AKV) Put(ctx context.Context, _ *Profile, key, value string, overwrite bool) error {
	_, err := b.client.SetSecret(ctx, key, azsecrets.SetSecretParameters{Value: &value}, nil)
	if err != nil {
		return err
	}
	return nil
}

func (b *AKV) Delete(ctx context.Context, _ *Profile, key string) error {
	_, err := b.client.DeleteSecret(ctx, key, nil)
	if err != nil {
		return err
	}
	return nil
}

func (b *AKV) SetParameter(key string, value interface{}) {
	//no-op
}

func (b *AKV) Close() error {
	return nil
}
