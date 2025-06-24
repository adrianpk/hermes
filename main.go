package main

import (
	"context"
	"embed"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/adrianpk/hermes/internal/core"
	"github.com/adrianpk/hermes/internal/feat/auth"
	"github.com/adrianpk/hermes/internal/repo/sqlite"
)

const (
	name      = "hermes"
	version   = "v1"
	namespace = "HERMES"
	engine    = "sqlite"
)

var (
	//go:embed assets
	assetsFS embed.FS
)

func main() {
	ctx := context.Background()
	log := am.NewLogger("info")
	cfg := am.LoadCfg(namespace, am.Flags)
	opts := am.DefOpts(log, cfg)

	// FlashManager
	flashManager := am.NewFlashManager()

	// Append WithFlashMiddleware to opts
	opts = append(opts, am.WithFlashMiddleware(flashManager))

	// Create the app with the updated opts
	app := core.NewApp(name, version, assetsFS, opts...)

	queryManager := am.NewQueryManager(assetsFS, engine)
	templateManager := am.NewTemplateManager(assetsFS)

	// Migrator
	migrator := am.NewMigrator(assetsFS, engine)

	// Seeder
	seeder := am.NewSeeder(assetsFS, engine)

	// FileServer
	fileServer := am.NewFileServer(assetsFS)
	app.MountFileServer("/", fileServer)

	// Auth feature
	authRepo := sqlite.NewAuthRepo(queryManager)
	authService := auth.NewService(authRepo)
	authWebHandler := auth.NewWebHandler(templateManager, flashManager, authService)
	authWebRouter := auth.NewWebRouter(authWebHandler)
	authAPIHandler := auth.NewAPIHandler(authService)
	authAPIRouter := auth.NewAPIRouter(authAPIHandler)
	authSeeder := auth.NewSeeder(assetsFS, engine, authRepo)

	app.MountWeb("/auth", authWebRouter)
	app.MountAPI(version, "/auth", authAPIRouter)

	// Add deps
	app.Add(migrator)
	app.Add(seeder)
	app.Add(flashManager)
	app.Add(fileServer)
	app.Add(queryManager)
	app.Add(templateManager)
	app.Add(authRepo)
	app.Add(authService)
	app.Add(authWebHandler)
	app.Add(authAPIHandler)
	app.Add(authWebRouter)
	app.Add(authAPIRouter)
	app.Add(authSeeder)

	err := app.Setup(ctx)
	if err != nil {
		log.Error("Failed to setup the app: ", err)
		return
	}

	// templateManager.Debug()
	// queryManager.Debug()

	err = app.Start(ctx)
	if err != nil {
		log.Error("Failed to start the app: ", err)
	}
}
