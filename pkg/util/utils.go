package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"os/exec"
)

// type UtilService interface {
// 	SafeClose(io.Closer)
// 	GetVideoAspectRatio(string) (string, error)
// }

type Utils struct{}

func NewUtils() *Utils {
	util := Utils{}
	return &util
}

func (u *Utils) SafeClose(c io.Closer) {
	if c == nil {
		return
	}
	if err := c.Close(); err != nil {
		log.Printf("Warning: failed to close resource: %v", err)
	}
}

func (u *Utils) GetVideoAspectRatio(filepath string) (string, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-print_format", "json", "-show_streams", filepath)

	buf := &bytes.Buffer{}
	cmd.Stdout = buf

	if err := cmd.Run(); err != nil {
		return "", err
	}

	type Stream struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	}

	type Data struct {
		Streams []Stream `json:"streams"`
	}

	var data Data

	if err := json.Unmarshal(buf.Bytes(), &data); err != nil {
		return "", err
	}

	var width, height int
	found := false
	for _, s := range data.Streams {
		if s.Width > 0 && s.Height > 0 {
			width = s.Width
			height = s.Height
			found = true
			break
		}
	}
	if !found {
		return "", fmt.Errorf("no video stream found")
	}

	ratio := aspectRatioType(width, height)

	return ratio, nil
}

func (u *Utils) ProcessVideoForFastStart(filepath string) (string, error) {

	outputFilepath := filepath + ".processing"

	cmd := exec.Command("ffmpeg", "-i", filepath, "-c", "copy", "-movflags", "faststart", "-f", "mp4", outputFilepath)

	if err := cmd.Run(); err != nil {
		return "", nil
	}

	return outputFilepath, nil
}

func aspectRatioType(width, height int) string {
	if width == 0 || height == 0 {
		return "other"
	}

	ratio := float64(width) / float64(height)
	const tolerance = 0.02 // 2% tolerance

	if math.Abs(ratio-16.0/9.0) < tolerance {
		return "16:9"
	}
	if math.Abs(ratio-9.0/16.0) < tolerance {
		return "9:16"
	}
	return "other"
}
