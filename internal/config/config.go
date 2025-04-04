package config

type Config struct {
    DBHost     string
    DBUser     string
    DBPassword string
    DBName     string
    DBPort     string
}

func LoadConfig() *Config {
    return &Config{
        DBHost:     "localhost8080",
        DBUser:     "postgres",
        DBPassword: "mobin2005G",
        DBName:     "booking",
        DBPort:     "5432",
    }
}