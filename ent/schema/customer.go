package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema/field"
)

// Customer holds the schema definition for the Customer entity.
type Customer struct {
	ent.Schema
}

// Fields of the Customer.
func (Customer) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").Unique().Immutable().Positive(),
		field.String("name").MaxLen(100).NotEmpty(),
		field.String("surname").MaxLen(100).NotEmpty(),
		field.Int("number").Positive(),
		field.Enum("gender").Values("Male", "Female"),
		field.String("country").MaxLen(50).NotEmpty(),
		field.Int("dependants").Default(0).NonNegative(),
		field.Time("birth_date").SchemaType(map[string]string{
			dialect.Postgres: "date",
		}),
	}
}

// Edges of the Customer.
func (Customer) Edges() []ent.Edge {
	return nil
}
