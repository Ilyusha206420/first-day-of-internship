package seed

import (
	"context"
	"encoding/json"
	"os"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
)

type GraphSeederNodeRepository interface {
}

type GraphSeederEdgeRepository interface {
}

type GraphSeederContentRepository interface {
	FindByData(data string) (*domain.Content, error)
}

type dialogBlock struct {
	InternalMessage string   `json:"internal message,omitempty"`
	PlayerMessages  []string `json:"player messages,omitempty"`
	NPCMessages     []string `json:"npc messages,omitempty"`
}

func SeedCompareWithJson(
	ctx context.Context, filePath string,
	nodeRepo GraphSeederNodeRepository,
	edgeRepo GraphSeederEdgeRepository,
	contentRepo GraphSeederContentRepository,
) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	var scenarios [][]dialogBlock
	decoder := json.NewDecoder(file)

	if err := decoder.Decode(&scenarios); err != nil {
		return err
	}

	for _, scenario := range scenarios {
		for _, block := range scenario {
			if block.InternalMessage != "" {
				c, err := contentRepo.FindByData(block.InternalMessage)
				if err != nil {
					return err
				}
			} else {

			}
		}
	}

	return nil
}
