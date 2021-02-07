package s3

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	ourfmt "github.com/gomicro/flow/fmt"
	"github.com/gosuri/uiprogress"
	"github.com/spf13/cobra"
)

var (
	bucket string
	bar    *uiprogress.Bar
)

func init() {
	S3Cmd.AddCommand(SyncCmd)

	SyncCmd.Flags().StringVarP(&bucket, "bucket", "b", "", "bucket to sync to")
}

// SyncCmd represents the S3 sync commands
var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync a directory to S3",
	Long:  "Sync a directory to S3",
	Run:   syncFunc,
	Args:  cobra.MaximumNArgs(1),
}

func syncFunc(cmd *cobra.Command, args []string) {
	path := args[0]
	if path == "" {
		path = "."
	}

	uiprogress.Start()

	iter, err := newSyncFolderIterator(path, bucket)
	if err != nil {
		ourfmt.Printf("unexpected error occurred creating folder iterator: %v", err)
		os.Exit(1)
	}

	bar = uiprogress.AddBar(len(iter.fileInfos))
	bar.AppendCompleted()
	bar.PrependElapsed()

	err = s3Uploader.UploadWithIterator(aws.BackgroundContext(), iter)
	if err != nil {
		ourfmt.Printf("unexpected error has occurred: %v", err)
		os.Exit(1)
	}

	err = iter.Err()
	if err != nil {
		ourfmt.Printf("unexpected error occurred during file walking: %v", err)
		os.Exit(1)
	}

	uiprogress.Stop()
	ourfmt.Printf("Upload complete...\n")
}

type syncFolderIterator struct {
	bucket    string
	fileInfos []fileInfo
	err       error
}

type fileInfo struct {
	key      string
	fullpath string
}

func newSyncFolderIterator(path, bucket string) (*syncFolderIterator, error) {
	metadata := []fileInfo{}
	err := filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			key := strings.TrimPrefix(p, path)
			metadata = append(metadata, fileInfo{key, p})
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("filepath walk: %w", err)
	}

	return &syncFolderIterator{
		bucket:    bucket,
		fileInfos: metadata,
		err:       nil,
	}, nil
}

func (iter *syncFolderIterator) Next() bool {
	bar.Incr()
	return len(iter.fileInfos) > 0
}

func (iter *syncFolderIterator) Err() error {
	return iter.err
}

func (iter *syncFolderIterator) UploadObject() s3manager.BatchUploadObject {
	fi := iter.fileInfos[0]
	iter.fileInfos = iter.fileInfos[1:]

	body, err := os.Open(fi.fullpath)
	if err != nil {
		iter.err = err
	}

	extension := filepath.Ext(fi.key)
	mimeType := mime.TypeByExtension(extension)

	if mimeType == "" {
		mimeType = "binary/octet-stream"
	}

	input := s3manager.UploadInput{
		Bucket:      &iter.bucket,
		Key:         &fi.key,
		Body:        body,
		ContentType: &mimeType,
	}

	return s3manager.BatchUploadObject{
		Object: &input,
	}
}
