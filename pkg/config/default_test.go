package config

import (
	"errors"
	"os"
	"testing"

	"github.com/hawkv6/clab-telemetry-linker/pkg/helpers"
	"github.com/hawkv6/clab-telemetry-linker/pkg/logging"
	"github.com/knadh/koanf"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestDefaultConfig_setUserHome(t *testing.T) {
	tests := []struct {
		name      string
		want      string
		wantError bool
	}{
		{
			name:      "Test set correct user home",
			want:      "/home/hawkv6",
			wantError: false,
		},
		{
			name:      "Test set incorrect user home",
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			helper := helpers.NewMockHelper(ctrl)
			config := &DefaultConfig{helper: helper}
			if !tt.wantError {
				helper.EXPECT().GetUserHome().Return(nil, tt.want)
				assert.NoError(t, config.setUserHome())
				assert.Equal(t, tt.want, config.userHome)
			}
			if tt.wantError {
				helper.EXPECT().GetUserHome().Return(errors.New("artificial error"), "")
				assert.Error(t, config.setUserHome())
			}
		})
	}
}

func TestDefaultConfig_setConfigFileName(t *testing.T) {
	type args struct {
		configName string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test with empty config name",
			args: args{
				configName: "",
			},
			want: "config.yaml",
		},
		{
			name: "Test without yaml suffix",
			args: args{configName: "config"},
			want: "config.yaml",
		},
		{
			name: "Test with yaml suffix",
			args: args{configName: "config.yaml"},
			want: "config.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{}
			config.setConfigFileName(tt.args.configName)
			assert.Equal(t, tt.want, config.fileName)
		})
	}
}

func TestDefaultConfig_setConfigPath(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test with correct user home",
			want: "/home/hawkv6/.clab-telemetry-linker",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{userHome: "/home/hawkv6"}
			config.setConfigPath()
			assert.Equal(t, tt.want, config.configPath)
		})
	}
}

func TestDefaultConfig_setConfigFileLocation(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "Test with correct file location",
			want: "/home/hawkv6/.clab-telemetry-linker/config.yaml",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{configPath: "/home/hawkv6/.clab-telemetry-linker", fileName: "config.yaml"}
			config.setConfigFileLocation()
			assert.Equal(t, tt.want, config.fullFileLocation)
		})
	}
}

func TestDefaultConfig_setClabName(t *testing.T) {
	type fields struct {
		name string
	}
	type args struct {
		clabName string
	}
	tests := []struct {
		fields fields
		name   string
		args   args
		want   string
	}{
		{
			name: "Test with no clab name in config",
			args: args{
				clabName: "",
			},
			want: "clab-hawkv6",
		},
		{
			fields: fields{name: "clab-hawkv6"},
			name:   "Test with clab-hawkv6 in config (no override)",
			args: args{
				clabName: "",
			},
			want: "clab-hawkv6",
		},
		{
			fields: fields{name: "clab-hawkv6"},
			name:   "Test with clab-hawkv6 in config (override)",
			args: args{
				clabName: "test",
			},
			want: "clab-hawkv6",
		},
		{
			fields: fields{name: "clab-hawkv6"},
			name:   "Test with identical name",
			args: args{
				clabName: "clab-hawkv6",
			},
			want: "clab-hawkv6",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:           logging.DefaultLogger.WithField("subsystem", "config_test"),
				koanfInstance: koanf.New("."),
				helper:        helpers.NewDefaultHelper(),
			}
			config.clabNameKey = config.helper.GetDefaultClabNameKey()
			if tt.fields.name != "" {
				if err := config.koanfInstance.Set(config.clabNameKey, tt.fields.name); err != nil {
					t.Fatal(err)
				}
			}
			assert.NoError(t, config.setClabName(tt.args.clabName))
			assert.Equal(t, tt.want, config.clabName)
		})
	}
}

func TestDefaultConfig_doesConfigExists(t *testing.T) {
	type fields struct {
		fullFileLocation string
	}
	tests := []struct {
		name      string
		fields    fields
		exists    bool
		wantError bool
	}{
		{
			name: "Test with file which exists",
			fields: fields{
				fullFileLocation: "/",
			},
			exists:    true,
			wantError: false,
		},
		{
			name: "Test with file which doesn't exists",
			fields: fields{
				fullFileLocation: "/doesnotexist",
			},
			wantError: false,
			exists:    false,
		},
		{
			name: "Test with file without permissions",
			fields: fields{
				fullFileLocation: "/root/.profile",
			},
			wantError: true,
			exists:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:              logging.DefaultLogger.WithField("subsystem", "config_test"),
				fullFileLocation: tt.fields.fullFileLocation,
			}
			err, exists := config.doesConfigExist()
			if tt.wantError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.exists, exists)
			}
		})
	}
}

func TestDefaultConfig_createConfig(t *testing.T) {
	type fields struct {
		configPath       string
		fullFileLocation string
	}
	tests := []struct {
		fields    fields
		name      string
		wantError bool
	}{
		{
			name: "Test valid file creation",
			fields: fields{
				configPath:       "/tmp/hawkv6",
				fullFileLocation: "/tmp/hawkv6/config.yaml",
			},
			wantError: false,
		},
		{
			name: "Test invalid folder creation",
			fields: fields{
				configPath:       "/root/hawkv6",
				fullFileLocation: "/root/hawkv6/config.yaml",
			},
			wantError: true,
		},
		{
			name: "Test invalid file creation",
			fields: fields{
				configPath:       "/",
				fullFileLocation: "/config.yaml",
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:              logging.DefaultLogger.WithField("subsystem", "config_test"),
				configPath:       tt.fields.configPath,
				fullFileLocation: tt.fields.fullFileLocation,
			}
			if tt.wantError {
				assert.Error(t, config.createConfig())
			} else {
				assert.NoError(t, config.createConfig())
				assert.NoError(t, os.RemoveAll(config.configPath))
			}
		})
	}
}
func TestDefaultConfig_readConfig(t *testing.T) {
	type fields struct {
		configPath       string
		fullFileLocation string
	}
	tests := []struct {
		fields    fields
		name      string
		want      string
		wantError bool
	}{
		{
			name: "Test read valid config",
			fields: fields{
				configPath:       ".",
				fullFileLocation: "config-example.yaml",
			},
			want:      "clab-hawkv6",
			wantError: false,
		},
		{
			name: "Test read invalid config",
			fields: fields{
				fullFileLocation: "/nofileexists",
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:              logging.DefaultLogger.WithField("subsystem", "config_test"),
				koanfInstance:    koanf.New("."),
				configPath:       tt.fields.configPath,
				fullFileLocation: tt.fields.fullFileLocation,
				helper:           helpers.NewDefaultHelper(),
			}
			if tt.wantError {
				assert.Error(t, config.readConfig())
			} else {
				assert.NoError(t, config.readConfig())
				assert.Equal(t, config.koanfInstance.String(config.helper.GetDefaultClabNameKey()), tt.want)
			}
		})
	}
}

func TestDefaultConfig_initConfig(t *testing.T) {
	type fields struct {
		configPath       string
		fullFileLocation string
	}
	tests := []struct {
		fields    fields
		name      string
		wantError bool
	}{
		{
			name: "Test init valid config",
			fields: fields{
				configPath:       "/tmp/hawkv6",
				fullFileLocation: "/tmp/hawkv6/config.yaml",
			},
			wantError: false,
		},
		{
			name: "Test init invalid config",
			fields: fields{
				fullFileLocation: "/nofileexists",
			},
			wantError: true,
		},
		{
			name: "Test with file without permissions",
			fields: fields{
				fullFileLocation: "/root/.profile",
			},
			wantError: true,
		},
		{
			name: "Test non parseable config",
			fields: fields{
				fullFileLocation: "/root",
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:              logging.DefaultLogger.WithField("subsystem", "config_test"),
				koanfInstance:    koanf.New("."),
				configPath:       tt.fields.configPath,
				fullFileLocation: tt.fields.fullFileLocation,
			}
			if tt.wantError {
				assert.Error(t, config.InitConfig())
			} else {
				assert.NoError(t, config.InitConfig())
			}
		})
	}
}

func TestDefaultConfig_SetValue(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name  string
		args  args
		wants string
	}{
		{
			name: "Test set value",
			args: args{
				key:   "test",
				value: "test",
			},
			wants: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:           logging.DefaultLogger.WithField("subsystem", "config_test"),
				koanfInstance: koanf.New("."),
			}
			assert.NoError(t, config.SetValue(tt.args.key, tt.args.value))
			assert.Equal(t, tt.wants, config.koanfInstance.String(tt.args.key))
		})
	}
}

func TestDefaultConfig_GetValue(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name  string
		args  args
		wants string
	}{
		{
			name: "Test get value",
			args: args{
				key:   "test",
				value: "test",
			},
			wants: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:           logging.DefaultLogger.WithField("subsystem", "config_test"),
				koanfInstance: koanf.New("."),
			}
			assert.NoError(t, config.koanfInstance.Set(tt.args.key, tt.args.value))
			assert.Equal(t, tt.wants, config.GetValue(tt.args.key))
		})
	}
}

func TestDefaultConfig_DeleteValue(t *testing.T) {
	type args struct {
		key   string
		value interface{}
	}
	tests := []struct {
		name  string
		args  args
		wants string
	}{
		{
			name: "Test delete value",
			args: args{
				key:   "test",
				value: "test",
			},
			wants: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:           logging.DefaultLogger.WithField("subsystem", "config_test"),
				koanfInstance: koanf.New("."),
			}
			assert.NoError(t, config.koanfInstance.Set(tt.args.key, tt.args.value))
			assert.Equal(t, tt.args.value, config.GetValue(tt.args.key))
			config.DeleteValue(tt.args.key)
			assert.Equal(t, tt.wants, config.GetValue(tt.args.key))
		})
	}
}
func TestDefaultConfig_WriteConfig(t *testing.T) {
	type fields struct {
		fullFileLocation string
		configPath       string
	}
	type args struct {
		key   string
		value string
	}
	tests := []struct {
		fields    fields
		args      args
		name      string
		wantError bool
	}{
		{
			name: "Write valid config",
			fields: fields{
				configPath:       "/tmp/hawkv6/",
				fullFileLocation: "/tmp/hawkv6/config.yaml",
			},
			args: args{
				key:   "test",
				value: "test",
			},
			wantError: false,
		},
		{
			name: "Test init invalid config",
			fields: fields{
				fullFileLocation: "/root/.profile",
			},
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DefaultConfig{
				log:              logging.DefaultLogger.WithField("subsystem", "config_test"),
				koanfInstance:    koanf.New("."),
				configPath:       tt.fields.configPath,
				fullFileLocation: tt.fields.fullFileLocation,
			}
			if !tt.wantError {
				assert.NoError(t, config.InitConfig())
				assert.NoError(t, config.koanfInstance.Set(tt.args.key, tt.args.value))
				assert.NoError(t, config.WriteConfig())
				config.koanfInstance = koanf.New(".")
				assert.Equal(t, "", config.koanfInstance.String(tt.args.key))
				assert.NoError(t, config.readConfig())
				assert.Equal(t, tt.args.key, config.koanfInstance.String(tt.args.key))
				assert.NoError(t, os.RemoveAll(tt.fields.configPath))
			} else {
				assert.Error(t, config.WriteConfig())
			}
		})
	}
}

func TestDefaultConfig_createDefaultConfig(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
		want      *DefaultConfig
	}{
		{
			name:      "Create Valid config",
			wantError: false,
			want: &DefaultConfig{
				userHome:         "/tmp/hawkv6",
				fileName:         "config.yaml",
				configPath:       "/tmp/hawkv6/.clab-telemetry-linker",
				fullFileLocation: "/tmp/hawkv6/.clab-telemetry-linker/config.yaml",
				clabNameKey:      "clab-name",
				clabName:         "clab-hawkv6",
			},
		},
		{
			name:      "Wrong user home",
			wantError: true,
			want: &DefaultConfig{
				userHome:    "",
				fileName:    "config.yaml",
				clabNameKey: "clab-name",
				clabName:    "clab-hawkv6",
			},
		},
		{
			name:      "Invalid config",
			wantError: true,
			want: &DefaultConfig{
				userHome:    "/root",
				fileName:    "config.yaml",
				clabNameKey: "clab-name",
				clabName:    "clab-hawkv6",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			helper := helpers.NewMockHelper(ctrl)
			helper.EXPECT().GetDefaultClabNameKey().Return(tt.want.clabNameKey)
			if tt.want.userHome == "" {
				helper.EXPECT().GetUserHome().Return(errors.New("artificial error"), "")
			} else {
				helper.EXPECT().GetUserHome().Return(nil, tt.want.userHome)
			}
			if tt.wantError {
				if tt.want.userHome == "/root" {
					helper.EXPECT().GetDefaultClabName().Return(tt.want.clabName)
				}
				err, _ := CreateDefaultConfig(tt.want.fileName, tt.want.clabName, tt.want.clabNameKey, helper)
				assert.Error(t, err)
			} else {
				helper.EXPECT().GetDefaultClabName().Return(tt.want.clabName)
				err, defaultConfig := CreateDefaultConfig(tt.want.fileName, tt.want.clabName, tt.want.clabNameKey, helper)
				assert.NoError(t, err)
				assert.NotNil(t, defaultConfig)
				assert.Equal(t, tt.want.userHome, defaultConfig.userHome)
				assert.Equal(t, tt.want.fileName, defaultConfig.fileName)
				assert.Equal(t, tt.want.configPath, defaultConfig.configPath)
				assert.Equal(t, tt.want.fullFileLocation, defaultConfig.fullFileLocation)
				assert.Equal(t, tt.want.clabNameKey, defaultConfig.clabNameKey)
				assert.Equal(t, tt.want.clabName, defaultConfig.clabName)

			}
		})
	}
}
func TestDefaultConfig_NewDefaultConfig(t *testing.T) {
	tests := []struct {
		name      string
		wantError bool
	}{
		{
			name:      "Test NewDefaultConfig",
			wantError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err, config := NewDefaultConfig()
			if tt.wantError {
				assert.Error(t, err)
				assert.Nil(t, config)
			}
		})
	}
}
