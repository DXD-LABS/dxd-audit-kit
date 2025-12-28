package verify

import (
	"context"
	"os"
	"testing"
)

func TestVerifyDocument(t *testing.T) {
	// Tạo file tạm để test
	tmpFile, err := os.CreateTemp("", "test-verify-*.txt")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	content := []byte("hello world")
	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("failed to write to temp file: %v", err)
	}
	tmpFile.Close()

	ctx := context.Background()

	t.Run("Success case", func(t *testing.T) {
		result, err := VerifyDocument(ctx, tmpFile.Name())
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		// SHA-256 của "hello world"
		expectedHash := "b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9"
		if result.Hash != expectedHash {
			t.Errorf("expected hash %s, got %s", expectedHash, result.Hash)
		}

		if result.Size != int64(len(content)) {
			t.Errorf("expected size %d, got %d", len(content), result.Size)
		}

		if result.HashAlgo != "sha256" {
			t.Errorf("expected algo sha256, got %s", result.HashAlgo)
		}

		if result.Path != tmpFile.Name() {
			t.Errorf("expected path %s, got %s", tmpFile.Name(), result.Path)
		}
	})

	t.Run("File not found", func(t *testing.T) {
		_, err := VerifyDocument(ctx, "non-existent-file.txt")
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})
}
