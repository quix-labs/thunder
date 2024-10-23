package yaml

import (
	"encoding/json"
	"github.com/quix-labs/thunder"
	"gopkg.in/yaml.v2"
	"io"
	"strings"
)

const ExporterId = "thunder.yaml"

func init() {
	thunder.Exporters.Register(ExporterId, &Exporter{})
}

type Exporter struct {
	w io.Writer
}

func (e *Exporter) Name() string {
	return "YAML"
}

func (e *Exporter) MimeType() string {
	return "application/x-yaml"
}

func (e *Exporter) Load(w io.Writer) error {
	e.w = w
	return nil
}

func (e *Exporter) BeforeAll() error {
	_, err := e.w.Write([]byte("items:\n"))
	return err
}

func (e *Exporter) WriteDocument(doc *thunder.Document, itemPosition uint64) error {
	var data map[string]interface{}
	if err := json.Unmarshal(doc.Json, &data); err != nil {
		return err
	}

	item := yaml.MapSlice{
		{Key: "pkey", Value: doc.Pkey},
		{Key: "data", Value: data},
	}

	yamlData, err := yaml.Marshal([]yaml.MapSlice{item})
	if err != nil {
		return err
	}

	for _, line := range strings.Split(string(yamlData), "\n") {
		if line == "" {
			continue
		}
		if _, err := e.w.Write([]byte("  " + line + "\n")); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) AfterAll() error {
	return nil
}

var (
	_ thunder.Exporter = (*Exporter)(nil)
)
