package main

import "os"

type BucketConfig struct {
	ConnectionString string
	AccessID string
	AccessKey string
}

type BookBucketConfig struct {
	BucketConfig
}


type SiteBucketConfig struct {
	BucketConfig
}

func LoadBookBucketConfig() BookBucketConfig {
	host := getEnv("BOOK_BUCKET_HOST", "s3://my-books")
	accessID := os.Getenv("ACCESS_ID")
	key := os.Getenv("ACCESS_KEY")
	return BookBucketConfig{
		BucketConfig{
			ConnectionString: host,
			AccessID:         accessID,
			AccessKey:        key,
		},
	}
}

func LoadSiteBucketConfig() SiteBucketConfig {
	host := getEnv("SITE_BUCKET_HOST", "s3://my-books")
	accessID := os.Getenv("ACCESS_ID")
	key := os.Getenv("ACCESS_KEY")
	return SiteBucketConfig{
		BucketConfig{
			ConnectionString: host,
			AccessID:         accessID,
			AccessKey:        key,
		},
	}
}

func getEnv(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
