package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func Now() time.Time {
	return time.Now().In(time.UTC)
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id"),
		field.String("username"),
		field.String("email").Unique(),
		field.Time("created_at").Default(Now),
		field.Time("updated_at").Default(Now).UpdateDefault(Now),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return nil
}
