package thunder

// FIELDS

type SimpleField struct {
	Column string
	Name   *string // Use Column as Name if empty or nil
}

// RELATIONS

type RelationPivot struct {
	Table      string
	LocalKey   string
	ForeignKey string

	Fields []SimpleField
}

type Relation struct {
	Name string

	Many bool

	Table       string
	PrimaryKeys []string

	LocalKey   string
	Pivot      *RelationPivot
	ForeignKey string

	Mapping Mapping

	// Recursive navigation
	Parent   *Relation
	Children []*Relation
}

func (r *Relation) Path() string {
	path := r.Name
	parent := r.Parent
	for parent != nil {
		path = parent.Name + "." + path
		parent = parent.Parent
	}
	return path
}

// COMPLETE MAPPING

type Mapping struct {
	Fields    []SimpleField
	Relations []Relation
}
