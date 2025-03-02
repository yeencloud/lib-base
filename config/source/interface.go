package source

type ConfigInterface interface {
	ReadString(key string) (string, error)
}
