package thunder

type Source struct {
	ID     int
	Driver SourceDriver
	Config DynamicConfig
}
