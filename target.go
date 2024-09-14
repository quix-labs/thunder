package thunder

type Target struct {
	ID     int
	Driver TargetDriver
	Config DynamicConfig
}
