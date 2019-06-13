package persistence

import (
	"github.com/jinzhu/gorm"
	"github.com/pborman/uuid"
	"github.com/splisson/opstic/entities"
)

type EventStoreInterface interface {
	GetEvents() ([]entities.Event, error)
	GetEventByCommitAndCategory(commit string, category string) (*entities.Event, error)
	CreateEvent(event entities.Event) (*entities.Event, error)
}

type EventStoreDB struct {
	db *gorm.DB
}

func NewEventDBStore(db *gorm.DB) *EventStoreDB {
	store := new(EventStoreDB)
	store.db = db
	db.LogMode(true)
	return store
}

func (s *EventStoreDB) GetEvents() ([]entities.Event, error) {
	events := []entities.Event{}
	db := s.db.Table("events").Select("*")
	db.Find(&events)
	return events, nil
}

func (s *EventStoreDB) GetEventByCommitAndCategory(commit string, category string) (*entities.Event, error) {
	event := entities.Event{}
	db := s.db.Table("events").Select("*").Where("commit = ? AND category = ?", commit, category)
	db.Find(&event)
	return &event, nil
}

func (s *EventStoreDB) CreateEvent(event entities.Event) (*entities.Event, error) {
	event.ID = uuid.New()

	tx := s.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil, err
	}

	if err := tx.Create(&event).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &event, tx.Commit().Error
}