package performance

import (
	"bytes"
	"os"
	"testing"

	"github.com/sirupsen/logrus"
)

type testHook struct {
	fatalCalled bool
	msg         string
}

func (h *testHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *testHook) Fire(entry *logrus.Entry) error {
	if entry.Level == logrus.FatalLevel {
		h.fatalCalled = true
		h.msg = entry.Message
	}
	return nil
}

func TestProfiling(t *testing.T) {
	// Création d'un logger mock
	logger := logrus.New()
	hook := &testHook{}
	logger.AddHook(hook)
	logger.SetOutput(&bytes.Buffer{}) // éviter d'écrire sur stdout

	// Rediriger os.Create vers des fichiers temporaires pour ne pas polluer le disque
	tmpCPUFile, err := os.CreateTemp("", "cpu-perf-*")
	if err != nil {
		t.Fatalf("unable to create temp CPU file: %v", err)
	}
	defer os.Remove(tmpCPUFile.Name())
	defer tmpCPUFile.Close()

	tmpMemFile, err := os.CreateTemp("", "mem-perf-*")
	if err != nil {
		t.Fatalf("unable to create temp MEM file: %v", err)
	}
	defer os.Remove(tmpMemFile.Name())
	defer tmpMemFile.Close()

	// Patch os.Create temporairement
	origCreate := performanceOsCreate
	performanceOsCreate = func(name string) (*os.File, error) {
		if name == "cpu-perf" {
			return tmpCPUFile, nil
		}
		if name == "mem-perf" {
			return tmpMemFile, nil
		}
		return origCreate(name)
	}
	defer func() { performanceOsCreate = origCreate }()

	// Appel de la fonction
	Profiling(logger)

	// Vérification qu’aucun Fatal n’a été appelé
	if hook.fatalCalled {
		t.Fatalf("logger.Fatal a été appelé: %s", hook.msg)
	}
}
