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
	// Environments
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"

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
	DB struct {
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

	// Определяем флаг для окружения
	env := flag.String("env", "", "application environment (local/dev/prod)")
	flag.Parse()

	// Если флаг не установлен, проверяем переменную окружения
	if *env == "" {
		*env = os.Getenv("APP_ENV")
	}

	wd, err := os.Getwd()

	if err != nil {
		log.Fatalf("Failed to get work directory: %v", err)
	}

	envFile := getEnvFileName(*env)
	fullPath := filepath.Join(wd, "../../", envFile)

	if err := godotenv.Load(fullPath); err != nil {
		log.Fatalf("Failed to load .env file %s ", err)
	}

	var cfg Config

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

// getEnvFileName возвращает имя файла конфигурации для указанного окружения.
// Принимает строку с именем окружения (local, dev, prod).
// Возвращает соответствующее имя .env файла.
func getEnvFileName(env string) string {
	switch env {
	case envDev:
		log.Printf("Using %s environment", envDev)
		return ".env.dev"
	case envProd:
		log.Printf("Using %s environment", envProd)
		return ".env.prod"
	case envLocal:
		log.Printf("Using %s environment", envLocal)
		return ".env.local"
	default:
		log.Println("Environment not specified or not recognized, using default (.env)")
		return ".env"
	}
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
