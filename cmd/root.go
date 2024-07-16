package cmd

import (
	"context"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/transparency-dev/merkle/proof"
	"github.com/transparency-dev/merkle/rfc6962"
)

var period time.Duration
var api string

var rootCmd = &cobra.Command{
	Use:   "eprescription-observer",
	Short: "EPrescription Observer allows to witness and validate consistency the ePrescription events.",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		previousProof, err := FetchConsistencyProof(ctx, api, 0)
		if err != nil {
			slog.ErrorContext(ctx, "Failed to fetch initial consistency proof", "error", err)
			return err
		}

		slog.Info("Initial consistency proof fetched", "size", previousProof.CurrentRoot.Size, "root_hash", previousProof.CurrentRoot.RootHash)

		for {
			newProof, err := FetchConsistencyProof(ctx, api, previousProof.CurrentRoot.Size)
			if err != nil {
				slog.Error("Failed to fetch consistency proof", "error", err, "previous_size", previousProof.CurrentRoot.Size)
				continue
			}

			previousHashDecoded, err := hex.DecodeString(previousProof.CurrentRoot.RootHash)
			if err != nil {
				slog.Error("Failed to decode last root hash", "error", err, "previous_size", previousProof.CurrentRoot.Size)
				return err
			}

			newHashDecoded, err := hex.DecodeString(newProof.CurrentRoot.RootHash)
			if err != nil {
				slog.Error("Failed to decode new root hash", "error", err, "previous_size", previousProof.CurrentRoot.Size, "current_size", newProof.CurrentRoot.Size)
				continue
			}

			decodedProofs, err := decodeProofs(newProof.Proof)
			if err != nil {
				slog.Error("Failed to decode proof hash", "error", err, "previous_size", previousProof.CurrentRoot.Size, "current_size", newProof.CurrentRoot.Size)
				continue
			}

			err = proof.VerifyConsistency(
				rfc6962.DefaultHasher,
				uint64(previousProof.CurrentRoot.Size),
				uint64(newProof.CurrentRoot.Size),
				decodedProofs,
				previousHashDecoded,
				newHashDecoded,
			)
			if err != nil {
				slog.Error("Failed to verify consistency proof", "error", err, "previous_size", previousProof.CurrentRoot.Size, "current_size", newProof.CurrentRoot.Size)
				continue
			}

			slog.Info("Successfully verified consistency proof", "size", previousProof.CurrentRoot.Size, "root_hash", previousProof.CurrentRoot.RootHash, "new_size", newProof.CurrentRoot.Size, "new_root_hash", newProof.CurrentRoot.RootHash)

			previousProof = newProof
			slog.Info("Waiting for next consistency proof", "period", period)
			time.Sleep(period)
		}
	},
}

func decodeProofs(proofs []string) ([][]byte, error) {
	decodedProofs := make([][]byte, 0, len(proofs))
	for _, hexProof := range proofs {
		decodedProof, err := hex.DecodeString(hexProof)
		if err != nil {
			return nil, err
		}
		decodedProofs = append(decodedProofs, decodedProof)
	}
	return decodedProofs, nil
}

type Root struct {
	RootHash string `json:"root_hash"`
	Size     int    `json:"size"`
}

type ConsistencyProof struct {
	Proof       []string `json:"proof"`
	CurrentRoot Root     `json:"current_root"`
}

func FetchConsistencyProof(ctx context.Context, api string, previousSize int) (ConsistencyProof, error) {
	url := strings.TrimSuffix(api, "/") + "/eprescription/consistency-proof"
	if previousSize != 0 {
		url += "?last_version=" + fmt.Sprintf("%d", previousSize)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return ConsistencyProof{}, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return ConsistencyProof{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return ConsistencyProof{}, errors.New("failed to fetch consistency proof")
	}

	proof := ConsistencyProof{}
	err = json.NewDecoder(resp.Body).Decode(&proof)
	if err != nil {
		return ConsistencyProof{}, err
	}

	return proof, nil
}

func Execute() {
	rootCmd.Flags().DurationVar(&period, "period", time.Hour, "Time to wait between consistency proofs checks(1h by default)")
	rootCmd.Flags().StringVar(&api, "api", "https://api.certifaction.io/", "API URL to fetch consistency proofs")

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
