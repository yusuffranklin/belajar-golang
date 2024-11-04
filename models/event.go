package models

import (
	"database/sql"
	"fmt"
	"rest-api-practice/database"
	"time"
)

type Event struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description" binding:"required"`
	Location    string    `json:"location" binding:"required"`
	DateTime    time.Time `json:"dateTime" binding:"required"`
	UserID      int64     `json:"userId"`
}

var events = []Event{}

func (e *Event) Save() error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id)
	VALUES (?, ?, ?, ?, ?)`

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(e.Name, e.Description, e.Location, e.DateTime, e.UserID)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	e.ID = id
	return nil
}

func GetAllEvents() ([]Event, error) {
	query := "SELECT * FROM events"

	rows, err := database.Db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event
		var dateTimeRaw sql.NullString

		err := rows.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &dateTimeRaw, &event.UserID)
		if err != nil {
			return nil, err
		}

		if dateTimeRaw.Valid {
			event.DateTime, err = time.Parse("2006-01-02 15:04:05", dateTimeRaw.String)
			if err != nil {
				return nil, fmt.Errorf("failed to parse DateTime: %w", err)
			}
		}

		events = append(events, event)
	}

	// Check for any error that may have occurred during iteration
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration failed: %w", err)
	}

	return events, nil
}

func GetEventByID(id int64) (*Event, error) {
	query := "SELECT id, name, description, location, dateTime, user_id FROM events WHERE id = ?"
	row := database.Db.QueryRow(query, id)

	var event Event
	var dateTimeRaw sql.NullString

	err := row.Scan(&event.ID, &event.Name, &event.Description, &event.Location, &dateTimeRaw, &event.UserID)
	if err != nil {
		return nil, err
	}

	if dateTimeRaw.Valid {
		event.DateTime, err = time.Parse("2006-01-02 15:04:05", dateTimeRaw.String)
		if err != nil {
			return nil, fmt.Errorf("failed to parse DateTime: %w", err)
		}
	}

	return &event, nil
}

func (event Event) Update() error {
	query := `
	UPDATE events
	SET name = ?, description = ?, location = ?, dateTime = ?
	WHERE id = ?
	`

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.Name, event.Description, event.Location, event.DateTime, event.ID)
	return err
}

func (event Event) Delete() error {
	query := "DELETE FROM events WHERE id = ?"

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(event.ID)
	return err
}

func (event Event) Register(userId int64) error {
	query := "INSERT INTO registrations(event_id, user_id) VALUES (?, ?)"

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(&event.ID, userId)
	return err
}

func (event Event) Unregister(userId int64) error {
	query := "DELETE FROM registrations WHERE event_id = ? AND user_id = ?"

	stmt, err := database.Db.Prepare(query)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(&event.ID, userId)
	return err
}
