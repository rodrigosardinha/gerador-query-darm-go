package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sirupsen/logrus"
)

// TestDarmProcessor testa o processador principal
func TestDarmProcessor(t *testing.T) {
	// Configurar logging para testes
	logrus.SetLevel(logrus.ErrorLevel)

	t.Run("NewDarmProcessor", testNewDarmProcessor)
	t.Run("Init", testInit)
	t.Run("ExtractDarmData", testExtractDarmData)
	t.Run("GenerateSQLInsert", testGenerateSQLInsert)
	t.Run("ParseMonetaryValue", testParseMonetaryValue)
	t.Run("RemoveLeadingZeros", testRemoveLeadingZeros)
}

// testNewDarmProcessor testa criação do processador
func testNewDarmProcessor(t *testing.T) {
	processor := NewDarmProcessor()

	if processor == nil {
		t.Fatal("Processador não deveria ser nil")
	}

	if processor.BaseDir == "" {
		t.Error("BaseDir não deveria estar vazio")
	}

	if processor.DarmsDir == "" {
		t.Error("DarmsDir não deveria estar vazio")
	}

	if processor.OutputDir == "" {
		t.Error("OutputDir não deveria estar vazio")
	}

	if processor.ProcessedGuias == nil {
		t.Error("ProcessedGuias deveria ser inicializado")
	}

	if processor.GuiasProcessadas == nil {
		t.Error("GuiasProcessadas deveria ser inicializado")
	}

	if processor.AllSQLInserts == nil {
		t.Error("AllSQLInserts deveria ser inicializado")
	}
}

// testInit testa inicialização do processador
func testInit(t *testing.T) {
	// Criar diretório temporário para teste
	tempDir := t.TempDir()

	processor := NewDarmProcessor()
	processor.BaseDir = tempDir
	processor.DarmsDir = filepath.Join(tempDir, "darms")
	processor.OutputDir = filepath.Join(tempDir, "inserts")

	err := processor.Init()
	if err != nil {
		t.Fatalf("Init falhou: %v", err)
	}

	// Verificar se diretórios foram criados
	if _, err := os.Stat(processor.DarmsDir); os.IsNotExist(err) {
		t.Error("Diretório darms deveria ter sido criado")
	}

	if _, err := os.Stat(processor.OutputDir); os.IsNotExist(err) {
		t.Error("Diretório inserts deveria ter sido criado")
	}
}

// testExtractDarmData testa extração de dados do DARM
func testExtractDarmData(t *testing.T) {
	processor := NewDarmProcessor()

	// Texto de exemplo com dados de DARM
	text := `
	02. INSCRIÇÃO MUNICIPAL 123456
	01. RECEITA 262-3
	06. VALOR DO TRIBUTO R$ 1.234,56
	09. VALOR TOTAL R$ 1.234,56
	03. DATA VENCIMENTO 15/12/2024
	04. ANO DE REFERÊNCIA 2025
	05. GUIA NØ
	123456789
	`

	data := processor.extractDarmData(text)

	if data == nil {
		t.Fatal("Dados não deveriam ser nil")
	}

	if data.Inscricao != "123456" {
		t.Errorf("Inscrição esperada: 123456, obtida: %s", data.Inscricao)
	}

	if data.CodigoReceita != "2623" {
		t.Errorf("Código de receita esperado: 2623, obtido: %s", data.CodigoReceita)
	}

	if data.ValorPrincipal != "1.234,56" {
		t.Errorf("Valor principal esperado: 1.234,56, obtido: %s", data.ValorPrincipal)
	}

	if data.ValorTotal != "1.234,56" {
		t.Errorf("Valor total esperado: 1.234,56, obtido: %s", data.ValorTotal)
	}

	if data.DataVencimento != "15/12/2024" {
		t.Errorf("Data de vencimento esperada: 15/12/2024, obtida: %s", data.DataVencimento)
	}

	if data.Exercicio != "2025" {
		t.Errorf("Exercício esperado: 2025, obtido: %s", data.Exercicio)
	}

	if data.NumeroGuia != "123456789" {
		t.Errorf("Número da guia esperado: 123456789, obtido: %s", data.NumeroGuia)
	}
}

// testGenerateSQLInsert testa geração de SQL INSERT
func testGenerateSQLInsert(t *testing.T) {
	processor := NewDarmProcessor()

	data := &DarmData{
		Inscricao:      "123456",
		CodigoReceita:  "2623",
		ValorPrincipal: "1.234,56",
		ValorTotal:     "1.234,56",
		DataVencimento: "15/12/2024",
		Exercicio:      "2025",
		NumeroGuia:     "123456789",
	}

	sql := processor.generateSQLInsert(data)

	if sql == "" {
		t.Fatal("SQL não deveria estar vazio")
	}

	// Verificar se contém elementos essenciais
	if !contains(sql, "INSERT INTO FarrDarmsPagos") {
		t.Error("SQL deveria conter INSERT INTO FarrDarmsPagos")
	}

	if !contains(sql, "123456") {
		t.Error("SQL deveria conter a inscrição")
	}

	if !contains(sql, "123456789") {
		t.Error("SQL deveria conter o número da guia")
	}

	if !contains(sql, "1234.56") {
		t.Error("SQL deveria conter o valor formatado")
	}

	if !contains(sql, "2024-12-15") {
		t.Error("SQL deveria conter a data formatada")
	}
}

// testParseMonetaryValue testa parse de valores monetários
func testParseMonetaryValue(t *testing.T) {
	processor := NewDarmProcessor()

	tests := []struct {
		input    string
		expected string
	}{
		{"R$ 1.234,56", "1234.56"},
		{"1.234,56", "1234.56"},
		{"1234,56", "1234.56"},
		{"1234.56", "1234.56"},
		{"", "0.00"},
		{"R$ 0,00", "0.00"},
		{"R$ 1.000.000,00", "1000000.00"},
	}

	for _, test := range tests {
		result := processor.parseMonetaryValue(test.input)
		if result != test.expected {
			t.Errorf("parseMonetaryValue(%s) = %s, esperado %s", test.input, result, test.expected)
		}
	}
}

// testRemoveLeadingZeros testa remoção de zeros à esquerda
func testRemoveLeadingZeros(t *testing.T) {
	processor := NewDarmProcessor()

	tests := []struct {
		input    string
		expected string
	}{
		{"00123", "123"},
		{"123", "123"},
		{"000", "0"},
		{"0", "0"},
		{"", ""},
	}

	for _, test := range tests {
		result := processor.removeLeadingZeros(test.input)
		if result != test.expected {
			t.Errorf("removeLeadingZeros(%s) = %s, esperado %s", test.input, result, test.expected)
		}
	}
}

// TestUtils testa utilitários
func TestUtils(t *testing.T) {
	t.Run("FileUtils", testFileUtils)
	t.Run("StringUtils", testStringUtils)
	t.Run("DateUtils", testDateUtils)
	t.Run("ValidationUtils", testValidationUtils)
	t.Run("SQLUtils", testSQLUtils)
}

// testFileUtils testa utilitários de arquivo
func testFileUtils(t *testing.T) {
	fileUtils := NewFileUtils()
	tempDir := t.TempDir()

	// Test EnsureDir
	dir := filepath.Join(tempDir, "testdir")
	err := fileUtils.EnsureDir(dir)
	if err != nil {
		t.Errorf("EnsureDir falhou: %v", err)
	}

	// Test FileExists
	if !fileUtils.FileExists(dir) {
		t.Error("FileExists deveria retornar true para diretório existente")
	}

	if fileUtils.FileExists(filepath.Join(tempDir, "nonexistent")) {
		t.Error("FileExists deveria retornar false para arquivo inexistente")
	}

	// Test GetFileSize
	size, err := fileUtils.GetFileSize(dir)
	if err != nil {
		t.Errorf("GetFileSize falhou: %v", err)
	}
	if size < 0 {
		t.Error("Tamanho do arquivo deveria ser >= 0")
	}
}

// testStringUtils testa utilitários de string
func testStringUtils(t *testing.T) {
	stringUtils := NewStringUtils()

	// Test RemoveAccents
	if stringUtils.RemoveAccents("café") != "cafe" {
		t.Error("RemoveAccents falhou")
	}

	if stringUtils.RemoveAccents("São Paulo") != "Sao Paulo" {
		t.Error("RemoveAccents falhou")
	}

	// Test CleanString
	if stringUtils.CleanString("café@#$%") != "cafe" {
		t.Error("CleanString falhou")
	}

	// Test FormatCurrency
	if stringUtils.FormatCurrency(1234.56) != "1234.56" {
		t.Error("FormatCurrency falhou")
	}

	// Test ParseCurrency
	value, err := stringUtils.ParseCurrency("R$ 1.234,56")
	if err != nil {
		t.Errorf("ParseCurrency falhou: %v", err)
	}
	if value != 1234.56 {
		t.Errorf("ParseCurrency retornou %f, esperado 1234.56", value)
	}
}

// testDateUtils testa utilitários de data
func testDateUtils(t *testing.T) {
	dateUtils := NewDateUtils()

	// Test ParseDateBR
	date, err := dateUtils.ParseDateBR("15/12/2024")
	if err != nil {
		t.Errorf("ParseDateBR falhou: %v", err)
	}

	// Test FormatDateBR
	formatted := dateUtils.FormatDateBR(date)
	if formatted != "15/12/2024" {
		t.Errorf("FormatDateBR retornou %s, esperado 15/12/2024", formatted)
	}

	// Test ConvertDateBRToSQL
	sqlDate, err := dateUtils.ConvertDateBRToSQL("15/12/2024")
	if err != nil {
		t.Errorf("ConvertDateBRToSQL falhou: %v", err)
	}
	if sqlDate != "2024-12-15" {
		t.Errorf("ConvertDateBRToSQL retornou %s, esperado 2024-12-15", sqlDate)
	}
}

// testValidationUtils testa utilitários de validação
func testValidationUtils(t *testing.T) {
	validationUtils := NewValidationUtils()

	// Test IsValidCPF
	if !validationUtils.IsValidCPF("123.456.789-09") {
		t.Error("CPF válido foi rejeitado")
	}

	if validationUtils.IsValidCPF("123.456.789-10") {
		t.Error("CPF inválido foi aceito")
	}

	// Test IsValidCNPJ
	if !validationUtils.IsValidCNPJ("11.222.333/0001-81") {
		t.Error("CNPJ válido foi rejeitado")
	}

	if validationUtils.IsValidCNPJ("11.222.333/0001-82") {
		t.Error("CNPJ inválido foi aceito")
	}

	// Test IsValidEmail
	if !validationUtils.IsValidEmail("test@example.com") {
		t.Error("Email válido foi rejeitado")
	}

	if validationUtils.IsValidEmail("invalid-email") {
		t.Error("Email inválido foi aceito")
	}

	// Test IsValidDate
	if !validationUtils.IsValidDate("15/12/2024") {
		t.Error("Data válida foi rejeitada")
	}

	if validationUtils.IsValidDate("32/13/2024") {
		t.Error("Data inválida foi aceita")
	}
}

// testSQLUtils testa utilitários SQL
func testSQLUtils(t *testing.T) {
	sqlUtils := NewSQLUtils()

	// Test EscapeString
	if sqlUtils.EscapeString("O'Connor") != "O''Connor" {
		t.Error("EscapeString falhou")
	}

	// Test QuoteString
	if sqlUtils.QuoteString("test") != "'test'" {
		t.Error("QuoteString falhou")
	}

	// Test FormatSQLValue
	if sqlUtils.FormatSQLValue("test") != "'test'" {
		t.Error("FormatSQLValue falhou para string")
	}

	if sqlUtils.FormatSQLValue(123) != "123" {
		t.Error("FormatSQLValue falhou para int")
	}

	if sqlUtils.FormatSQLValue(123.45) != "123.45" {
		t.Error("FormatSQLValue falhou para float")
	}

	if sqlUtils.FormatSQLValue(true) != "1" {
		t.Error("FormatSQLValue falhou para bool")
	}

	if sqlUtils.FormatSQLValue(nil) != "NULL" {
		t.Error("FormatSQLValue falhou para nil")
	}

	// Test GeneratePlaceholders
	placeholders := sqlUtils.GeneratePlaceholders(3)
	if placeholders != "?, ?, ?" {
		t.Errorf("GeneratePlaceholders retornou %s, esperado '?, ?, ?'", placeholders)
	}
}

// TestConfig testa configurações
func TestConfig(t *testing.T) {
	t.Run("DefaultConfig", testDefaultConfig)
	t.Run("LoadConfig", testLoadConfig)
	t.Run("SaveConfig", testSaveConfig)
	t.Run("ValidateConfig", testValidateConfig)
}

// testDefaultConfig testa configuração padrão
func testDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config == nil {
		t.Fatal("Config não deveria ser nil")
	}

	if config.Database.Host != "localhost" {
		t.Error("Host padrão deveria ser localhost")
	}

	if config.Database.Port != 3306 {
		t.Error("Porta padrão deveria ser 3306")
	}

	if config.Database.Database != "silfae" {
		t.Error("Database padrão deveria ser silfae")
	}

	if config.Paths.DarmsDir != "darms" {
		t.Error("DarmsDir padrão deveria ser darms")
	}

	if config.Paths.OutputDir != "inserts" {
		t.Error("OutputDir padrão deveria ser inserts")
	}

	if config.SQL.Encoding != "latin1" {
		t.Error("Encoding padrão deveria ser latin1")
	}

	if config.Logging.Level != "info" {
		t.Error("Level padrão deveria ser info")
	}
}

// testLoadConfig testa carregamento de configuração
func testLoadConfig(t *testing.T) {
	// Test com arquivo inexistente
	_, err := LoadConfig("nonexistent.json")
	if err == nil {
		t.Error("LoadConfig deveria falhar com arquivo inexistente")
	}

	// Test com configuração padrão
	config, err := LoadConfig("")
	if err != nil {
		t.Errorf("LoadConfig falhou: %v", err)
	}

	if config == nil {
		t.Fatal("Config não deveria ser nil")
	}
}

// testSaveConfig testa salvamento de configuração
func testSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.json")

	config := DefaultConfig()

	err := SaveConfig(config, configPath)
	if err != nil {
		t.Errorf("SaveConfig falhou: %v", err)
	}

	// Verificar se arquivo foi criado
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Arquivo de configuração deveria ter sido criado")
	}

	// Carregar configuração salva
	loadedConfig, err := LoadConfig(configPath)
	if err != nil {
		t.Errorf("LoadConfig falhou: %v", err)
	}

	if loadedConfig.Database.Host != config.Database.Host {
		t.Error("Configuração carregada não corresponde à salva")
	}
}

// testValidateConfig testa validação de configuração
func testValidateConfig(t *testing.T) {
	config := DefaultConfig()

	err := ValidateConfig(config)
	if err != nil {
		t.Errorf("ValidateConfig falhou com configuração válida: %v", err)
	}

	// Test com configuração inválida
	invalidConfig := DefaultConfig()
	invalidConfig.Paths.BaseDir = ""

	err = ValidateConfig(invalidConfig)
	if err == nil {
		t.Error("ValidateConfig deveria falhar com configuração inválida")
	}
}

// Funções auxiliares
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// BenchmarkDarmProcessor testa performance
func BenchmarkDarmProcessor(b *testing.B) {
	processor := NewDarmProcessor()

	text := `
	02. INSCRIÇÃO MUNICIPAL 123456
	01. RECEITA 262-3
	06. VALOR DO TRIBUTO R$ 1.234,56
	09. VALOR TOTAL R$ 1.234,56
	03. DATA VENCIMENTO 15/12/2024
	04. ANO DE REFERÊNCIA 2025
	05. GUIA NØ
	123456789
	`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.extractDarmData(text)
	}
}

// BenchmarkParseMonetaryValue testa performance do parse de valores monetários
func BenchmarkParseMonetaryValue(b *testing.B) {
	processor := NewDarmProcessor()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		processor.parseMonetaryValue("R$ 1.234,56")
	}
}
