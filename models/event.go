package models

import (
	"fmt"
	"learninggo/db"
	"time"
)

type Event struct {
	ID          int64
	Name        string    `binding:"required"`
	Description string    `binding:"required"`
	Location    string    `binding:"required"`
	DateTime    time.Time `binding:"required"`
	UserID      int
}

var events []Event = []Event{}

func (e *Event) Save() error {
	query := `INSERT INTO events(name, description, location, dateTime, user_id)
              VALUES (?, ?, ?, ?, ?)`
	stmt, err := db.DB.Prepare(query)
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
	query := `SELECT * FROM events`
	rows, err := db.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var event Event
		err := rows.Scan(
			&event.ID,
			&event.Name,
			&event.Description, // ✅ Include this!
			&event.Location,
			&event.DateTime,
			&event.UserID,
		)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}
	return events, nil
}

func GetEventById(id int64) (*Event, error) {
	query := `SELECT * FROM events WHERE id=?`
	row := db.DB.QueryRow(query, id)
	var event Event
	err := row.Scan(
		&event.ID,
		&event.Name,
		&event.Description,
		&event.Location,
		&event.DateTime,
		&event.UserID,
	)
	if err != nil {
		return nil, err
	}
	return &event, nil
}

func DeleteById(id int64) error {
	query := `DELETE FROM events WHERE id = ?`
	result, err := db.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no event found with ID %d", id)
	}

	return nil
}

func EditById(id int64, event Event) error {
	query := `UPDATE events SET name=?, description=?, location=?, user_id=? WHERE id=?`
	result, err := db.DB.Exec(query, id, event.Name, event.Description, event.Location, event.UserID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no event found with ID %d", id)
	}
	return nil
}

func (e Event) Register(userId int64) error {
	query := `INSERT INTO registration(event_id, user_id) VALUES (?, ?)`
	stmt, err := db.DB.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(e.ID, userId)
	return err
}
