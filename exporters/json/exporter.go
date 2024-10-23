package json

import (
	"encoding/json"
	"fmt"
	"github.com/quix-labs/thunder"
	"io"
)

const ExporterId = "thunder.json"

func init() {
	thunder.Exporters.Register(ExporterId, &Exporter{})
}

type Exporter struct {
	w io.Writer
}

func (e *Exporter) Name() string {
	return "JSON"
}

func (e *Exporter) MimeType() string {
	return "application/json"
}

func (e *Exporter) Load(w io.Writer) error {
	e.w = w
	return nil
}

func (e *Exporter) BeforeAll() error {
	_, err := e.w.Write([]byte("["))
	return err
}

func (e *Exporter) WriteDocument(doc *thunder.Document, itemPosition uint64) error {
	// Prepend comma if not first line
	if itemPosition > 1 {
		if _, err := e.w.Write([]byte(",")); err != nil {
			return err
		}
	}

	sanitizedPkey, err := json.Marshal(doc.Pkey)
	if err != nil {
		return err
	}

	_, err = e.w.Write([]byte(fmt.Sprintf(`{"pkey":%s,"data":%s}`, sanitizedPkey, doc.Json)))
	return err
}

func (e *Exporter) AfterAll() error {
	_, err := e.w.Write([]byte("]"))
	return err
}

var (
	_ thunder.Exporter = (*Exporter)(nil)
)
