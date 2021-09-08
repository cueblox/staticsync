package main

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cueblox/blox"
	"github.com/cueblox/blox/content"
	"github.com/hashicorp/go-hclog"
	"gocloud.dev/blob"

	// Import the blob packages we want to be able to open.
	_ "gocloud.dev/blob/azureblob"
	_ "gocloud.dev/blob/gcsblob"
	_ "gocloud.dev/blob/s3blob"
)

// Here is a real implementation of Greeter
type StaticSync struct {
	logger hclog.Logger
	cfg    *blox.Config
}

func (g *StaticSync) Process(bloxConfig string) error {
	g.logger.Debug("message from CloudSync.Process")
	cfg, err := blox.NewConfig(content.BaseConfig)
	if err != nil {
		g.logger.Error("loading base config", "error", err.Error())
		return err
	}
	g.cfg = cfg

	err = g.cfg.LoadConfigString(bloxConfig)
	if err != nil {

		g.logger.Error("loading config", "error", err.Error())
		return err
	}
	staticDir, err := g.cfg.GetString("static_dir")
	if err != nil {
		g.logger.Info("no static directory present, skipping image linking")
		return nil
	}
	err = g.syncDirectory(staticDir)

	g.logger.Debug("Return Value", "error", err)
	return err
}

func (g *StaticSync) syncDirectory(staticDir string) error {
	cwd, err := os.Getwd()
	if err != nil {
		g.logger.Error("Working directory error", "error", err.Error())
	}
	g.logger.Debug("Working directory", "dirname", cwd)

	g.logger.Debug("processing", "dir", staticDir)
	fi, err := os.Stat(staticDir)
	if errors.Is(err, os.ErrNotExist) {
		g.logger.Error("no static directory found, skipping")
		return nil
	}
	if !fi.IsDir() {
		return errors.New("given static directory is not a directory")
	}
	err = filepath.Walk(staticDir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			g.logger.Debug("Current Node", "info", info.Name())

			g.logger.Debug("Processing", "path", path)
			if !info.IsDir() {

				g.logger.Debug("Reading File", "path", path)
				buf, err := ioutil.ReadFile(path)
				if err != nil {

					g.logger.Error("reading file", "path", path, "error", err.Error())
					return err
				}

				relpath, err := filepath.Rel(staticDir, path)
				if err != nil {
					return err
				}

				g.logger.Info("Starting Sync", "path", relpath)
				bucketURL := os.Getenv("STATIC_BUCKET")
				if bucketURL == "" {

					g.logger.Error("Bucket Missing", "bucketURL", bucketURL)
					return errors.New("no STATIC_BUCKET environment variable set")
				}

				ctx := context.Background()
				// Open a connection to the bucket.

				g.logger.Debug("Opening Bucket", "bucket", bucketURL)
				b, err := blob.OpenBucket(ctx, bucketURL)
				if err != nil {
					return fmt.Errorf("failed to setup bucket: %s", err)
				}
				defer b.Close()

				w, err := b.NewWriter(ctx, relpath, nil)
				if err != nil {
					return fmt.Errorf("sync failed to obtain writer: %s", err)
				}
				_, err = w.Write(buf)
				if err != nil {
					return fmt.Errorf("sync failed to write to bucket: %s", err)
				}
				if err = w.Close(); err != nil {
					return fmt.Errorf("sync failed to close: %s", err)
				}
			}

			return nil
		})
	return err
}
