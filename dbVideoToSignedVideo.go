package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/bootdotdev/learn-file-storage-s3-golang-starter/internal/database"
)

// func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
// 	log.Printf("+%v", video)
// 	parts := strings.Split(*video.VideoURL, ",")
// 	if len(parts) != 2 {
// 		return video, fmt.Errorf("invalid video url")
// 	}

// 	bucket := parts[0]
// 	key := parts[1]

// 	url, err := cfg.utils.GeneratePresignedURL(cfg.s3Client, bucket, key, time.Hour)
// 	if err != nil {
// 		return video, fmt.Errorf("error creating presigned url: %w", err)
// 	}

// 	video.VideoURL = &url

// 	return video, nil
// }

func (cfg *apiConfig) dbVideoToSignedVideo(video database.Video) (database.Video, error) {
	if video.VideoURL == nil {
		return video, fmt.Errorf("missing video url")
	}

	parts := strings.Split(*video.VideoURL, ",")
	if len(parts) != 2 {
		return video, fmt.Errorf("invalid videl url format for signing")
	}

	bucket := parts[0]
	key := parts[1]

	presignedURL, err := cfg.utils.GeneratePresignedURL(cfg.s3Client, bucket, key, time.Hour)
	if err != nil {
		return video, fmt.Errorf("couldn't generate presigned url")
	}

	// // DEBUG
	// log.Printf("\nPresignedURL: %s\n", presignedURL)

	video.VideoURL = &presignedURL

	return video, nil
}
