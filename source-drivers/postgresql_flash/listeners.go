package postgresql_flash

import (
	"github.com/quix-labs/flash"
	"github.com/quix-labs/thunder"
)

type RealtimeConfigItem struct {
	Table       string
	PrimaryKeys []string

	ListenerConfig *flash.ListenerConfig
}
type RealtimeConfig map[*thunder.Relation]*RealtimeConfigItem

func GetRealtimeConfigForProcessor(p *thunder.Processor) (RealtimeConfig, error) {
	var configs = make(RealtimeConfig)

	// TODO CONDITIONS
	// Recursively append configs
	if err := appendRelationConfig(nil, p.Table, &p.Mapping, &configs, p.PrimaryKeys); err != nil {
		return nil, err
	}

	return configs, nil
}

func appendRelationConfig(
	relation *thunder.Relation,
	table string,
	mapping *thunder.Mapping,
	configs *RealtimeConfig,
	primaryKeys []string,
) error {
	baseFields, err := extractMappingColumns(mapping, primaryKeys)
	if err != nil {
		return err
	}
	scopedConfig := &flash.ListenerConfig{
		Table:  table,
		Fields: baseFields,
	}

	for _, relation := range mapping.Relations {
		// TODO PIVOT TABLES
		if err := appendRelationConfig(&relation, relation.Table, &relation.Mapping, configs, relation.PrimaryKeys); err != nil {
			return err
		}

	}

	(*configs)[relation] = &RealtimeConfigItem{
		Table:          table,
		PrimaryKeys:    primaryKeys,
		ListenerConfig: scopedConfig,
	}

	return nil
}

func extractMappingColumns(m *thunder.Mapping, additional []string) ([]string, error) {
	var columns = additional
	for _, field := range m.Fields {
		columns = append(columns, field.Column)
	}
	return columns, nil
}
