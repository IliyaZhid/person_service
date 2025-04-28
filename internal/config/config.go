package config

import (
	"flag"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

const (
	// EnvLocal локальное окружение
	EnvLocal = "local"
	// EnvDev dev окружение
	EnvDev = "dev"
	// EnvProd prod окружение
	EnvProd = "prod"

	// Database defaults
	defDBHost = "localhost"
	defDBPort = 5432
	defDBUser = "postgres"
	defDBPass = ""
	defDBName = "postgres"

	// Server defaults
	defSrvHost = "localhost"
	defSrvPort = 8080
)

// Config содержит все настройки приложения
type Config struct {
	Environment string
	DB          struct {
		Host string
		Port int
		User string
		Pass string
		Name string
	}
	Server struct {
		Port int
		Host string
	}
}

// MustLoad загружает конфигурацию приложения.
// Принимает окружение через флаг -env или переменную окружения APP_ENV.
// Если окружение не указано, используется локальное (local).
// В случае ошибки загрузки конфигурации вызывает log.Fatal.
// Возвращает указатель на загруженную конфигурацию.
func MustLoad() *Config {
	var cfg Config

	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get work directory: %v", err)
	}

	cfg.Environment = determineEnvironment()
	fullPath := filepath.Join(wd, "../../", getEnvFileName(cfg.Environment))

	if err := godotenv.Load(fullPath); err != nil {
		log.Fatalf("Failed to load .env file %s ", err)
	}

	// Загружаем DB конфигурацию
	cfg.DB.Host = getEnv("DB_HOST", defDBHost)
	cfg.DB.Port = getEnvAsInt("DB_PORT", defDBPort)
	cfg.DB.User = getEnv("DB_USER", defDBUser)
	cfg.DB.Pass = getEnv("DB_PASSWORD", defDBPass)
	cfg.DB.Name = getEnv("DB_NAME", defDBName)

	// Загружаем Server конфигурацию
	cfg.Server.Host = getEnv("SERVER_PORT", defSrvHost)
	cfg.Server.Port = getEnvAsInt("SERVER_PORT", defSrvPort)

	return &cfg
}

func determineEnvironment() string {
	// Определяем флаг для окружения
	env := flag.String("env", "", "application environment (local/dev/prod)")
	flag.Parse()

	// Приоритеты: флаг > переменная окружения > значение по умолчанию
	if *env == "" {
		*env = os.Getenv("APP_ENV")
	}

	switch *env {
	case EnvDev:
		log.Printf("Using %s environment", EnvDev)
		return EnvDev
	case EnvProd:
		log.Printf("Using %s environment", EnvProd)
		return EnvProd
	case EnvLocal:
		log.Printf("Using %s environment", EnvLocal)
		return EnvLocal
	default:
		log.Println("Environment not specified or not recognized, using default")
	}

	return EnvLocal
}

func getEnvFileName(env string) string {
	if env == "" || env == EnvLocal {
		if _, err := os.Stat(".env.local"); err == nil {
			return ".env.local"
		} else {
			return ".env"
		}
	}
	return ".env." + env
}

// getEnv получает значение переменной окружения.
// Если переменная не найдена, возвращает значение по умолчанию.
// key - имя переменной окружения
// defaultVal - значение по умолчанию
func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	log.Printf("Env var %s is missing. Use default value", key)
	return defaultVal
}

// getEnvAsInt получает числовое значение переменной окружения.
// Если переменная не найдена или не может быть преобразована в число,
// возвращает значение по умолчанию.
// key - имя переменной окружения
// defaultVal - значение по умолчанию
func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	log.Printf("Env var %s is missing or not typecastable . Use defaul value", key)
	return defaultVal
}
