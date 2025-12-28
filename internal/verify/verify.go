package verify

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/google/uuid"
)

// VerifyResult chứa kết quả verify tài liệu
type VerifyResult struct {
	DocumentID uuid.UUID `json:"document_id"`
	Path       string    `json:"path"`
	Hash       string    `json:"hash"`
	HashAlgo   string    `json:"hash_algo"`
	Size       int64     `json:"size"`
	VerifiedAt time.Time `json:"verified_at"`
}

// VerifyDocument đọc file, tính hash và trả về VerifyResult
func VerifyDocument(ctx context.Context, path string) (VerifyResult, error) {
	file, err := os.Open(path)
	if err != nil {
		return VerifyResult{}, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return VerifyResult{}, fmt.Errorf("failed to get file info: %w", err)
	}

	if stat.IsDir() {
		return VerifyResult{}, fmt.Errorf("path is a directory: %s", path)
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return VerifyResult{}, fmt.Errorf("failed to compute hash: %w", err)
	}

	return VerifyResult{
		DocumentID: uuid.New(),
		Path:       path,
		Hash:       fmt.Sprintf("%x", hash.Sum(nil)),
		HashAlgo:   "sha256",
		Size:       stat.Size(),
		VerifiedAt: time.Now(),
	}, nil
}
