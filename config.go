package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config representa as configurações do sistema
type Config struct {
	Database DatabaseConfig `json:"database"`
	Paths    PathsConfig    `json:"paths"`
	SQL      SQLConfig      `json:"sql"`
	Logging  LoggingConfig  `json:"logging"`
}

// DatabaseConfig configurações do banco de dados
type DatabaseConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Database string `json:"database"`
	Username string `json:"username"`
	Password string `json:"password"`
	Charset  string `json:"charset"`
}

// PathsConfig configurações de caminhos
type PathsConfig struct {
	BaseDir   string `json:"base_dir"`
	DarmsDir  string `json:"darms_dir"`
	OutputDir string `json:"output_dir"`
	TempDir   string `json:"temp_dir"`
}

// SQLConfig configurações SQL
type SQLConfig struct {
	Encoding       string `json:"encoding"`
	BatchSize      int    `json:"batch_size"`
	UseTransaction bool   `json:"use_transaction"`
	UseIgnore      bool   `json:"use_ignore"`
}

// LoggingConfig configurações de logging
type LoggingConfig struct {
	Level      string `json:"level"`
	Format     string `json:"format"`
	OutputFile string `json:"output_file"`
}

// DefaultConfig retorna configuração padrão
func DefaultConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     3306,
			Database: "silfae",
			Username: "root",
			Password: "",
			Charset:  "latin1",
		},
		Paths: PathsConfig{
			BaseDir:   ".",
			DarmsDir:  "darms",
			OutputDir: "inserts",
			TempDir:   "temp",
		},
		SQL: SQLConfig{
			Encoding:       "latin1",
			BatchSize:      100,
			UseTransaction: true,
			UseIgnore:      true,
		},
		Logging: LoggingConfig{
			Level:      "info",
			Format:     "text",
			OutputFile: "",
		},
	}
}

// LoadConfig carrega configuração do arquivo
func LoadConfig(configPath string) (*Config, error) {
	// Se não especificado, usar configuração padrão
	if configPath == "" {
		return DefaultConfig(), nil
	}

	// Verificar se arquivo existe
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("arquivo de configuração não encontrado: %s", configPath)
	}

	// Ler arquivo
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler arquivo de configuração: %v", err)
	}

	// Parsear JSON
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("erro ao parsear arquivo de configuração: %v", err)
	}

	return &config, nil
}

// SaveConfig salva configuração em arquivo
func SaveConfig(config *Config, configPath string) error {
	// Serializar para JSON
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("erro ao serializar configuração: %v", err)
	}

	// Criar diretório se não existir
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório: %v", err)
	}

	// Escrever arquivo
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("erro ao escrever arquivo de configuração: %v", err)
	}

	return nil
}

// ValidateConfig valida configuração
func ValidateConfig(config *Config) error {
	// Validar caminhos
	if config.Paths.BaseDir == "" {
		return fmt.Errorf("base_dir não pode estar vazio")
	}

	if config.Paths.DarmsDir == "" {
		return fmt.Errorf("darms_dir não pode estar vazio")
	}

	if config.Paths.OutputDir == "" {
		return fmt.Errorf("output_dir não pode estar vazio")
	}

	// Validar banco de dados
	if config.Database.Host == "" {
		return fmt.Errorf("host do banco de dados não pode estar vazio")
	}

	if config.Database.Port <= 0 {
		return fmt.Errorf("porta do banco de dados deve ser maior que 0")
	}

	if config.Database.Database == "" {
		return fmt.Errorf("nome do banco de dados não pode estar vazio")
	}

	// Validar SQL
	if config.SQL.BatchSize <= 0 {
		return fmt.Errorf("batch_size deve ser maior que 0")
	}

	// Validar logging
	validLevels := map[string]bool{
		"debug":   true,
		"info":    true,
		"warning": true,
		"error":   true,
		"fatal":   true,
	}

	if !validLevels[config.Logging.Level] {
		return fmt.Errorf("nível de logging inválido: %s", config.Logging.Level)
	}

	return nil
}

// GetConfigPath retorna caminho padrão do arquivo de configuração
func GetConfigPath() string {
	// Tentar encontrar config.json no diretório atual
	configPath := "config.json"
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}

	// Tentar encontrar config.json no diretório home
	homeDir, err := os.UserHomeDir()
	if err == nil {
		configPath = filepath.Join(homeDir, ".darm-processor", "config.json")
		if _, err := os.Stat(configPath); err == nil {
			return configPath
		}
	}

	// Retornar caminho padrão
	return "config.json"
}

// CreateDefaultConfigFile cria arquivo de configuração padrão
func CreateDefaultConfigFile(configPath string) error {
	config := DefaultConfig()
	return SaveConfig(config, configPath)
}
