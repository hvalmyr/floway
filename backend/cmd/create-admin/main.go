// Command create-admin bootstraps the first admin_users row so that someone
// can log in to /api/v1/admin/login. There is intentionally no HTTP endpoint
// for creating admin accounts (see internal/service/admin_user_service.go).
package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"floway-backend/internal/config"
	"floway-backend/internal/repository"
	"floway-backend/internal/service"
)

func main() {
	login := flag.String("login", "", "admin login (required)")
	password := flag.String("password", "", "admin password (if omitted, read from stdin)")
	flag.Parse()

	if strings.TrimSpace(*login) == "" {
		log.Fatal("-login is required")
	}

	pw := *password
	if pw == "" {
		var err error
		pw, err = readPassword()
		if err != nil {
			log.Fatalf("failed to read password: %v", err)
		}
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	svc := service.NewAdminUserService(repository.NewAdminUserRepository(pool))
	user, err := svc.Create(ctx, *login, pw)
	if err != nil {
		log.Fatalf("failed to create admin user: %v", err)
	}

	fmt.Printf("Created admin user #%d (%s)\n", user.ID, user.Login)
}

func readPassword() (string, error) {
	fmt.Print("Password: ")
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return "", err
	}
	return strings.TrimSpace(line), nil
}
