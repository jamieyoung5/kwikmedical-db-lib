package config

const (
	DbUserName        = EnvVarPrefix + "POSTGRESQL_USERNAME"
	DbUserNameDefault = ""

	DbPassword        = EnvVarPrefix + "POSTGRESQL_PASSWORD"
	DbPasswordDefault = ""

	DbHost        = EnvVarPrefix + "POSTGRESQL_HOST"
	DbHostDefault = "localhost"

	DbDatabaseName = "neondb"
)

type Config struct {
	UserName     string
	Password     string
	Engine       string
	Host         string
	DatabaseName string
}

func NewConfig() *Config {
	av := NewAppViper()
	av.SetAndBindDefaults(&map[string]interface{}{
		DbUserName: DbUserNameDefault,
		DbPassword: DbPasswordDefault,
		DbHost:     DbHostDefault,
	})
	config := Config{
		UserName:     av.GetString(DbUserName),
		Password:     av.GetString(DbPassword),
		Host:         av.GetString(DbHost),
		DatabaseName: DbDatabaseName,
	}

	return &config
}
