package events

// Events
const (
	EVENT_USER_CREATED = "user_created"
	EVENT_USER_UPDATED = "user_updated"
	EVENT_USER_DELETED = "user_deleted"
)

type Event interface {
	Attrs() map[string]string
}

type UserCreated struct {
	ID string
}

func (u UserCreated) Attrs() map[string]string {
	return map[string]string{
		"event":   EVENT_USER_CREATED,
		"user_id": u.ID,
	}
}

type UserUpdated struct {
	ID string
}

func (u UserUpdated) Attrs() map[string]string {
	return map[string]string{
		"event":   EVENT_USER_UPDATED,
		"user_id": u.ID,
	}
}

type UserDeleted struct {
	ID string
}

func (u UserDeleted) Attrs() map[string]string {
	return map[string]string{
		"event":   EVENT_USER_DELETED,
		"user_id": u.ID,
	}
}
