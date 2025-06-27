package main

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

// ExemploUso demonstra o uso básico do processador
func ExemploUso() {
	logrus.Info("🎯 Exemplo básico de uso do processador de DARMs")

	// 1. Criar processador
	processor := NewDarmProcessor()

	// 2. Inicializar
	if err := processor.Init(); err != nil {
		logrus.Errorf("❌ Erro ao inicializar: %v", err)
		return
	}

	// 3. Processar DARMs
	if err := processor.ProcessDarms(); err != nil {
		logrus.Errorf("❌ Erro durante o processamento: %v", err)
		return
	}

	// 4. Mostrar resultados
	showResults(processor)

	logrus.Info("✅ Exemplo básico concluído!")
}

// ExemploConfiguracao demonstra como configurar o processador
func ExemploConfiguracao() {
	logrus.Info("⚙️ Exemplo de configuração personalizada")

	// 1. Carregar configuração padrão
	config, err := LoadConfig("config.json")
	if err != nil {
		logrus.Errorf("❌ Erro ao carregar configuração: %v", err)
		return
	}

	// 2. Personalizar configuração
	config.Database.Host = "192.168.1.100"
	config.Database.Port = 3306
	config.Database.Database = "silfae"
	config.SQL.BatchSize = 50
	config.Logging.Level = "debug"

	// 3. Salvar configuração
	if err := SaveConfig(config, "config_personalizada.json"); err != nil {
		logrus.Errorf("❌ Erro ao salvar configuração: %v", err)
	}

	logrus.Info("✅ Configuração personalizada salva em config_personalizada.json")
}

// ExemploUtilitarios demonstra como usar os utilitários
func ExemploUtilitarios() {
	logrus.Info("🛠️ Exemplo de uso dos utilitários")

	// StringUtils
	stringUtils := NewStringUtils()
	logrus.Infof("📝 RemoveAccents: %s", stringUtils.RemoveAccents("café"))
	logrus.Infof("📝 CleanString: %s", stringUtils.CleanString("café@#$%"))
	logrus.Infof("💰 FormatCurrency: %s", stringUtils.FormatCurrency(1234.56))

	// DateUtils
	dateUtils := NewDateUtils()
	sqlDate, err := dateUtils.ConvertDateBRToSQL("15/12/2024")
	if err != nil {
		logrus.Errorf("❌ Erro ao converter data: %v", err)
	} else {
		logrus.Infof("📅 Data convertida: %s", sqlDate)
	}

	// ValidationUtils
	validationUtils := NewValidationUtils()
	logrus.Infof("🔍 CPF válido: %t", validationUtils.IsValidCPF("123.456.789-09"))
	logrus.Infof("🔍 Email válido: %t", validationUtils.IsValidEmail("test@example.com"))

	// SQLUtils
	sqlUtils := NewSQLUtils()
	logrus.Infof("💾 QuoteString: %s", sqlUtils.QuoteString("O'Connor"))
	logrus.Infof("💾 FormatSQLValue: %s", sqlUtils.FormatSQLValue(123.45))

	logrus.Info("✅ Exemplo de utilitários concluído!")
}

// ExemploTestes demonstra como executar testes
func ExemploTestes() {
	logrus.Info("🧪 Exemplo de execução de testes")

	// Criar dados de teste
	testData := &DarmData{
		Inscricao:      "123456",
		CodigoReceita:  "2623",
		ValorPrincipal: "1.234,56",
		ValorTotal:     "1.234,56",
		DataVencimento: "15/12/2024",
		Exercicio:      "2025",
		NumeroGuia:     "123456789",
	}

	// Testar extração de dados
	processor := NewDarmProcessor()
	testText := `
	02. INSCRIÇÃO MUNICIPAL 123456
	01. RECEITA 262-3
	06. VALOR DO TRIBUTO R$ 1.234,56
	09. VALOR TOTAL R$ 1.234,56
	03. DATA VENCIMENTO 15/12/2024
	04. ANO DE REFERÊNCIA 2025
	05. GUIA NØ
	123456789
	`

	extractedData := processor.extractDarmData(testText)
	if extractedData != nil {
		logrus.Infof("✅ Dados extraídos: Inscrição=%s, Guia=%s",
			extractedData.Inscricao, extractedData.NumeroGuia)
	} else {
		logrus.Error("❌ Falha na extração de dados")
	}

	// Testar geração de SQL
	sql := processor.generateSQLInsert(testData)
	if sql != "" {
		logrus.Info("✅ SQL gerado com sucesso")
		logrus.Debugf("SQL: %s", sql[:100]+"...")
	} else {
		logrus.Error("❌ Falha na geração de SQL")
	}

	logrus.Info("✅ Exemplo de testes concluído!")
}

// ExemploPerformance demonstra benchmarks
func ExemploPerformance() {
	logrus.Info("⚡ Exemplo de performance")

	processor := NewDarmProcessor()
	testText := `
	02. INSCRIÇÃO MUNICIPAL 123456
	01. RECEITA 262-3
	06. VALOR DO TRIBUTO R$ 1.234,56
	09. VALOR TOTAL R$ 1.234,56
	03. DATA VENCIMENTO 15/12/2024
	04. ANO DE REFERÊNCIA 2025
	05. GUIA NØ
	123456789
	`

	// Benchmark de extração de dados
	logrus.Info("🔄 Executando benchmark de extração...")
	for i := 0; i < 1000; i++ {
		processor.extractDarmData(testText)
	}

	// Benchmark de parse de valores monetários
	logrus.Info("🔄 Executando benchmark de parse monetário...")
	for i := 0; i < 10000; i++ {
		processor.parseMonetaryValue("R$ 1.234,56")
	}

	logrus.Info("✅ Exemplo de performance concluído!")
}

// Funções auxiliares
func isPDFFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".pdf" || ext == ".PDF"
}

func showResults(processor *DarmProcessor) {
	logrus.Infof("📊 Total de guias processadas: %d", len(processor.GuiasProcessadas))

	if len(processor.GuiasProcessadas) > 0 {
		logrus.Info("📋 Guias processadas:")
		for i, guia := range processor.GuiasProcessadas {
			logrus.Infof("  %d. Guia %s", i+1, guia)
		}
	}

	// Verificar arquivos gerados
	outputFiles, err := os.ReadDir(processor.OutputDir)
	if err != nil {
		logrus.Errorf("❌ Erro ao ler diretório de saída: %v", err)
		return
	}

	logrus.Info("📄 Arquivos gerados:")
	for _, file := range outputFiles {
		if !file.IsDir() {
			logrus.Infof("  - %s", file.Name())
		}
	}
}
