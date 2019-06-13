package services

import (
	"github.com/splisson/opstic/entities"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)
var (
	testEvent = entities.Event{
		Category: entities.EVENT_CATEGORY_DEPLOY,
		Timestamp: time.Now(),
		PipelineId:		"test_pipeline",
		Status:			"success",
		Commit:			"1234567890",
		Environment:     "unit_test",
		}
)

func TestCreateEventWithLeadTime(t *testing.T) {

	eventService := NewEventService(testEventStore)

	t.Run("success deploy should fill lead time if there is a build for that commit", func(t *testing.T) {
		newEvent := testEvent
		newEvent.Category = entities.EVENT_CATEGORY_BUILD
		newEvent.Status = "success"
		newEvent.Commit = "123success"
		newEvent.Timestamp = time.Now().Add(-5 * time.Minute)
		event, err := eventService.CreateEvent(newEvent)
		assert.Nil(t, err, "no error")
		newEvent = testEvent
		newEvent.Category = entities.EVENT_CATEGORY_DEPLOY
		newEvent.Status = "success"
		newEvent.Commit = "123success"
		newEvent.Timestamp = time.Now()
		event, err = eventService.CreateEvent(newEvent)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, event.ID, "id should not be nil")
		assert.NotEmpty(t, event.ID, "id should not be empty")
		assert.True(t, event.LeadTime > 0, "lead time should be > 0")
	})

	t.Run("success deploy should not fill lead time if there is not a build for that commit", func(t *testing.T) {
		newEvent := testEvent
		newEvent.Category = entities.EVENT_CATEGORY_BUILD
		newEvent.Status = "success"
		newEvent.Commit = "different"
		newEvent.Timestamp = time.Now().Add(-5 * time.Minute)
		event, err := eventService.CreateEvent(newEvent)
		assert.Nil(t, err, "no error")
		newEvent = testEvent
		newEvent.Category = entities.EVENT_CATEGORY_DEPLOY
		newEvent.Status = "success"
		newEvent.Commit = "same"
		newEvent.Timestamp = time.Now()
		event, err = eventService.CreateEvent(newEvent)
		assert.Nil(t, err, "no error")
		assert.NotNil(t, event.ID, "id should not be nil")
		assert.NotEmpty(t, event.ID, "id should not be empty")
		assert.True(t, event.LeadTime == 0, "lead time should be == 0")
	})
}