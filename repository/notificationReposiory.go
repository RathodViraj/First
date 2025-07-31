package repository

import (
	"First/model"
	"database/sql"
)

type NotificationRepository interface {
	SaveNotification(n model.Notification) error
	GetNotificationsByUser(userID int) ([]model.Notification, error)
}

type notificationRepo struct {
	db *sql.DB
}

func NewNotificationRepo(db *sql.DB) NotificationRepository {
	return &notificationRepo{db}
}

func (r *notificationRepo) SaveNotification(n model.Notification) error {
	query := `INSERT INTO notifications (type, from_user, to_user, post_id, message, is_read, timestamp)
              VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, n.Type, n.FromUser, n.ToUser, n.PostID, n.Message, n.Seen, n.Timestamp)
	return err
}

func (r *notificationRepo) GetNotificationsByUser(userID int) ([]model.Notification, error) {
	query := `SELECT id, type, from_user, to_user, post_id, message, is_read, timestamp 
              FROM notifications 
              WHERE to_user = ? 
              ORDER BY timestamp DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifs []model.Notification
	for rows.Next() {
		var n model.Notification
		var postID sql.NullInt64

		err := rows.Scan(&n.ID, &n.Type, &n.FromUser, &n.ToUser, &postID, &n.Message, &n.Seen, &n.Timestamp)
		if err != nil {
			return nil, err
		}

		if postID.Valid {
			pid := int(postID.Int64)
			n.PostID = &pid
		}

		notifs = append(notifs, n)
	}
	return notifs, nil
}
