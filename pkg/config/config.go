package config

type Config interface {
	InitConfig()
	GetValue(string) string
	DeleteValue(string)
	SetValue(string, interface{})
	WriteConfig() error
}
