package event

import (
	"context"
	"errors"

	"github.com/dipithedipi/password-manager/prisma/db"
)

func NewEvent(clientPostgresDb *db.PrismaClient, eventType string, title string, description string, ip string, userId string) error {
	if title == "" || ip == "" || userId == ""{
		return errors.New("title, IP, and User ID are required for settings up a new event")
	}
	
	// check db connection
	if clientPostgresDb == nil {
		return errors.New("database connection is required")
	}

	clientPostgresDb.Event.CreateOne(
		db.Event.Title.Set(title),
		db.Event.IPAddress.Set(ip),
		db.Event.Type.Set(eventType),
		db.Event.User.Link(
			db.User.ID.Equals(userId),
		),
		db.Event.Description.SetOptional(&description),
	).Exec(context.Background())

	return nil
}