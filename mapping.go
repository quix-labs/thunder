package thunder

type MappingField struct {
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

	UsePivotTable   bool   `json:"use_pivot_table,omitempty"`
	PivotTable      string `json:"pivot_table,omitempty"`
	ForeignPivotKey string `json:"foreign_pivot_key,omitempty"`
	LocalPivotKey   string `json:"local_pivot_key,omitempty"`

	ForeignKey string `json:"foreign_key,omitempty"`

	Mapping Mapping `json:"mapping,omitempty"`
}

type Mapping []MappingField
