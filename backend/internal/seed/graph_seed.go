package seed

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
)

const (
	startPointMethod    = "START_POINT"
	internalMsgMethod   = "INTERNAL_MSG"
	playerReplicaMethod = "PLAYER_REPLICA"
	npcReplicaMethod    = "NPC_REPLICA"
)

type GraphSeederRepository interface {
	Sync(ctx context.Context, nodes []*domain.Node) error
}

type dialogBlock struct {
	InternalMessage string   `json:"internal message,omitempty"`
	PlayerMessages  []string `json:"player messages,omitempty"`
	NPCMessages     []string `json:"npc messages,omitempty"`
}

func CompareAndSeed(ctx context.Context, filePath string, repo GraphSeederRepository) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read graph configuration: %w", err)
	}

	var scenarios [][]dialogBlock
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&scenarios); err != nil {
		return fmt.Errorf("decode graph configuration: %w", err)
	}
	if err := decoder.Decode(new(any)); !errors.Is(err, io.EOF) {
		if err == nil {
			return errors.New("decode graph configuration: multiple JSON values")
		}
		return fmt.Errorf("decode graph configuration: %w", err)
	}

	nodes, err := buildGraph(scenarios)
	if err != nil {
		return err
	}
	if err := repo.Sync(ctx, nodes); err != nil {
		return fmt.Errorf("synchronize graph: %w", err)
	}
	return nil
}

func buildGraph(scenarios [][]dialogBlock) ([]*domain.Node, error) {
	if len(scenarios) == 0 {
		return nil, errors.New("graph configuration must contain at least one scenario")
	}

	var nodes []*domain.Node
	addNode := func(method string, content string) *domain.Node {
		var cont *domain.Content
		if content != "" {
			cont = domain.NewContent(0, content)
		}
		node := domain.NewNode(int64(len(nodes)+1), method, cont)
		nodes = append(nodes, node)
		return node
	}
	connect := func(from, to []*domain.Node) {
		for _, source := range from {
			source.Next = append(source.Next, to...)
		}
	}

	start := addNode(startPointMethod, "")
	for scenarioIndex, scenario := range scenarios {
		if len(scenario) == 0 {
			return nil, fmt.Errorf("scenario %d is empty", scenarioIndex+1)
		}
		first := scenario[0]
		if strings.TrimSpace(first.InternalMessage) == "" ||
			len(first.PlayerMessages) != 0 ||
			len(first.NPCMessages) != 0 {
			return nil, fmt.Errorf("scenario %d must start with an internal message", scenarioIndex+1)
		}

		internal := addNode(internalMsgMethod, first.InternalMessage)
		connect([]*domain.Node{start}, []*domain.Node{internal})
		previous := []*domain.Node{internal}

		for blockIndex, block := range scenario[1:] {
			if strings.TrimSpace(block.InternalMessage) != "" ||
				len(block.PlayerMessages) == 0 ||
				len(block.NPCMessages) == 0 {
				return nil, fmt.Errorf(
					"scenario %d block %d must contain player and NPC messages",
					scenarioIndex+1, blockIndex+2,
				)
			}
			for _, message := range block.PlayerMessages {
				if strings.TrimSpace(message) == "" {
					return nil, fmt.Errorf("scenario %d block %d contains an empty player message", scenarioIndex+1, blockIndex+2)
				}
			}
			for _, message := range block.NPCMessages {
				if strings.TrimSpace(message) == "" {
					return nil, fmt.Errorf("scenario %d block %d contains an empty NPC message", scenarioIndex+1, blockIndex+2)
				}
			}

			players := make([]*domain.Node, 0, len(block.PlayerMessages))
			for _, message := range block.PlayerMessages {
				players = append(players, addNode(playerReplicaMethod, message))
			}
			connect(previous, players)

			npcMessages := make([]*domain.Node, 0, len(block.NPCMessages))
			for _, message := range block.NPCMessages {
				npcMessages = append(npcMessages, addNode(npcReplicaMethod, message))
			}
			connect(players, npcMessages)
			previous = npcMessages
		}
	}

	return nodes, nil
}
