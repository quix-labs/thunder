package csv

import (
	"encoding/csv"
	"github.com/quix-labs/thunder"
	"io"
)

const ExporterId = "thunder.csv"

func init() {
	thunder.Exporters.Register(ExporterId, &Exporter{})
}

type Exporter struct {
	csvWriter *csv.Writer
}

func (e *Exporter) Name() string {
	return "CSV"
}

func (e *Exporter) MimeType() string {
	return "text/csv"
}

func (e *Exporter) Load(w io.Writer) error {
	e.csvWriter = csv.NewWriter(w)
	return nil
}

func (e *Exporter) BeforeAll() error {
	return e.csvWriter.Write([]string{"pkey", "data"})
}

func (e *Exporter) WriteDocument(doc *thunder.Document, itemPosition uint64) error {
	return e.csvWriter.Write([]string{doc.Pkey, string(doc.Json)})
}

func (e *Exporter) AfterAll() error {
	e.csvWriter.Flush()
	return nil
}

var (
	_ thunder.Exporter = (*Exporter)(nil)
)
