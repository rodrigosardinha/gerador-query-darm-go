package main

import (
	"runtime"

	"github.com/sirupsen/logrus"
)

const version = "1.0.0"

func main() {
	// Configurar logging
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
		ForceColors:   true,
	})
	logrus.SetLevel(logrus.InfoLevel)

	logrus.Infof("🚀 Processador de DARMs - Versão Go %s", version)
	logrus.Infof("💻 Sistema: %s/%s", runtime.GOOS, runtime.GOARCH)

	// Criar processador
	processor := NewDarmProcessor()

	// Inicializar
	if err := processor.Init(); err != nil {
		logrus.Fatalf("❌ Erro ao inicializar: %v", err)
	}

	// Processar DARMs
	if err := processor.ProcessDarms(); err != nil {
		logrus.Fatalf("❌ Erro durante o processamento: %v", err)
	}

	logrus.Info("✅ Processamento concluído com sucesso!")
}
