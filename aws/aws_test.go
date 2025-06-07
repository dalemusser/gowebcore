package awsutil

import (
	"context"
	"testing"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
)

func fakeRegion(o *awsconfig.LoadOptions) error { o.Region = "us-east-1"; return nil }

func TestConstructors(t *testing.T) {
	ctx := context.Background()
	_, _ = NewS3(ctx, fakeRegion)         // just ensure no panic
	_, _ = NewCloudFront(ctx, fakeRegion) // same
}
