package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// FileUtils contém utilitários para manipulação de arquivos
type FileUtils struct{}

// NewFileUtils cria nova instância de FileUtils
func NewFileUtils() *FileUtils {
	return &FileUtils{}
}

// EnsureDir cria diretório se não existir
func (fu *FileUtils) EnsureDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}

// FileExists verifica se arquivo existe
func (fu *FileUtils) FileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

// GetFileSize retorna tamanho do arquivo
func (fu *FileUtils) GetFileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// CopyFile copia arquivo
func (fu *FileUtils) CopyFile(src, dst string) error {
	// Ler arquivo fonte
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}

	// Criar diretório de destino se não existir
	dir := filepath.Dir(dst)
	if err := fu.EnsureDir(dir); err != nil {
		return err
	}

	// Escrever arquivo destino
	return os.WriteFile(dst, data, 0644)
}

// StringUtils contém utilitários para manipulação de strings
type StringUtils struct{}

// NewStringUtils cria nova instância de StringUtils
func NewStringUtils() *StringUtils {
	return &StringUtils{}
}

// RemoveAccents remove acentos de string
func (su *StringUtils) RemoveAccents(s string) string {
	// Mapeamento de caracteres acentuados para não acentuados
	accents := map[rune]rune{
		'á': 'a', 'à': 'a', 'ã': 'a', 'â': 'a', 'ä': 'a',
		'é': 'e', 'è': 'e', 'ê': 'e', 'ë': 'e',
		'í': 'i', 'ì': 'i', 'î': 'i', 'ï': 'i',
		'ó': 'o', 'ò': 'o', 'õ': 'o', 'ô': 'o', 'ö': 'o',
		'ú': 'u', 'ù': 'u', 'û': 'u', 'ü': 'u',
		'ç': 'c',
		'Á': 'A', 'À': 'A', 'Ã': 'A', 'Â': 'A', 'Ä': 'A',
		'É': 'E', 'È': 'E', 'Ê': 'E', 'Ë': 'E',
		'Í': 'I', 'Ì': 'I', 'Î': 'I', 'Ï': 'I',
		'Ó': 'O', 'Ò': 'O', 'Õ': 'O', 'Ô': 'O', 'Ö': 'O',
		'Ú': 'U', 'Ù': 'U', 'Û': 'U', 'Ü': 'U',
		'Ç': 'C',
	}

	result := make([]rune, len(s))
	for i, r := range s {
		if replacement, exists := accents[r]; exists {
			result[i] = replacement
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// CleanString limpa string removendo caracteres especiais
func (su *StringUtils) CleanString(s string) string {
	// Remover acentos
	s = su.RemoveAccents(s)

	// Manter apenas letras, números, espaços e alguns caracteres especiais
	re := regexp.MustCompile(`[^a-zA-Z0-9\s\-_\.\/]`)
	return re.ReplaceAllString(s, "")
}

// FormatCurrency formata valor como moeda
func (su *StringUtils) FormatCurrency(value float64) string {
	return fmt.Sprintf("%.2f", value)
}

// ParseCurrency parseia string de moeda para float
func (su *StringUtils) ParseCurrency(s string) (float64, error) {
	// Remover R$, espaços e pontos de milhares
	clean := regexp.MustCompile(`[R$\s]`).ReplaceAllString(s, "")

	// Se tem vírgula, tratar como separador decimal brasileiro
	if strings.Contains(clean, ",") {
		if strings.Contains(clean, ".") {
			// Formato brasileiro: 9.014,06
			clean = strings.ReplaceAll(clean, ".", "")
			clean = strings.ReplaceAll(clean, ",", ".")
		} else {
			// Só vírgula: 9014,06
			clean = strings.ReplaceAll(clean, ",", ".")
		}
	}

	return strconv.ParseFloat(clean, 64)
}

// DateUtils contém utilitários para manipulação de datas
type DateUtils struct{}

// NewDateUtils cria nova instância de DateUtils
func NewDateUtils() *DateUtils {
	return &DateUtils{}
}

// ParseDateBR parseia data no formato brasileiro (DD/MM/YYYY)
func (du *DateUtils) ParseDateBR(dateStr string) (time.Time, error) {
	return time.Parse("02/01/2006", dateStr)
}

// FormatDateBR formata data para formato brasileiro
func (du *DateUtils) FormatDateBR(date time.Time) string {
	return date.Format("02/01/2006")
}

// ParseDateSQL parseia data para formato SQL (YYYY-MM-DD)
func (du *DateUtils) ParseDateSQL(dateStr string) (time.Time, error) {
	return time.Parse("2006-01-02", dateStr)
}

// FormatDateSQL formata data para formato SQL
func (du *DateUtils) FormatDateSQL(date time.Time) string {
	return date.Format("2006-01-02")
}

// ConvertDateBRToSQL converte data BR para SQL
func (du *DateUtils) ConvertDateBRToSQL(dateBR string) (string, error) {
	date, err := du.ParseDateBR(dateBR)
	if err != nil {
		return "", err
	}
	return du.FormatDateSQL(date), nil
}

// ValidationUtils contém utilitários para validação
type ValidationUtils struct{}

// NewValidationUtils cria nova instância de ValidationUtils
func NewValidationUtils() *ValidationUtils {
	return &ValidationUtils{}
}

// IsValidCPF valida CPF
func (vu *ValidationUtils) IsValidCPF(cpf string) bool {
	// Remover caracteres não numéricos
	cpf = regexp.MustCompile(`\D`).ReplaceAllString(cpf, "")

	if len(cpf) != 11 {
		return false
	}

	// Verificar se todos os dígitos são iguais
	if cpf == "00000000000" || cpf == "11111111111" || cpf == "22222222222" ||
		cpf == "33333333333" || cpf == "44444444444" || cpf == "55555555555" ||
		cpf == "66666666666" || cpf == "77777777777" || cpf == "88888888888" ||
		cpf == "99999999999" {
		return false
	}

	// Calcular primeiro dígito verificador
	sum := 0
	for i := 0; i < 9; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (10 - i)
	}
	remainder := sum % 11
	if remainder < 2 {
		if string(cpf[9]) != "0" {
			return false
		}
	} else {
		expected := 11 - remainder
		actual, _ := strconv.Atoi(string(cpf[9]))
		if actual != expected {
			return false
		}
	}

	// Calcular segundo dígito verificador
	sum = 0
	for i := 0; i < 10; i++ {
		digit, _ := strconv.Atoi(string(cpf[i]))
		sum += digit * (11 - i)
	}
	remainder = sum % 11
	if remainder < 2 {
		if string(cpf[10]) != "0" {
			return false
		}
	} else {
		expected := 11 - remainder
		actual, _ := strconv.Atoi(string(cpf[10]))
		if actual != expected {
			return false
		}
	}

	return true
}

// IsValidCNPJ valida CNPJ
func (vu *ValidationUtils) IsValidCNPJ(cnpj string) bool {
	// Remover caracteres não numéricos
	cnpj = regexp.MustCompile(`\D`).ReplaceAllString(cnpj, "")

	if len(cnpj) != 14 {
		return false
	}

	// Verificar se todos os dígitos são iguais
	if cnpj == "00000000000000" || cnpj == "11111111111111" || cnpj == "22222222222222" ||
		cnpj == "33333333333333" || cnpj == "44444444444444" || cnpj == "55555555555555" ||
		cnpj == "66666666666666" || cnpj == "77777777777777" || cnpj == "88888888888888" ||
		cnpj == "99999999999999" {
		return false
	}

	// Calcular primeiro dígito verificador
	weights := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum := 0
	for i := 0; i < 12; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights[i]
	}
	remainder := sum % 11
	if remainder < 2 {
		if string(cnpj[12]) != "0" {
			return false
		}
	} else {
		expected := 11 - remainder
		actual, _ := strconv.Atoi(string(cnpj[12]))
		if actual != expected {
			return false
		}
	}

	// Calcular segundo dígito verificador
	weights = []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	sum = 0
	for i := 0; i < 13; i++ {
		digit, _ := strconv.Atoi(string(cnpj[i]))
		sum += digit * weights[i]
	}
	remainder = sum % 11
	if remainder < 2 {
		if string(cnpj[13]) != "0" {
			return false
		}
	} else {
		expected := 11 - remainder
		actual, _ := strconv.Atoi(string(cnpj[13]))
		if actual != expected {
			return false
		}
	}

	return true
}

// IsValidEmail valida email
func (vu *ValidationUtils) IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// IsValidDate valida data no formato DD/MM/YYYY
func (vu *ValidationUtils) IsValidDate(dateStr string) bool {
	_, err := time.Parse("02/01/2006", dateStr)
	return err == nil
}

// SQLUtils contém utilitários para SQL
type SQLUtils struct{}

// NewSQLUtils cria nova instância de SQLUtils
func NewSQLUtils() *SQLUtils {
	return &SQLUtils{}
}

// EscapeString escapa string para SQL
func (su *SQLUtils) EscapeString(s string) string {
	// Substituir aspas simples por duas aspas simples
	return strings.ReplaceAll(s, "'", "''")
}

// QuoteString coloca string entre aspas simples
func (su *SQLUtils) QuoteString(s string) string {
	return "'" + su.EscapeString(s) + "'"
}

// FormatSQLValue formata valor para SQL
func (su *SQLUtils) FormatSQLValue(value interface{}) string {
	switch v := value.(type) {
	case string:
		if v == "" {
			return "NULL"
		}
		return su.QuoteString(v)
	case int, int32, int64:
		return fmt.Sprintf("%v", v)
	case float32, float64:
		return fmt.Sprintf("%.2f", v)
	case bool:
		if v {
			return "1"
		}
		return "0"
	case nil:
		return "NULL"
	default:
		return su.QuoteString(fmt.Sprintf("%v", v))
	}
}

// GeneratePlaceholders gera placeholders para SQL
func (su *SQLUtils) GeneratePlaceholders(count int) string {
	placeholders := make([]string, count)
	for i := 0; i < count; i++ {
		placeholders[i] = "?"
	}
	return strings.Join(placeholders, ", ")
}
