package helper

import (
	"encoding/base64"
)

func ImageGet(url string, imageBase64 string) []byte {
	if url != "" {
		data, err := URLDownload(url)
		if err == nil {
			return data
		}
	}

	if imageBase64 != "" {
		data, err := base64.StdEncoding.DecodeString(imageBase64)
		if err == nil {
			return data
		}
	}

	return nil
}
