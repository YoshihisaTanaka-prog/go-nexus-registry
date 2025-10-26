package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"fmt"
	"time"
)

func getBaseFields(defaultStatus string, allowedStatus []string) []ent.Field {
	return []ent.Field{
		field.String("id").NotEmpty().Unique(),
		field.String("kind").NotEmpty(),
		field.String("name").NotEmpty(),
		field.String("version").NotEmpty(),
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
	baseFields := getBaseFields("uploading", []string{"uploading", "uploaded"})
	return append(
		baseFields,
		field.String("requested_by").NotEmpty(),
	)
}

// Indexes of the RequestedLibrary.
func (RequestedLibrary) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("kind", "name", "version").Unique().
		StorageKey("requested_libraries_kind_name_version_idx"),
	}
}

// Edges of the RequestedLibrary.
func (RequestedLibrary) Edges() []ent.Edge {
	return nil
}

// SavedLibrary holds the schema definition for the SavedLibrary entity.
type SavedLibrary struct {
	ent.Schema
}

// Fields of the SavedLibrary.
func (SavedLibrary) Fields() []ent.Field {
	baseFields := getBaseFields("uploading", []string{"uploading", "uploaded", "failed", "updating"})
	return append(
		baseFields,
		field.Int("v1").NonNegative(),
		field.Int("v2").NonNegative(),
		field.Int("v3").NonNegative(),
		field.Bool("isPublished").Default(false),
	)
}

// Indexes of the SavedLibrary.
func (SavedLibrary) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("kind", "name", "version").Unique().
		StorageKey("saved_libraries_kind_name_version_idx"),
	}
}

// Edges of the SavedLibrary.
func (SavedLibrary) Edges() []ent.Edge {
	return nil
}
