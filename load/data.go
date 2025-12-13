// Package load provides functions to load configuration and environment settings.
package load

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/kubex-ecosystem/kbx/get"
	"github.com/kubex-ecosystem/kbx/tools"
	"github.com/kubex-ecosystem/kbx/types"

	gl "github.com/kubex-ecosystem/logz"
)

type MailConfig = types.MailConfig
type MailConnection = types.MailConnection
type Email = types.Email

// ------------------------------- New Mail Srv Params Functions -----------------------------//

type MailSrvParams struct {
	ConfigPath         string `json:"config_path,omitempty" yaml:"config_path,omitempty" xml:"config_path,omitempty" toml:"config_path,omitempty" mapstructure:"config_path,omitempty"`
	*types.Attachment  `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	*types.Email       `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	*types.MailConfig  `json:",inline" yaml:",inline" xml:"-" toml:",inline" mapstructure:",squash"`
	types.MailProvider `json:"-" yaml:"-" xml:"-" toml:"-" mapstructure:"-"`
}

func NewMailSrvParams(configPath string) *MailSrvParams {
	return &MailSrvParams{ConfigPath: configPath, MailConfig: types.NewMailConfig(configPath), Attachment: &types.Attachment{}, Email: &types.Email{}}
}

// ------------------------------- New Mail Params Functions -----------------------------//

func NewMailConfig(configPath string) *MailConfig {
	return &MailConfig{
		ConfigPath:  configPath,
		Provider:    "",
		Connections: make([]*MailConnection, 0),
	}
}

// ------------------------------- New Logz Params Functions -----------------------------//

type LogzConfig = types.LogzConfig

func NewLogzParams() *LogzConfig { return &LogzConfig{} }

func ParseLogzArgs(level string, minLevel string, maxLevel string, output string) *LogzConfig {
	LogzArgs := NewLogzParams()
	LogzArgs.Level = gl.Level(get.ValOrType(level, "info"))
	LogzArgs.MinLevel = gl.Level(get.ValOrType(minLevel, "debug"))
	LogzArgs.MaxLevel = gl.Level(get.ValOrType(maxLevel, "fatal"))
	return LogzArgs
}

// ------------------------------- New Srv Params Functions -----------------------------//

type SrvConfig = types.SrvConfig

func NewSrvArgs() *SrvConfig { return &SrvConfig{} }

func ParseSrvArgs(bind string, pubCertKeyPath string, pubKeyPath string, privKeyPath string, accessTokenTTL int, refreshTokenTTL int, issuer string) *SrvConfig {
	SrvArgs := NewSrvArgs()
	SrvArgs.Bind = get.ValOrType(bind, ":8080")
	SrvArgs.PubCertKeyPath = get.ValOrType(pubCertKeyPath, "")
	SrvArgs.PubKeyPath = get.ValOrType(pubKeyPath, "")
	SrvArgs.PrivKeyPath = get.ValOrType(privKeyPath, "")
	SrvArgs.AccessTokenTTL = time.Duration(get.ValOrType(accessTokenTTL, 15)) * time.Minute
	SrvArgs.RefreshTokenTTL = time.Duration(get.ValOrType(refreshTokenTTL, 60)) * time.Minute
	SrvArgs.Issuer = get.ValOrType(issuer, "kubex-ecosystem")
	return SrvArgs
}

type GlobalRef struct {
	ID   uuid.UUID `json:"id,omitempty"`
	Name string    `json:"name,omitempty"`
}

func NewGlobalRef(name string) *GlobalRef {
	return &GlobalRef{
		ID:   uuid.New(),
		Name: name,
	}
}

func (gr *GlobalRef) GetGlobalRef() *GlobalRef { return gr }
func (gr *GlobalRef) GetName() string          { return gr.Name }
func (gr *GlobalRef) GetID() uuid.UUID         { return gr.ID }
func (gr *GlobalRef) SetName(name string)      { gr.Name = name }
func (gr *GlobalRef) SetID(id uuid.UUID)       { gr.ID = id }
func (gr *GlobalRef) String() string {
	return gr.Name + "-" + gr.ID.String()
}

// ------------------------------- KBX Config Registry -----------------------------//

var configRegistry = map[reflect.Type]bool{
	reflect.TypeFor[MailSrvParams](): true,
	reflect.TypeFor[MailConfig]():    true,
	reflect.TypeFor[LogzConfig]():    true,
	reflect.TypeFor[SrvConfig]():     true,
	reflect.TypeFor[GlobalRef]():     true,
}

var defaultFactories = map[reflect.Type]func() any{
	reflect.TypeFor[MailSrvParams](): func() any { return NewMailSrvParams("") },
	reflect.TypeFor[MailConfig]():    func() any { return NewMailConfig("") },
	reflect.TypeFor[LogzConfig]():    func() any { return NewLogzParams() },
	reflect.TypeFor[SrvConfig]():     func() any { return NewSrvArgs() },
	reflect.TypeFor[GlobalRef]():     func() any { return NewGlobalRef("default") },
}

// LoadConfig loads a configuration of type T from the specified file path.

func LoadConfig[T any](cfgPath string) (*T, error) {
	if configRegistry[reflect.TypeFor[T]()] {
		cfgLoader := get.Loader[T](cfgPath)
		return cfgLoader.DeserializeFromFile(get.FileExt(cfgPath))
	}
	return nil, gl.Errorf("configuration type not registered")
}

func LoadConfigOrDefault[T MailConfig | MailConnection | LogzConfig | SrvConfig | MailSrvParams | Email](cfgPath string, genFile bool) (*T, error) {
	// Só entra aqui se o tipo for algum já registrado, então não me preocupo em checar o erro, só logo retorno o default
	cfgMapper := tools.NewEmptyMapperType[T](cfgPath)
	cfg, err := cfgMapper.DeserializeFromFile(get.FileExt(cfgPath))
	if err == nil {
		return cfg, nil
	}
	gl.Warnf("failed to load config from '%s', using default: %v", cfgPath, err)
	defaultCfg := defaultFactories[reflect.TypeFor[T]()]().(*T)
	if genFile {
		cfgMapper.SetValue(defaultCfg)
		cfgMapper.SerializeToFile(get.FileExt(cfgPath))
	}
	return defaultCfg, nil
}
