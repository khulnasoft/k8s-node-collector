package collector

import (
	"bytes"
	"compress/bzip2"
	"encoding/base64"
	"io"
	"strings"
)

// decompressBzip2 accept bzip2 compressed bytes and decompress it
// #nosec
func decompressBzip2(compressedBytes []byte) (io.Reader, error) {
	bz2Reader := bzip2.NewReader(bytes.NewReader(compressedBytes))
	uncompressedWriter := new(bytes.Buffer)
	//nolint:gosec
	_, err := io.Copy(uncompressedWriter, bz2Reader)
	if err != nil {
		return nil, err
	}
	return uncompressedWriter, nil
}

// base64Decode accept encoded reader and base64 decode it
func base64Decode(encodedReader io.Reader) ([]byte, error) {
	encodedBytes, err := io.ReadAll(encodedReader)
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(string(bytes.TrimSpace(encodedBytes)))
}

func uncompressAndDecode(kubeletConfig string) ([]byte, error) {
	encodedReader := io.NopCloser(strings.NewReader(kubeletConfig))
	// base64 decode logs
	compressedLogsBytes, err := base64Decode(encodedReader)
	if err != nil {
		return nil, err
	}
	// bzip2 decompress logs
	unCompressedLogsReader, err := decompressBzip2(compressedLogsBytes)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(unCompressedLogsReader)
}
