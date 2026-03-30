package blob

import (
	"bytes"
	"context"
	"sync/atomic"

	"golang.org/x/sync/errgroup"

	"github.com/dv1x3r/amazing-core/internal/lib/wrap"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Config struct {
	Endpoint        string `json:"endpoint"`
	Region          string `json:"region"`
	Bucket          string `json:"bucket"`
	PathPrefix      string `json:"path_prefix"`
	AccessKeyID     string `json:"access_key_id"`
	SecretAccessKey string `json:"secret_access_key"`
}

type S3SyncResult struct {
	SyncedFiles  int `json:"synced_files"`
	SkippedFiles int `json:"skipped_files"`
}

func newS3Client(cfg S3Config) *s3.Client {
	awsCfg := aws.Config{
		Region: cfg.Region,
		Credentials: credentials.NewStaticCredentialsProvider(
			cfg.AccessKeyID,
			cfg.SecretAccessKey,
			"",
		),
	}

	opts := []func(*s3.Options){}
	if cfg.Endpoint != "" {
		opts = append(opts, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.Endpoint)
			o.UsePathStyle = true
		})
	}

	return s3.NewFromConfig(awsCfg, opts...)
}

func (s *Service) SyncToS3(ctx context.Context, cfg S3Config) (*S3SyncResult, error) {
	const op = "blob.Service.SyncToS3"

	client := newS3Client(cfg)

	// Step 1: List all existing S3 objects with their sizes
	existing := map[string]int64{} // key -> size
	paginator := s3.NewListObjectsV2Paginator(client, &s3.ListObjectsV2Input{
		Bucket: aws.String(cfg.Bucket),
		Prefix: aws.String(cfg.PathPrefix),
	})

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return nil, wrap.IfErr(op, err)
		}
		for _, obj := range page.Contents {
			if obj.Key != nil && obj.Size != nil {
				existing[*obj.Key] = *obj.Size
			}
		}
	}

	// s.logger.Debug(op, "existing", len(existing))

	// Step 2: Query files from database
	const query = "SELECT cdnid, blob FROM asset_file;"
	rows, err := s.store.DB().QueryContext(ctx, query)
	if err != nil {
		return nil, wrap.IfErr(op, err)
	}
	defer rows.Close()

	// Step 3: Upload files concurrently using errgroup
	var (
		syncedFiles  atomic.Int64
		skippedFiles atomic.Int64
	)

	g, gctx := errgroup.WithContext(ctx)
	g.SetLimit(20) // max concurrent uploads

	for rows.Next() {
		// Stop feeding new jobs if a previous upload failed
		if gctx.Err() != nil {
			break
		}

		var cdnid string
		var data []byte

		if err := rows.Scan(&cdnid, &data); err != nil {
			return nil, wrap.IfErr(op, err)
		}

		s3Key := cdnid
		if cfg.PathPrefix != "" {
			s3Key = cfg.PathPrefix + cdnid
		}

		// Check if file exists with same size
		if existingSize, exists := existing[s3Key]; exists && existingSize == int64(len(data)) {
			skipped := skippedFiles.Add(1)
			s.logger.Debug(op, "cdnid", cdnid, "status", "skipped", "synced", syncedFiles.Load(), "skipped", skipped)
			continue
		}

		g.Go(func() error {
			_, err := client.PutObject(gctx, &s3.PutObjectInput{
				Bucket:      aws.String(cfg.Bucket),
				Key:         aws.String(s3Key),
				Body:        bytes.NewReader(data),
				ContentType: aws.String("application/octet-stream"),
			})
			if err != nil {
				return wrap.IfErr(op, err)
			}

			synced := syncedFiles.Add(1)
			s.logger.Debug(op, "cdnid", cdnid, "status", "synced", "synced", synced, "skipped", skippedFiles.Load())
			return nil
		})
	}

	if err := rows.Err(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	if err := g.Wait(); err != nil {
		return nil, wrap.IfErr(op, err)
	}

	result := &S3SyncResult{
		SyncedFiles:  int(syncedFiles.Load()),
		SkippedFiles: int(skippedFiles.Load()),
	}

	s.logger.Info("sync cache files with s3: finished",
		"synced_files", result.SyncedFiles,
		"skipped_files", result.SkippedFiles,
	)

	return result, nil
}
