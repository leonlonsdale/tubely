package util

import (
	"io"
	"log"
)

func SafeClose(c io.Closer) {
	if c == nil {
		return
	}
	if err := c.Close(); err != nil {
		log.Printf("Warning: failed to close resource: %v", err)
	}
}
