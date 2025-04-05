package notification

import (
	"context"
	"encoding/json"
	"time"

	model "github.com/ShmaykhelDuo/battler/internal/model/notification"
	"github.com/ShmaykhelDuo/battler/internal/pkg/db/postgres"
	"github.com/google/uuid"
)

type Notification struct {
	ID         uuid.UUID       `db:"id"`
	Type       int             `db:"type_id"`
	Payload    json.RawMessage `db:"payload"`
	CreateTime time.Time       `db:"created_at"`
}

type PostgresRepository struct {
	db *postgres.DB
}

func NewPostgresRepository(db *postgres.DB) *PostgresRepository {
	return &PostgresRepository{db: db}
}

func (r *PostgresRepository) CreateNotification(ctx context.Context, userID uuid.UUID, n model.Notification) error {
	sql := "INSERT INTO notifications (id, user_id, type_id, payload, received, created_at) VALUES ($1, $2, $3, $4, false, $5)"

	_, err := r.db.Exec(ctx, sql, n.ID, userID, int(n.Type), n.Payload, n.CreateTime)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostgresRepository) PendingNotifications(ctx context.Context, userID uuid.UUID) ([]model.Notification, error) {
	sql := "SELECT id, type_id, payload, created_at FROM notifications WHERE user_id = $1 AND received = false AND created_at < now() ORDER BY created_at;"

	var dto []Notification
	err := r.db.Select(ctx, &dto, sql, userID)
	if err != nil {
		return nil, err
	}

	res := make([]model.Notification, len(dto))
	for i, n := range dto {
		res[i] = model.Notification{
			ID:         n.ID,
			Type:       model.Type(n.Type),
			Payload:    n.Payload,
			CreateTime: n.CreateTime,
		}
	}

	return res, nil
}

func (r *PostgresRepository) MarkNotificationReceived(ctx context.Context, id uuid.UUID) error {
	sql := "UPDATE notifications SET received = true WHERE id = $1;"

	_, err := r.db.Exec(ctx, sql, id)
	if err != nil {
		return err
	}

	return nil
}
