package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/pressly/goose"
	"log"
	"os"
	"transaction-server/app"
	"transaction-server/app/boot"
	_ "transaction-server/internal/database/migrations"
)

var (
	flags   = flag.NewFlagSet("goose", flag.ExitOnError)
	dir     = flags.String("dir", "internal/database/migrations", "Directory with migration files")
	verbose = flags.Bool("v", false, "Enable verbose mode")
)

func main() {

	flags.Usage = usage
	if err := flags.Parse(os.Args[1:]); err != nil {
		log.Fatalf("error parsing the flags: %v", err)
	}
	args := flags.Args()
	if *verbose {
		goose.SetVerbose(true)
	}

	// I.e. no command provided, hence print usage and return.
	if len(args) < 1 {
		flags.Usage()
		return
	}

	// Prepares command and arguments for goose's run.
	command := args[0]
	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	// If command is create or fix, no need to connect to db and hence the
	// specific case handling.
	switch command {
	case "create":
		if err := goose.Run("create", nil, *dir, arguments...); err != nil {
			log.Fatalf("failed to run command: %v", err)
		}
		return
	case "fix":
		if err := goose.Run("fix", nil, *dir); err != nil {
			log.Fatalf("failed to run the command: %v", err)
		}
		return
	}

	// For other commands boot application (hence getting db and config ready).
	// Read application's dialect and get sqldb instance.
	if err := boot.Initialize(context.Background()); err != nil {
		log.Fatalf("failed to run command: %v", err)
	}

	dialect := app.Context().Config().Db.Dialect
	if err := goose.SetDialect(dialect); err != nil {
		log.Fatalf("failed to run command: %v", err)
	}
	sqlDb, errDb := app.Context().DB().Instance(context.Background()).DB()
	if errDb != nil {
		log.Fatalf("failed to run command: %v", errDb)
	}

	dirs := []string{*dir}

	for _, dir := range dirs {
		// Finally, executes the goose's command.
		if err := goose.Run(command, sqlDb, dir, arguments...); err != nil {
			log.Fatalf("failed to run command: %v", err)
		}
	}
}

func usage() {
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}

var usageCommands = `
Commands:
	up                   Migrate the DB to the most recent version available
	up-to VERSION        Migrate the DB to a specific VERSION
	down                 Roll back the version by 1
	down-to VERSION      Roll back to a specific VERSION
	redo                 Re-run the latest migration
	reset                Roll back all migrations
	status               Dump the migration status for the current DB
	version              Print the current version of the database
	create NAME          Creates new migration file with the current timestamp
	fix                  Apply sequential ordering to migrations
`
