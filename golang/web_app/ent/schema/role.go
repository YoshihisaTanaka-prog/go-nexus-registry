package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"fmt"
	"time"
)

// Role holds the schema definition for the Role entity.
type Role struct {
	ent.Schema
}

// Fields of the Role.
func (Role) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
		field.String("mode").NotEmpty().Validate(func(s string) error {
			for _, v := range []string{"viewers", "editors", "admins", "apis", "custom"} {
				if s == v {
					return nil
				}
			}
			return fmt.Errorf("invalid mode: %s", s)
		}),
		field.Bool("isForNexus").Immutable(),
		field.String("requestedBy").NotEmpty(),
		field.Time("createdAt").Default(time.Now).Immutable(),
		field.Time("updatedAt").Default(time.Now).UpdateDefault(time.Now),
	}
}

// Indexes of the Role.
func (Role) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name", "isForNexus").Unique().
		StorageKey("roles_name_is_for_nexus"),
	}
}

// Edges of the Role.
func (Role) Edges() []ent.Edge {
	return nil
}