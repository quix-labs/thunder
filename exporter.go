package thunder

import (
	"github.com/quix-labs/thunder/utils"
	"io"
)

// Exporter is an interface for transforming a set of documents into a file.
type Exporter interface {
	Name() string
	MimeType() string
	Load(w io.Writer) error
	BeforeAll() error
	WriteDocument(doc *Document, itemPosition uint64) error
	AfterAll() error
}

// Exporters is a registry that allows external library to register their own exporter
var Exporters = utils.NewRegistry[Exporter]("exporter")
