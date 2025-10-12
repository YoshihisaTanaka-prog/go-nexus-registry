package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"fmt"
	"time"
)

func getBaseFields(defaultStatus string, allowedStatus []string) []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Unique(),
		field.String("name").NotEmpty(),
		field.Int("v1").NonNegative(),
		field.Int("v2").NonNegative(),
		field.Int("v3").NonNegative(),
		field.String("status").Default(defaultStatus).Validate(func(s string) error {
			for _, v := range allowedStatus {
				if s == v {
					return nil
				}
			}
			return fmt.Errorf("invalid status: %s", s)
		}),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now).UpdateDefault(time.Now),
	}
}

// RequestedLibrary holds the schema definition for the RequestedLibrary entity.
type RequestedLibrary struct {
	ent.Schema
}

// Fields of the RequestedLibrary.
func (RequestedLibrary) Fields() []ent.Field {
	baseFields := getBaseFields("status1", []string{"status1", "status2", "status3"})
	return append(
		baseFields,
		field.String("version"),
		field.String("requested_by").NotEmpty(),
	)
}

// Edges of the RequestedLibrary.
func (RequestedLibrary) Edges() []ent.Edge {
	return nil
}

// SavedLibrary holds the schema definition for the SavedLibrary entity.
type SavedLibrary struct {
	ent.Schema
}

// Fields of the Library.
func (SavedLibrary) Fields() []ent.Field {
	baseFields := getBaseFields("status1", []string{"status1", "status2", "status3"})
	return append(
		baseFields,
		field.String("version").NotEmpty(),
	)
}

// Edges of the SavedLibrary.
func (SavedLibrary) Edges() []ent.Edge {
	return nil
}
