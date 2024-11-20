package config

const (
	DbUserName        = EnvVarPrefix + "POSTGRESQL_USERNAME"
	DbUserNameDefault = ""

	DbPassword        = EnvVarPrefix + "POSTGRESQL_PASSWORD"
	DbPasswordDefault = ""

	DbHost        = EnvVarPrefix + "POSTGRESQL_HOST"
	DbHostDefault = "localhost"

	DbPort        = EnvVarPrefix + "POSTGRESQL_PORT"
	DbPortDefault = int32(123)

	DbDatabaseName = "kwikmedicaldb"
)

type Config struct {
	UserName     string
	Password     string
	Engine       string
	Host         string
	Port         int32
	DatabaseName string
}

func NewConfig() *Config {
	av := NewAppViper()
	av.SetAndBindDefaults(&map[string]interface{}{
		DbUserName: DbUserNameDefault,
		DbPassword: DbPasswordDefault,
		DbHost:     DbHostDefault,
		DbPort:     DbPortDefault,
	})
	config := Config{
		UserName:     av.GetString(DbUserName),
		Password:     av.GetString(DbPassword),
		Host:         av.GetString(DbHost),
		Port:         av.GetInt32(DbPort),
		DatabaseName: DbDatabaseName,
	}

	return &config
}
