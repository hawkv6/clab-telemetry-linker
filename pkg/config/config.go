package config

type Config interface {
	InitConfig() error
	GetValue(string) string
	DeleteValue(string)
	SetValue(string, interface{}) error
	WriteConfig() error
}
