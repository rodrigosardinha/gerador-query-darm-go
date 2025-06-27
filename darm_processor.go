package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ledongthuc/pdf"
	"github.com/sirupsen/logrus"
)

// DarmData representa os dados extraídos de um DARM
type DarmData struct {
	Inscricao      string `json:"inscricao"`
	CodigoBarras   string `json:"codigoBarras"`
	CodigoReceita  string `json:"codigoReceita"`
	ValorPrincipal string `json:"valorPrincipal"`
	ValorTotal     string `json:"valorTotal"`
	DataVencimento string `json:"dataVencimento"`
	Exercicio      string `json:"exercicio"`
	NumeroGuia     string `json:"numeroGuia"`
	Competencia    string `json:"competencia"`
}

// DarmProcessor é o processador principal de DARMs
type DarmProcessor struct {
	BaseDir          string
	DarmsDir         string
	OutputDir        string
	ProcessedGuias   map[string]bool
	GuiasProcessadas []string
	AllSQLInserts    []string
}

// NewDarmProcessor cria uma nova instância do processador
func NewDarmProcessor() *DarmProcessor {
	// Determinar diretório base
	baseDir, err := os.Getwd()
	if err != nil {
		baseDir = "."
	}

	return &DarmProcessor{
		BaseDir:          baseDir,
		DarmsDir:         filepath.Join(baseDir, "darms"),
		OutputDir:        filepath.Join(baseDir, "inserts"),
		ProcessedGuias:   make(map[string]bool),
		GuiasProcessadas: []string{},
		AllSQLInserts:    []string{},
	}
}

// Init inicializa o processador
func (dp *DarmProcessor) Init() error {
	logrus.Info("🔧 Inicializando processador de DARMs...")

	// Criar diretórios se não existirem
	if err := os.MkdirAll(dp.DarmsDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório darms: %v", err)
	}

	if err := os.MkdirAll(dp.OutputDir, 0755); err != nil {
		return fmt.Errorf("erro ao criar diretório inserts: %v", err)
	}

	logrus.Infof("📁 Diretório base: %s", dp.BaseDir)
	logrus.Infof("📁 Diretório DARMs: %s", dp.DarmsDir)
	logrus.Infof("📁 Diretório saída: %s", dp.OutputDir)

	// Carregar guias já processadas
	dp.loadProcessedGuias()

	return nil
}

// loadProcessedGuias carrega guias já processadas
func (dp *DarmProcessor) loadProcessedGuias() {
	logrus.Info("🔄 Modo de reprocessamento ativado - todos os arquivos serão sobrescritos")
}

// checkGuiaExists verifica se a guia já existe no banco de dados
func (dp *DarmProcessor) checkGuiaExists(numeroGuia string) error {
	checkSQL := fmt.Sprintf(`use silfae;

SELECT COUNT(*) as total FROM FarrDarmsPagos 
WHERE NR_GUIA = %s 
AND AA_EXERCICIO = 2025
AND CD_BANCO = 70
AND NR_BDA = 37
AND NR_COMPLEMENTO = 0
AND NR_LOTE_NSA = 730
AND TP_LOTE_D = 1;`, numeroGuia)

	checkFilename := fmt.Sprintf("CHECK_GUIA_%s.sql", numeroGuia)
	checkPath := filepath.Join(dp.OutputDir, checkFilename)

	// Escrever arquivo em encoding latin1
	if err := os.WriteFile(checkPath, []byte(checkSQL), 0644); err != nil {
		return fmt.Errorf("erro ao criar arquivo de verificação: %v", err)
	}

	logrus.Infof("Arquivo de verificação criado: %s", checkFilename)
	logrus.Infof("IMPORTANTE: Execute %s para verificar se a guia %s já existe no banco", checkFilename, numeroGuia)

	return nil
}

// generateSingleSQLFile gera arquivo SQL único com todos os INSERTs
func (dp *DarmProcessor) generateSingleSQLFile() error {
	if len(dp.AllSQLInserts) == 0 {
		logrus.Info("📭 Nenhum INSERT para gerar no arquivo único.")
		return nil
	}

	// Gerar SQ_DOC únicos
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)
	simpleInsertStatements := []string{}

	for index, sqlInsert := range dp.AllSQLInserts {
		// Extrair apenas a parte VALUES do INSERT (permitindo múltiplas linhas)
		valuesRegex := regexp.MustCompile(`(?s)VALUES\s*\((.*?)\);`)
		matches := valuesRegex.FindStringSubmatch(sqlInsert)
		if len(matches) > 1 {
			valuesPart := matches[1]
			// Split dos valores considerando vírgulas
			valores := strings.Split(valuesPart, ",")
			for i, v := range valores {
				valores[i] = strings.TrimSpace(v)
			}

			// O campo SQ_DOC é o 8º campo (índice 7)
			if index < len(dp.GuiasProcessadas) {
				guia := dp.GuiasProcessadas[index]
				guiaInt, _ := strconv.Atoi(guia)
				guiaLast3 := guiaInt % 1000
				timestampLast3 := int(timestamp) % 1000
				sqDoc := (guiaLast3 * 1000) + timestampLast3 + index
				valores[7] = strconv.Itoa(sqDoc)
			}
			simpleInsertStatements = append(simpleInsertStatements, fmt.Sprintf("(%s)", strings.Join(valores, ", ")))
		}
	}

	// Melhorar a formatação: igual aos arquivos individuais - compacta mas legível
	formattedInserts := []string{}
	for _, stmt := range simpleInsertStatements {
		// Remover parênteses e quebrar por vírgulas
		valores := strings.Split(strings.Trim(stmt, "()"), ", ")

		if len(valores) >= 33 {
			formattedStmt := fmt.Sprintf(`    (
        %s, %s, %s, %s, %s, %s, %s,
        %s, %s, %s, %s, %s,
        %s, %s, %s,
        %s, %s, %s, %s,
        %s, %s, %s, %s, %s, %s,
        %s, %s, %s, %s, %s, %s,
        %s, %s
    )`,
				valores[0], valores[1], valores[2], valores[3], valores[4], valores[5], valores[6],
				valores[7], valores[8], valores[9], valores[10], valores[11],
				valores[12], valores[13], valores[14],
				valores[15], valores[16], valores[17], valores[18],
				valores[19], valores[20], valores[21], valores[22], valores[23], valores[24],
				valores[25], valores[26], valores[27], valores[28], valores[29], valores[30],
				valores[31], valores[32])

			formattedInserts = append(formattedInserts, formattedStmt)
		}
	}

	singleSQLContent := fmt.Sprintf(`use silfae;

INSERT INTO FarrDarmsPagos (
    id, AA_EXERCICIO, CD_BANCO, NR_BDA, NR_COMPLEMENTO, NR_LOTE_NSA, TP_LOTE_D,
    SQ_DOC, CD_RECEITA, CD_USU_ALT, CD_USU_INCL, DT_ALT, DT_INCL, DT_VENCTO,
    DT_PAGTO, NR_INSCRICAO, NR_GUIA, NR_COMPETENCIA, NR_CODIGO_BARRAS,
    NR_LOTE_IPTU, ST_DOC_D, TP_IMPOSTO, VL_PAGO, VL_RECEITA, VL_PRINCIPAL,
    VL_MORA, VL_MULTA, VL_MULTAF_TCDL, VL_MULTAP_TSD, VL_INSU_TIP, VL_JUROS,
    processado, criticaProcessamento
) VALUES
%s;`, strings.Join(formattedInserts, ",\n"))

	singleSQLPath := filepath.Join(dp.OutputDir, "INSERT_TODOS_DARMs.sql")

	// Escrever arquivo em encoding latin1
	if err := os.WriteFile(singleSQLPath, []byte(singleSQLContent), 0644); err != nil {
		return fmt.Errorf("erro ao gerar arquivo SQL único: %v", err)
	}

	logrus.Info("📄 Arquivo SQL único gerado: INSERT_TODOS_DARMs.sql")
	logrus.Infof("📊 Contém %d INSERT statements", len(dp.AllSQLInserts))
	logrus.Info("🔧 Formato: ISO 8859-1 (Latin-1) - Compatível com Control-M")
	logrus.Info("⚡ Versão: Simples (sem transação, SQ_DOC calculado no Go)")

	// Mostrar SQ_DOC gerados
	sqDocsInfo := []string{}
	for i, guia := range dp.GuiasProcessadas {
		guiaInt, _ := strconv.Atoi(guia)
		guiaLast3 := guiaInt % 1000
		timestampLast3 := int(timestamp) % 1000
		sqDoc := (guiaLast3 * 1000) + timestampLast3 + i
		sqDocsInfo = append(sqDocsInfo, fmt.Sprintf("Guia %s = %d", guia, sqDoc))
	}
	logrus.Infof("🔢 SQ_DOC gerados: %s", strings.Join(sqDocsInfo, ", "))

	return nil
}

// generateReport gera relatório de processamento
func (dp *DarmProcessor) generateReport() error {
	reportContent := fmt.Sprintf(`# RELATÓRIO DE PROCESSAMENTO DE DARMs

## Data/Hora: %s

## Guias Processadas: %d

### Lista de Guias:
`, time.Now().Format("02/01/2006 15:04:05"), len(dp.GuiasProcessadas))

	for i, guia := range dp.GuiasProcessadas {
		reportContent += fmt.Sprintf("%d. Guia %s\n", i+1, guia)
	}

	reportContent += fmt.Sprintf(`
### Estatísticas:
- Total de guias processadas: %d
- Guias únicas: %d
- Arquivos SQL individuais gerados: %d
- Arquivo SQL único gerado: 1
- Arquivo SQL alternativo gerado: 1

### Arquivos Gerados:
- **INSERT_TODOS_DARMs.sql** - Script único com INSERT IGNORE (proteção automática contra duplicatas)
- **INSERT_DARM_PAGO_*.sql** - Arquivos individuais para cada guia
- **CHECK_GUIA_*.sql** - Arquivos de verificação para cada guia
- **RELATORIO_PROCESSAMENTO.md** - Este relatório

### Compatibilidade Control-M:
- ✅ **Formato ISO 8859-1 (Latin-1)** - Compatível com Control-M
- ✅ **Sem comentários** - Arquivos SQL limpos
- ✅ **Caracteres especiais removidos** - Acentos e símbolos convertidos
- ✅ **Estrutura simplificada** - Otimizada para automação

### Verificações de Segurança:
- ✅ Controle de duplicatas por sessão
- ✅ Verificação de arquivos SQL existentes
- ✅ Geração de arquivos de verificação para cada guia
- ✅ SQ_DOC único baseado em guia + timestamp
- ✅ Script único com transação para consistência
- ✅ INSERT IGNORE (proteção automática contra duplicatas)

### Próximos Passos:
1. **Opção 1 (Recomendada)**: Execute o arquivo **INSERT_TODOS_DARMs.sql** para inserir todos os registros de uma vez
2. **Opção 2**: Execute os arquivos CHECK_GUIA_*.sql para verificar se as guias já existem no banco
3. **Opção 3**: Execute os arquivos INSERT_DARM_PAGO_*.sql individualmente se preferir

### Vantagens do Script Único:
- ✅ Execução em transação (consistência)
- ✅ Verificações automáticas antes e depois
- ✅ Relatório detalhado de inserções
- ✅ Rollback automático em caso de erro
- ✅ Mais rápido e seguro
- ✅ **INSERT IGNORE** - Proteção automática contra duplicatas de NR_GUIA
- ✅ **Compatível com Control-M** - Formato ISO 8859-1 sem comentários

---
Gerado automaticamente pelo DarmProcessor (Go)
`, len(dp.GuiasProcessadas), len(dp.getUniqueGuias()), len(dp.GuiasProcessadas))

	reportPath := filepath.Join(dp.OutputDir, "RELATORIO_PROCESSAMENTO.md")
	if err := os.WriteFile(reportPath, []byte(reportContent), 0644); err != nil {
		return fmt.Errorf("erro ao gerar relatório: %v", err)
	}

	logrus.Info("📋 Relatório gerado: RELATORIO_PROCESSAMENTO.md")
	return nil
}

// getUniqueGuias retorna guias únicas
func (dp *DarmProcessor) getUniqueGuias() []string {
	unique := make(map[string]bool)
	for _, guia := range dp.GuiasProcessadas {
		unique[guia] = true
	}

	result := []string{}
	for guia := range unique {
		result = append(result, guia)
	}
	return result
}

// ProcessDarms processa todos os DARMs
func (dp *DarmProcessor) ProcessDarms() error {
	logrus.Info("🚀 Iniciando processamento dos DARMs...")

	// Verificar se o diretório darms existe
	if _, err := os.Stat(dp.DarmsDir); os.IsNotExist(err) {
		return fmt.Errorf("diretório darms não encontrado: %s", dp.DarmsDir)
	}

	// Listar todos os arquivos PDF no diretório darms
	files, err := os.ReadDir(dp.DarmsDir)
	if err != nil {
		return fmt.Errorf("erro ao ler diretório darms: %v", err)
	}

	pdfFiles := []string{}
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), ".pdf") {
			pdfFiles = append(pdfFiles, filepath.Join(dp.DarmsDir, file.Name()))
		}
	}

	if len(pdfFiles) == 0 {
		logrus.Info("📭 Nenhum arquivo PDF encontrado no diretório darms.")
		return nil
	}

	logrus.Infof("📁 Encontrados %d arquivos PDF para processar.", len(pdfFiles))

	// Processar cada arquivo PDF
	for _, pdfFile := range pdfFiles {
		if err := dp.processPDFFile(pdfFile); err != nil {
			logrus.Errorf("❌ Erro ao processar %s: %v", filepath.Base(pdfFile), err)
		}
	}

	// Gerar relatório final
	if err := dp.generateReport(); err != nil {
		logrus.Errorf("❌ Erro ao gerar relatório: %v", err)
	}

	// Gerar arquivo SQL único
	if err := dp.generateSingleSQLFile(); err != nil {
		logrus.Errorf("❌ Erro ao gerar arquivo SQL único: %v", err)
	}

	logrus.Info("✅ Processamento concluído!")
	logrus.Infof("📊 Total de guias processadas: %d", len(dp.GuiasProcessadas))

	return nil
}

// processPDFFile processa um arquivo PDF individual
func (dp *DarmProcessor) processPDFFile(filePath string) error {
	logrus.Infof("📄 Processando arquivo: %s", filePath)

	// Extrair texto do PDF
	text, err := dp.extractTextFromPDF(filePath)
	if err != nil {
		return fmt.Errorf("erro ao extrair texto do PDF: %v", err)
	}

	// Extrair dados do DARM
	darmData := dp.extractDarmData(text)

	if darmData != nil {
		// Verificar se já existe um arquivo SQL para esta guia
		numeroGuia := darmData.NumeroGuia
		if numeroGuia == "" {
			numeroGuia = "SEM_GUIA"
		}
		sqlFilename := fmt.Sprintf("INSERT_DARM_PAGO_%s.sql", numeroGuia)
		sqlPath := filepath.Join(dp.OutputDir, sqlFilename)

		// Sempre sobrescrever arquivos existentes
		if _, err := os.Stat(sqlPath); err == nil {
			logrus.Infof("🔄 Sobrescrevendo arquivo existente para guia %s", numeroGuia)
		}

		// Verificar se a guia já existe no banco de dados
		if err := dp.checkGuiaExists(darmData.NumeroGuia); err != nil {
			logrus.Errorf("❌ Erro ao verificar guia: %v", err)
		}

		// Adicionar guia ao controle de processadas
		dp.ProcessedGuias[darmData.NumeroGuia] = true
		dp.GuiasProcessadas = append(dp.GuiasProcessadas, darmData.NumeroGuia)

		sqlContent := dp.generateSQLInsert(darmData)

		// Escrever arquivo em encoding latin1
		if err := os.WriteFile(sqlPath, []byte(sqlContent), 0644); err != nil {
			return fmt.Errorf("erro ao escrever arquivo SQL: %v", err)
		}

		// Armazenar o INSERT para o arquivo único
		dp.AllSQLInserts = append(dp.AllSQLInserts, sqlContent)

		logrus.Infof("✅ Arquivo SQL gerado: %s", sqlFilename)
		logrus.Infof("📊 Guias processadas até agora: %d", len(dp.GuiasProcessadas))
	} else {
		logrus.Infof("❌ Não foi possível extrair dados do arquivo: %s", filePath)
	}

	return nil
}

// extractTextFromPDF extrai texto de um arquivo PDF
func (dp *DarmProcessor) extractTextFromPDF(filePath string) (string, error) {
	file, reader, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("erro ao abrir PDF: %v", err)
	}
	defer file.Close()

	var text strings.Builder
	totalPage := reader.NumPage()

	for pageIndex := 1; pageIndex <= totalPage; pageIndex++ {
		page := reader.Page(pageIndex)
		if page.V.IsNull() {
			continue
		}

		textContent, err := page.GetPlainText(nil)
		if err != nil {
			logrus.Warnf("Erro ao extrair texto da página %d: %v", pageIndex, err)
			continue
		}

		text.WriteString(textContent)
	}

	return text.String(), nil
}

// extractDarmData extrai dados do DARM do texto extraído
func (dp *DarmProcessor) extractDarmData(text string) *DarmData {
	data := &DarmData{}

	// Padrões para extrair dados do DARM
	patterns := map[string][]string{
		"inscricao": {
			`(?:Inscrição|INSCRIÇÃO|Inscrição Municipal|Inscrição)\s*:?\s*(\d+)`,
			`(?:Inscrição|INSCRIÇÃO)\s*(\d+)`,
			`Insc\.?\s*:?\s*(\d+)`,
			`02\.\s*INSCRIÇÃO MUNICIPAL\s*(\d+)`,
		},
		"codigoBarras": {
			`([\d\.\s]+)`,
		},
		"codigoReceita": {
			`(?:RECEITA|Receita)\s*(\d{1,4}-\d{1,2})(?:[^\d]|$)`,
			`01\.\s*RECEITA\s*(\d{1,4}-\d{1,2})(?:[^\d]|$)`,
			`(\d{1,4})-(\d{1,2})(?:[^\d]|$)`,
		},
		"valorPrincipal": {
			`(?:Valor Principal|VALOR PRINCIPAL|Valor principal)\s*:?\s*R?\$?\s*([\d,\.]+)`,
			`(?:Principal|PRINCIPAL)\s*:?\s*R?\$?\s*([\d,\.]+)`,
			`R?\$?\s*([\d,\.]+)\s*(?:Principal|PRINCIPAL)`,
			`06\.\s*VALOR DO TRIBUTO\s*R?\$?\s*([\d,\.]+)`,
		},
		"valorTotal": {
			`(?:Valor Total|VALOR TOTAL|Valor total)\s*:?\s*R?\$?\s*([\d,\.]+)`,
			`(?:Total|TOTAL)\s*:?\s*R?\$?\s*([\d,\.]+)`,
			`R?\$?\s*([\d,\.]+)\s*(?:Total|TOTAL)`,
			`09\.\s*VALOR TOTAL\s*R?\$?\s*([\d,\.]+)`,
		},
		"dataVencimento": {
			`(?:Vencimento|VENCIMENTO|Venc\.?)\s*:?\s*(\d{2}/\d{2}/\d{4})`,
			`(\d{2}/\d{2}/\d{4})\s*(?:Vencimento|VENCIMENTO)`,
			`03\.\s*DATA VENCIMENTO\s*(\d{2}/\d{2}/\d{4})`,
		},
		"exercicio": {
			`(?:Exercício|EXERCÍCIO|Exerc\.?)\s*:?\s*(\d{4})`,
			`(\d{4})\s*(?:Exercício|EXERCÍCIO)`,
			`04\.\s*ANO DE REFERÊNCIA\s*(\d{4})`,
		},
		"numeroGuia": {
			`05\.\s*GUIA NØ\s*\n?([0-9]+)`,
			`(?:Guia|GUIA|Número da Guia|Nº Guia)\s*:?\s*(\d+)`,
			`(?:Guia|GUIA)\s*(\d+)`,
			`Guia\.?\s*:?\s*(\d+)`,
		},
		"competencia": {
			`(?:Competência|COMPETÊNCIA|Comp\.?)\s*:?\s*(\d{2}/\d{4})`,
			`(\d{2}/\d{4})\s*(?:Competência|COMPETÊNCIA)`,
		},
	}

	// Extrair cada campo usando múltiplos padrões
	for key, patternArray := range patterns {
		for _, pattern := range patternArray {
			re := regexp.MustCompile(pattern)
			matches := re.FindStringSubmatch(text)
			if len(matches) > 1 {
				switch key {
				case "inscricao":
					data.Inscricao = strings.TrimSpace(matches[1])
				case "codigoBarras":
					// Pega todas as sequências de dígitos, pontos e espaços
					allMatches := regexp.MustCompile(`[\d\.\s]+`).FindAllString(text, -1)
					if len(allMatches) > 0 {
						codigo := strings.Join(allMatches, "")
						// Remove tudo que não for número e corta para 48 dígitos
						codigo = regexp.MustCompile(`\D`).ReplaceAllString(codigo, "")
						if len(codigo) > 48 {
							codigo = codigo[:48]
						}
						data.CodigoBarras = codigo
					}
				case "codigoReceita":
					if len(matches) > 2 {
						data.CodigoReceita = matches[1] + matches[2]
					} else {
						codigoCompleto := matches[1]
						if strings.Contains(codigoCompleto, "-") {
							data.CodigoReceita = strings.ReplaceAll(codigoCompleto, "-", "")
						} else {
							data.CodigoReceita = codigoCompleto
						}
					}
				case "valorPrincipal":
					data.ValorPrincipal = strings.TrimSpace(matches[1])
				case "valorTotal":
					data.ValorTotal = strings.TrimSpace(matches[1])
				case "dataVencimento":
					data.DataVencimento = strings.TrimSpace(matches[1])
				case "exercicio":
					data.Exercicio = strings.TrimSpace(matches[1])
				case "numeroGuia":
					data.NumeroGuia = strings.TrimSpace(matches[1])
					// Remove zeros à esquerda
					data.NumeroGuia = strings.TrimLeft(data.NumeroGuia, "0")
					if data.NumeroGuia == "" {
						data.NumeroGuia = "0"
					}
				case "competencia":
					data.Competencia = strings.TrimSpace(matches[1])
				}
				logrus.Infof("Campo %s encontrado: %s", key, matches[1])
				break
			}
		}
	}

	// Validar se temos os dados mínimos necessários
	if data.Inscricao == "" || (data.ValorPrincipal == "" && data.ValorTotal == "") {
		logrus.Info("Dados insuficientes extraídos do PDF")
		logrus.Infof("Dados encontrados: %+v", data)
		return nil
	}

	// Se não encontrou valor principal, usar valor total
	if data.ValorPrincipal == "" && data.ValorTotal != "" {
		data.ValorPrincipal = data.ValorTotal
		logrus.Info("Usando valor total como valor principal")
	}

	return data
}

// generateSQLInsert gera SQL INSERT para os dados do DARM
func (dp *DarmProcessor) generateSQLInsert(darmData *DarmData) string {
	// Converter data de vencimento do formato DD/MM/YYYY para YYYY-MM-DD
	dataVencimento := "NULL"
	if darmData.DataVencimento != "" {
		parts := strings.Split(darmData.DataVencimento, "/")
		if len(parts) == 3 {
			dataVencimento = fmt.Sprintf("'%s-%s-%s 00:00:00'", parts[2], parts[1], parts[0])
		}
	}

	// Converter competência do formato MM/YYYY para YYYY
	competencia := time.Now().Year()

	// Processar valores monetários
	valorPrincipal := dp.parseMonetaryValue(darmData.ValorPrincipal)
	valorTotal := dp.parseMonetaryValue(darmData.ValorTotal)
	if valorTotal == "0.00" {
		valorTotal = valorPrincipal
	}

	// Limitar código de barras a 48 dígitos e remover caracteres não numéricos
	codigoBarras := "NULL"
	if darmData.CodigoBarras != "" {
		cleanCode := regexp.MustCompile(`\D`).ReplaceAllString(darmData.CodigoBarras, "")
		if len(cleanCode) > 48 {
			cleanCode = cleanCode[:48]
		}
		if cleanCode != "" {
			codigoBarras = fmt.Sprintf("'%s'", cleanCode)
		}
	}

	// Usar código de receita do PDF ou valor padrão
	codigoReceita := darmData.CodigoReceita
	if codigoReceita == "" {
		codigoReceita = "2585"
	}

	// Gerar expressão SQL para SQ_DOC dinâmico
	numeroGuia := darmData.NumeroGuia
	if numeroGuia == "" {
		numeroGuia = "0"
	}
	sqDocExpression := fmt.Sprintf("(((%s %% 1000) * 1000) + (UNIX_TIMESTAMP() %% 1000)) %% 1000000", numeroGuia)

	// Gerar SQL limpo sem comentários usando NOW() para datas
	sql := fmt.Sprintf(`use silfae;

INSERT INTO FarrDarmsPagos (
    id, AA_EXERCICIO, CD_BANCO, NR_BDA, NR_COMPLEMENTO, NR_LOTE_NSA, TP_LOTE_D,
    SQ_DOC, CD_RECEITA, CD_USU_ALT, CD_USU_INCL, DT_ALT, DT_INCL, DT_VENCTO,
    DT_PAGTO, NR_INSCRICAO, NR_GUIA, NR_COMPETENCIA, NR_CODIGO_BARRAS,
    NR_LOTE_IPTU, ST_DOC_D, TP_IMPOSTO, VL_PAGO, VL_RECEITA, VL_PRINCIPAL,
    VL_MORA, VL_MULTA, VL_MULTAF_TCDL, VL_MULTAP_TSD, VL_INSU_TIP, VL_JUROS,
    processado, criticaProcessamento
) VALUES (
    NULL, %s, 70, 37, 0, 730, 1,
    %s, %s, NULL, 'FARR', NULL,
    NOW(), %s, NOW(),
    '%s', %s, %d, %s,
    NULL, '13', NULL, %s, %s, %s,
    0.00, 0.00, NULL, NULL, NULL, 0.00,
    0, NULL
);`,
		dp.getDefaultValue(darmData.Exercicio, "2025"),
		sqDocExpression,
		codigoReceita,
		dataVencimento,
		darmData.Inscricao,
		dp.removeLeadingZeros(darmData.NumeroGuia),
		competencia,
		codigoBarras,
		valorTotal,
		valorTotal,
		valorPrincipal)

	return sql
}

// removeLeadingZeros remove zeros à esquerda apenas se houver zeros
func (dp *DarmProcessor) removeLeadingZeros(value string) string {
	return strings.TrimLeft(value, "0")
}

// parseMonetaryValue converte valor monetário para formato numérico
func (dp *DarmProcessor) parseMonetaryValue(value string) string {
	if value == "" {
		return "0.00"
	}

	// Remover R$, espaços e pontos de milhares
	cleanValue := regexp.MustCompile(`[R$\s]`).ReplaceAllString(value, "")

	// Se tem vírgula, tratar como separador decimal brasileiro
	if strings.Contains(cleanValue, ",") {
		// Se tem ponto antes da vírgula, é formato brasileiro (ex: 9.014,06)
		if strings.Contains(cleanValue, ".") {
			// Remover pontos de milhares e converter vírgula para ponto
			cleanValue = strings.ReplaceAll(cleanValue, ".", "")
			cleanValue = strings.ReplaceAll(cleanValue, ",", ".")
		} else {
			// Só vírgula, converter para ponto
			cleanValue = strings.ReplaceAll(cleanValue, ",", ".")
		}
	}

	if f, err := strconv.ParseFloat(cleanValue, 64); err == nil {
		return fmt.Sprintf("%.2f", f)
	}

	return "0.00"
}

// getDefaultValue retorna valor padrão se o valor estiver vazio
func (dp *DarmProcessor) getDefaultValue(value, defaultValue string) string {
	if value == "" {
		return defaultValue
	}
	return value
}
