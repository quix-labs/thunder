package thunder

import (
	"errors"
	"fmt"
)

type JsonMappingField struct {
	FieldType string `json:"_type"`

	// Simple field
	Column string `json:"column,omitempty"`
	Alias  string `json:"alias,omitempty"`

	// Relation
	Name        string   `json:"name,omitempty"`
	Type        string   `json:"type,omitempty"`
	Table       string   `json:"table,omitempty"`
	PrimaryKeys []string `json:"primary_keys,omitempty"`
	LocalKey    string   `json:"local_key,omitempty"`

	UsePivotTable   bool               `json:"use_pivot_table,omitempty"`
	PivotTable      string             `json:"pivot_table,omitempty"`
	ForeignPivotKey string             `json:"foreign_pivot_key,omitempty"`
	LocalPivotKey   string             `json:"local_pivot_key,omitempty"`
	PivotFields     []JsonMappingField `json:"pivot_fields,omitempty"`

	ForeignKey string `json:"foreign_key,omitempty"`

	Mapping JsonMapping `json:"mapping,omitempty"`
}

type JsonMapping []JsonMappingField

func SerializeMapping(m *Mapping) (*JsonMapping, error) {
	var jm = JsonMapping{}

	for _, field := range m.Fields {
		name := ""
		if field.Name != nil {
			name = *field.Name
		}
		jm = append(jm, JsonMappingField{
			FieldType: "simple",
			Column:    field.Column,
			Name:      name,
		})
	}

	for _, rel := range m.Relations {

		jsonRel := JsonMappingField{
			FieldType:   "relation",
			Name:        rel.Name,
			Type:        "one-to-one",
			Table:       rel.Table,
			LocalKey:    rel.LocalKey,
			ForeignKey:  rel.ForeignKey,
			PrimaryKeys: rel.PrimaryKeys,
		}

		if rel.Many {
			jsonRel.Type = "has-many"
		}

		if rel.Pivot != nil {

			nested, err := SerializeMapping(&Mapping{Fields: rel.Pivot.Fields})
			if err != nil {
				return nil, err
			}
			jsonRel.UsePivotTable = true
			jsonRel.ForeignPivotKey = rel.Pivot.ForeignKey
			jsonRel.LocalPivotKey = rel.Pivot.LocalKey
			jsonRel.PivotTable = rel.Pivot.Table
			jsonRel.PivotFields = *nested
		}

		nestedMapping, err := SerializeMapping(&rel.Mapping)
		if err != nil {
			return nil, err
		}
		jsonRel.Mapping = *nestedMapping

		jm = append(jm, jsonRel)
	}

	return &jm, nil
}

func UnserializeMapping(jm *JsonMapping, parent *Relation) (*Mapping, error) {
	var mapping = Mapping{
		Fields:    []SimpleField{},
		Relations: []Relation{},
	}
	for _, field := range *jm {
		if field.FieldType == "simple" {
			var name *string = nil
			if field.Name != "" {
				name = &field.Name
			}

			mapping.Fields = append(mapping.Fields, SimpleField{
				Column: field.Column,
				Name:   name,
			})
			continue
		}

		if field.FieldType == "relation" {
			relation := Relation{
				Name:        field.Name,
				Many:        field.Type == "has-many",
				Table:       field.Table,
				PrimaryKeys: field.PrimaryKeys,

				LocalKey:   field.LocalKey,
				ForeignKey: field.ForeignKey,
				Parent:     parent,
			}

			// Add pivot if needed
			if field.UsePivotTable {
				nestedFields, err := UnserializeMapping(&field.Mapping, &relation)
				if err != nil {
					return nil, err
				}
				relation.Pivot = &RelationPivot{
					Table:      field.PivotTable,
					LocalKey:   field.LocalPivotKey,
					ForeignKey: field.ForeignPivotKey,
					Fields:     nestedFields.Fields, //TODO FRONT
				}
			}

			// Add nested mapping
			nestedMapping, err := UnserializeMapping(&field.Mapping, &relation)
			if err != nil {
				return nil, err
			}
			relation.Mapping = *nestedMapping

			// Inject relation
			mapping.Relations = append(mapping.Relations, relation)
			continue
		}
		return nil, errors.New(fmt.Sprintf("Unknown field type: %s", field.FieldType))
	}

	return &mapping, nil
}
