package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/dxdlabs/dxd-audit-kit/internal/audit"
	"github.com/dxdlabs/dxd-audit-kit/internal/config"
	"github.com/dxdlabs/dxd-audit-kit/internal/db"
	"github.com/dxdlabs/dxd-audit-kit/internal/verify"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var (
	cfg  config.Config
	repo audit.Repository
)

func main() {
	cfg = config.Load()

	rootCmd := &cobra.Command{
		Use:   "dxd-audit-cli",
		Short: "DXD Audit CLI is a toolkit for audit logs and verification",
	}

	rootCmd.AddCommand(verifyCmd())
	rootCmd.AddCommand(logEventCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initRepo() {
	database, err := db.Open(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	repo = audit.NewRepository(database)
}

func verifyCmd() *cobra.Command {
	var filePath string
	var algo string

	cmd := &cobra.Command{
		Use:   "verify",
		Short: "Verify a document and create an audit record if it doesn't exist",
		Run: func(cmd *cobra.Command, args []string) {
			if filePath == "" {
				log.Fatal("--file flag is required")
			}

			res, err := verify.VerifyDocument(context.Background(), filePath)
			if err != nil {
				log.Fatalf("Verification failed: %v", err)
			}

			initRepo()

			// Check if exists by hash
			doc, err := repo.GetDocumentByHash(context.Background(), res.Hash)
			if err != nil {
				// Create new if not found
				doc = audit.Document{
					ID:       res.DocumentID,
					Hash:     res.Hash,
					HashAlgo: res.HashAlgo,
					Size:     res.Size,
				}
				doc, err = repo.CreateDocument(context.Background(), doc)
				if err != nil {
					log.Fatalf("Failed to create document in DB: %v", err)
				}
				fmt.Println("New document registered.")
			} else {
				fmt.Println("Document already exists in DB.")
			}

			fmt.Printf("Document ID: %s\n", doc.ID)
			fmt.Printf("Hash:        %s\n", doc.Hash)
			fmt.Printf("Size:        %d bytes\n", doc.Size)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file to verify")
	cmd.Flags().StringVarP(&algo, "hash-algo", "a", "sha256", "Hash algorithm to use")
	return cmd
}

func logEventCmd() *cobra.Command {
	var docIDStr string
	var filePath string
	var email string
	var ip string
	var ua string

	cmd := &cobra.Command{
		Use:   "log-event",
		Short: "Log a signature event for a document",
		Run: func(cmd *cobra.Command, args []string) {
			initRepo()
			ctx := context.Background()
			var docID uuid.UUID

			if docIDStr != "" {
				var err error
				docID, err = uuid.Parse(docIDStr)
				if err != nil {
					log.Fatalf("Invalid document-id: %v", err)
				}
			} else if filePath != "" {
				res, err := verify.VerifyDocument(ctx, filePath)
				if err != nil {
					log.Fatalf("Failed to verify file: %v", err)
				}
				doc, err := repo.GetDocumentByHash(ctx, res.Hash)
				if err != nil {
					log.Fatalf("Document not found in DB. Please run 'verify' first.")
				}
				docID = doc.ID
			} else {
				log.Fatal("Either --document-id or --file is required")
			}

			if email == "" {
				log.Fatal("--signer-email is required")
			}

			event := audit.SignEvent{
				DocumentID:  docID,
				SignerEmail: email,
				IPAddress:   ip,
				UserAgent:   ua,
			}

			event, err := repo.LogSignEvent(ctx, event)
			if err != nil {
				log.Fatalf("Failed to log event: %v", err)
			}

			fmt.Printf("Event logged successfully.\n")
			fmt.Printf("Event ID:  %s\n", event.ID)
			fmt.Printf("Signed At: %s\n", event.SignedAt)
		},
	}

	cmd.Flags().StringVarP(&docIDStr, "document-id", "d", "", "ID of the document")
	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to the file to find document ID")
	cmd.Flags().StringVarP(&email, "signer-email", "e", "", "Email of the signer")
	cmd.Flags().StringVarP(&ip, "ip", "i", "", "IP address of the signer")
	cmd.Flags().StringVarP(&ua, "user-agent", "u", "", "User agent of the signer")

	return cmd
}
