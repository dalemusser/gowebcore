package awsutil

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront"
	"github.com/aws/aws-sdk-go-v2/service/cloudfront/types"
)

type CFClient struct {
	Client *cloudfront.Client
}

func NewCloudFront(ctx context.Context, cfgs ...func(*awsconfig.LoadOptions) error) (*CFClient, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, cfgs...)
	if err != nil {
		return nil, err
	}
	return &CFClient{cloudfront.NewFromConfig(cfg)}, nil
}

// Invalidate submits an invalidation request for the given paths.
func (c *CFClient) Invalidate(ctx context.Context, distID string, paths ...string) (string, error) {
	if len(paths) == 0 {
		paths = []string{"/*"}
	}
	qty := int32(len(paths))

	out, err := c.Client.CreateInvalidation(ctx, &cloudfront.CreateInvalidationInput{
		DistributionId: aws.String(distID),
		InvalidationBatch: &types.InvalidationBatch{
			CallerReference: aws.String("gowebcore-" + time.Now().UTC().Format("20060102T150405")),
			Paths: &types.Paths{
				Items:    paths,
				Quantity: &qty,
			},
		},
	})
	if err != nil {
		return "", err
	}
	return *out.Invalidation.Id, nil
}
