package postgresql_flash

import (
	"fmt"
	"github.com/quix-labs/flash"
	"github.com/quix-labs/thunder"
)

type RealtimeConfigItem struct {
	Table       string
	PrimaryKeys []string

	ListenerConfig *flash.ListenerConfig
}
type RealtimeConfig map[string]*RealtimeConfigItem

func GetRealtimeConfigForProcessor(p *thunder.Processor) (RealtimeConfig, error) {
	var configs = make(RealtimeConfig)

	// TODO CONDITIONS
	// Recursively append configs
	if err := appendRelationConfig("", p.Table, &p.Mapping, &configs, p.PrimaryKeys); err != nil {
		return nil, err
	}

	return configs, nil
}

func appendRelationConfig(
	path string,
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

	for _, field := range *mapping {
		if field.FieldType == "relation" {
			var nestedPath string
			if path != "" {
				nestedPath = path + "." + field.Name
			} else {
				nestedPath = field.Name
			}

			// TODO PIVOT TABLES
			if err := appendRelationConfig(nestedPath, field.Table, &field.Mapping, configs, field.PrimaryKeys); err != nil {
				return err
			}
		}
	}

	(*configs)[path] = &RealtimeConfigItem{
		Table:          table,
		PrimaryKeys:    primaryKeys,
		ListenerConfig: scopedConfig,
	}

	return nil
}

func extractMappingColumns(m *thunder.Mapping, additional []string) ([]string, error) {
	var columns = additional
	for _, field := range *m {
		switch field.FieldType {
		case "relation":
			break
		case "simple":
			if field.Name != "" {
				columns = append(columns, field.Name)
			} else if field.Column != "" {
				columns = append(columns, field.Column)
			}
			break
		default:
			return nil, fmt.Errorf("invalid field type: %s", field.FieldType)
		}
	}
	return columns, nil
}
