package awsutil

import (
	"context"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Client struct {
	Client *s3.Client
}

// NewS3 loads default AWS config.
func NewS3(ctx context.Context, cfgs ...func(*awsconfig.LoadOptions) error) (*S3Client, error) {
	cfg, err := awsconfig.LoadDefaultConfig(ctx, cfgs...)
	if err != nil {
		return nil, err
	}
	return &S3Client{s3.NewFromConfig(cfg)}, nil
}

func (c *S3Client) PresignPUT(ctx context.Context, bucket, key string, d time.Duration) (string, error) {
	ps := s3.NewPresignClient(c.Client)
	out, err := ps.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(d))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

func (c *S3Client) PresignGET(ctx context.Context, bucket, key string, d time.Duration) (string, error) {
	ps := s3.NewPresignClient(c.Client)
	out, err := ps.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}, s3.WithPresignExpires(d))
	if err != nil {
		return "", err
	}
	return out.URL, nil
}

// UploadStream uploads from any io.Reader using the high-level uploader.
func (c *S3Client) UploadStream(ctx context.Context, bucket, key string, body io.Reader) error {
	u := manager.NewUploader(c.Client)
	_, err := u.Upload(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   body,
	})
	return err
}
