package rds

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	gofmt "fmt"
	"net/url"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"github.com/gomicro/flow/fmt"

	"github.com/gomicro/blockit/dbblocker"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // needs the postgres driver for connecting
	"github.com/spf13/cobra"
)

var (
	rootConnStr string
	secretName  string
)

type secret struct {
	DBHost string `json:"dbhost"`
	DBPass string `json:"dbpass"`
	DBURL  string `json:"dburl"`
	DBUser string `json:"dbuser"`
}

func init() {
	dbInitCmd.Flags().StringVar(&rootConnStr, "root", "", "database root credentials")
	dbInitCmd.Flags().StringVar(&secretName, "secret", "", "desired name of the created secret to store")
}

var dbInitCmd = &cobra.Command{
	Use:   "init",
	Short: "init a new db",
	Long:  `Initialize a new postgres database and secrets`,
	Args:  cobra.ExactArgs(1),
	Run:   dbInitFunc,
}

func dbInitFunc(cmd *cobra.Command, args []string) {
	ctx := context.Background()

	u, err := url.Parse(rootConnStr)
	if err != nil {
		fmt.Printf("Cannot parse root credentials: %s", err)
		os.Exit(1)
	}

	rootUser := u.User.Username()

	secret := &secret{
		DBHost: u.Host,
		DBUser: strings.ToLower(args[0]),
		DBPass: uuid.New().String(),
	}
	secret.DBURL = gofmt.Sprintf("postgres://%s:%s@%s/%s", secret.DBUser, secret.DBPass, secret.DBHost, secret.DBUser)

	db, err := sql.Open("postgres", rootConnStr)
	if err != nil {
		fmt.Printf("Cannot open postgres connection: %s", err)
		os.Exit(1)
	}

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(5)

	fmt.Verbosef("Waiting for database to connect...")
	block := dbblocker.NewWithContext(ctx, db)
	<-block.Blockit()
	fmt.Verbosef("Connected to database with root credentials")

	dbname := strings.ToLower(secret.DBUser)

	fmt.Verbosef("Checking for database existence")

	exists, err := dbExists(ctx, db, dbname)
	if err != nil {
		fmt.Printf("Cannot determine database existence: %s", err)
		os.Exit(1)
	}

	if !exists {
		fmt.Verbosef("Database does not exist, creating")

		err = createDB(ctx, db, dbname)
		if err != nil {
			fmt.Printf("Cannot create database: %s", err)
			os.Exit(1)
		}
	} else {
		fmt.Verbosef("Database already existed")
	}

	fmt.Verbosef("Checking for user existence")

	exists, err = userExists(ctx, db, secret.DBUser)
	if err != nil {
		fmt.Printf("Cannot determine user existence: %s", err)
		os.Exit(1)
	}

	if !exists {
		fmt.Verbosef("User does not exists, creating")

		err = createUser(ctx, db, secret.DBUser, secret.DBPass)
		if err != nil {
			fmt.Printf("Cannot create user: %s", err)
			os.Exit(1)
		}
	} else {
		fmt.Printf("User already existed, halting")
		os.Exit(1)
	}

	fmt.Verbosef("Granting permissions")

	err = grantPermissions(ctx, db, rootUser, secret.DBUser, dbname)
	if err != nil {
		fmt.Printf("Cannot grant permissions: %s", err)
		os.Exit(1)
	}

	if secretName != "" {
		b, err := json.Marshal(secret)
		if err != nil {
			fmt.Printf("Cannot marshal secret contents: %s", err)
			os.Exit(1)
		}

		bs := string(b)

		input := &secretsmanager.CreateSecretInput{
			Name:         &secretName,
			SecretString: &bs,
		}

		_, err = asmSvc.CreateSecret(input)
		if err != nil {
			fmt.Printf("Error creating secret: %s", err)
		}
	}
}

func dbExists(ctx context.Context, db *sql.DB, dbname string) (bool, error) {
	q := "SELECT datname FROM pg_database WHERE datname = $1"

	var name string
	err := db.QueryRowContext(ctx, q, dbname).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func createDB(ctx context.Context, db *sql.DB, dbname string) error {
	q := gofmt.Sprintf("CREATE DATABASE %s", dbname)

	_, err := db.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
}

func userExists(ctx context.Context, db *sql.DB, user string) (bool, error) {
	q := "SELECT rolname FROM pg_roles WHERE rolname = $1"

	var name string
	err := db.QueryRowContext(ctx, q, user).Scan(&name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func createUser(ctx context.Context, db *sql.DB, user, pass string) error {
	q := gofmt.Sprintf("CREATE USER %s", user)

	_, err := db.ExecContext(ctx, q)
	if err != nil {
		return gofmt.Errorf("failed to create role: %w", err)
	}

	q = gofmt.Sprintf("ALTER USER %s PASSWORD '%s'", user, pass)

	_, err = db.ExecContext(ctx, q)
	if err != nil {
		return gofmt.Errorf("failed to set user password: %s", err)
	}

	return nil
}

func grantPermissions(ctx context.Context, db *sql.DB, rootUser, user, dbname string) error {
	q := gofmt.Sprintf("ALTER DATABASE %s OWNER TO %s", dbname, user)
	_, err := db.ExecContext(ctx, q)
	if err != nil {
		return gofmt.Errorf("failed to assign database ownership: %w", err)
	}

	q = gofmt.Sprintf("GRANT %s to %s", user, rootUser)
	_, err = db.ExecContext(ctx, q)
	if err != nil {
		return gofmt.Errorf("failed to grant root user permissions to db: %w", err)
	}

	q = gofmt.Sprintf("GRANT rds_superuser to %s", user)
	_, err = db.ExecContext(ctx, q)
	if err != nil {
		return gofmt.Errorf("failed to grant user rds permissions: %w", err)
	}

	return nil
}
