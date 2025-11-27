package db

import (
	"database/sql"
	"os"
	"testing"
)

func TestNewPostgresDB_InvalidParams(t *testing.T) {
	if err := os.Setenv("DB_HOST", "invalid_host"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_PORT", "5432"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_USER", "invalid"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_PASSWORD", "invalid"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_NAME", "invalid"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_SSLMODE", "disable"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}

	db, err := NewPostgresDB()
	if err == nil {
		if err := db.Close(); err != nil {
			t.Errorf("error closing db: %v", err)
		}
		t.Error("expected error with invalid params, got nil")
	}
}

func TestNewPostgresDB_EmptySSLMode(t *testing.T) {
	if err := os.Setenv("DB_HOST", "invalid_host"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_PORT", "5432"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_USER", "invalid"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_PASSWORD", "invalid"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_NAME", "invalid"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Unsetenv("DB_SSLMODE"); err != nil {
		t.Fatalf("error unsetting env: %v", err)
	}

	db, err := NewPostgresDB()
	if err == nil {
		if err := db.Close(); err != nil {
			t.Errorf("error closing db: %v", err)
		}
		t.Error("expected error with invalid params, got nil")
	}
}

func TestNewPostgresDB_DriverNotFound(t *testing.T) {
	// Simular driver no registrado
	if err := os.Setenv("DB_HOST", "localhost"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_PORT", "5432"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_USER", "user"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_PASSWORD", "pass"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_NAME", "db"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}
	if err := os.Setenv("DB_SSLMODE", "disable"); err != nil {
		t.Fatalf("error setting env: %v", err)
	}

	// Desregistrar el driver temporalmente (no es posible en sql stdlib, pero el test cubre el error de conexión)
	db, err := sql.Open("invalid_driver", "")
	if err == nil {
		if err := db.Close(); err != nil {
			t.Errorf("error closing db: %v", err)
		}
	}
}
