package pgsql_repo

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Ilyusha206420/first-day-of-internship/backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type GraphRepo struct {
	db *sqlx.DB
}

func NewGraphRepo(db *sqlx.DB) *GraphRepo {
	return &GraphRepo{db: db}
}

type graphNodeSignature struct {
	method     string
	content    string
	hasContent bool
}

type graphStoredNode struct {
	ID          int64          `db:"id"`
	Method      sql.NullString `db:"method_name"`
	ContentData sql.NullString `db:"content_data"`
}

type graphStoredEdge struct {
	SourceID int64 `db:"src_id"`
	TargetID int64 `db:"tgt_id"`
}

type graphEdgeKey struct {
	sourceID int64
	targetID int64
}

func (r *GraphRepo) Sync(ctx context.Context, graphNodes []*domain.Node) error {
	configNodes := make(map[int64]*domain.Node, len(graphNodes))
	for _, node := range graphNodes {
		if node == nil {
			return fmt.Errorf("graph contains a nil node")
		}
		if node.ID <= 0 {
			return fmt.Errorf("graph node has invalid configuration id %d", node.ID)
		}
		if node.Method == "" {
			return fmt.Errorf("graph node %d has no method", node.ID)
		}
		if _, exists := configNodes[node.ID]; exists {
			return fmt.Errorf("graph contains duplicate configuration id %d", node.ID)
		}
		configNodes[node.ID] = node
	}
	if len(configNodes) == 0 {
		return fmt.Errorf("graph contains no nodes")
	}

	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin graph synchronization: %w", err)
	}
	defer tx.Rollback()

	signatures := make(map[int64]graphNodeSignature, len(configNodes))
	contentIDs := make(map[int64]sql.NullInt64, len(configNodes))
	methodIDs := make(map[int64]int64, len(configNodes))
	for _, node := range graphNodes {
		configID := node.ID
		methodID, err := ensureGraphMethod(ctx, tx, node.Method)
		if err != nil {
			return fmt.Errorf("ensure graph method %q: %w", node.Method, err)
		}
		methodIDs[configID] = methodID

		signature := graphNodeSignature{method: node.Method}
		if node.Cont != nil {
			contentID, err := ensureGraphContent(ctx, tx, node.Cont.Data)
			if err != nil {
				return fmt.Errorf("ensure content for graph node %d: %w", configID, err)
			}
			signature.content = node.Cont.Data
			signature.hasContent = true
			contentIDs[configID] = sql.NullInt64{Int64: contentID, Valid: true}
		}
		signatures[configID] = signature
	}

	var existingNodes []graphStoredNode
	query := `
		SELECT n.id, m.method_name, c.data AS content_data
		FROM nodes n
		LEFT JOIN methods m ON m.id = n.method_id
		LEFT JOIN contents c ON c.id = n.content_id
		ORDER BY n.id`
	if err := tx.SelectContext(ctx, &existingNodes, query); err != nil {
		return fmt.Errorf("load existing graph nodes: %w", err)
	}

	existingBySignature := make(map[graphNodeSignature][]int64, len(existingNodes))
	for _, node := range existingNodes {
		if !node.Method.Valid {
			continue
		}
		signature := graphNodeSignature{method: node.Method.String}
		if node.ContentData.Valid {
			signature.content = node.ContentData.String
			signature.hasContent = true
		}
		existingBySignature[signature] = append(existingBySignature[signature], node.ID)
	}

	databaseIDs := make(map[int64]int64, len(configNodes))
	usedNodeIDs := make(map[int64]struct{}, len(configNodes))
	for _, node := range graphNodes {
		configID := node.ID
		signature := signatures[configID]
		if candidates := existingBySignature[signature]; len(candidates) > 0 {
			databaseIDs[configID] = candidates[0]
			usedNodeIDs[candidates[0]] = struct{}{}
			existingBySignature[signature] = candidates[1:]
			continue
		}

		var contentID any
		if id := contentIDs[configID]; id.Valid {
			contentID = id.Int64
		}
		var databaseID int64
		query := `INSERT INTO nodes (method_id, content_id) VALUES ($1, $2) RETURNING id`
		if err := tx.QueryRowxContext(ctx, query, methodIDs[configID], contentID).Scan(&databaseID); err != nil {
			return fmt.Errorf("insert graph node %d: %w", configID, err)
		}
		databaseIDs[configID] = databaseID
		usedNodeIDs[databaseID] = struct{}{}
	}

	desiredEdges := make(map[graphEdgeKey]struct{})
	for _, node := range graphNodes {
		configID := node.ID
		for _, next := range node.Next {
			if next == nil {
				return fmt.Errorf("graph node %d links to a nil node", configID)
			}
			targetID, exists := databaseIDs[next.ID]
			if !exists {
				return fmt.Errorf("graph node %d links to unknown configuration node %d", configID, next.ID)
			}
			desiredEdges[graphEdgeKey{
				sourceID: databaseIDs[configID],
				targetID: targetID,
			}] = struct{}{}
		}
	}

	var existingEdges []graphStoredEdge
	if err := tx.SelectContext(ctx, &existingEdges, `SELECT src_id, tgt_id FROM edges`); err != nil {
		return fmt.Errorf("load existing graph edges: %w", err)
	}
	existingEdgeSet := make(map[graphEdgeKey]struct{}, len(existingEdges))
	for _, edge := range existingEdges {
		key := graphEdgeKey{sourceID: edge.SourceID, targetID: edge.TargetID}
		existingEdgeSet[key] = struct{}{}
		if _, wanted := desiredEdges[key]; !wanted {
			query := `DELETE FROM edges WHERE src_id = $1 AND tgt_id = $2`
			if _, err := tx.ExecContext(ctx, query, edge.SourceID, edge.TargetID); err != nil {
				return fmt.Errorf("delete obsolete graph edge %d -> %d: %w", edge.SourceID, edge.TargetID, err)
			}
		}
	}
	for edge := range desiredEdges {
		if _, exists := existingEdgeSet[edge]; exists {
			continue
		}
		query := `INSERT INTO edges (src_id, tgt_id) VALUES ($1, $2)`
		if _, err := tx.ExecContext(ctx, query, edge.sourceID, edge.targetID); err != nil {
			return fmt.Errorf("insert graph edge %d -> %d: %w", edge.sourceID, edge.targetID, err)
		}
	}

	for _, node := range existingNodes {
		if _, used := usedNodeIDs[node.ID]; used {
			continue
		}
		if _, err := tx.ExecContext(ctx, `DELETE FROM nodes WHERE id = $1`, node.ID); err != nil {
			return fmt.Errorf("delete obsolete graph node %d: %w", node.ID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit graph synchronization: %w", err)
	}
	return nil
}

func ensureGraphMethod(ctx context.Context, tx *sqlx.Tx, method string) (int64, error) {
	query := `INSERT INTO methods (method_name) VALUES ($1) ON CONFLICT (method_name) DO NOTHING`
	if _, err := tx.ExecContext(ctx, query, method); err != nil {
		return 0, err
	}

	var id int64
	query = `SELECT id FROM methods WHERE method_name = $1`
	if err := tx.GetContext(ctx, &id, query, method); err != nil {
		return 0, err
	}
	return id, nil
}

func ensureGraphContent(ctx context.Context, tx *sqlx.Tx, data string) (int64, error) {
	query := `INSERT INTO contents (data) VALUES ($1) ON CONFLICT (data) DO NOTHING`
	if _, err := tx.ExecContext(ctx, query, data); err != nil {
		return 0, err
	}

	var id int64
	query = `SELECT id FROM contents WHERE data = $1`
	if err := tx.GetContext(ctx, &id, query, data); err != nil {
		return 0, err
	}
	return id, nil
}
