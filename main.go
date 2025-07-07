package main

import (
	"context"
	"embed"

	"github.com/adrianpk/hermes/internal/am"
	"github.com/adrianpk/hermes/internal/core"
	"github.com/adrianpk/hermes/internal/feat/auth"
	"github.com/adrianpk/hermes/internal/feat/ssg"
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

	fm := am.NewFlashManager()

	app := core.NewApp(name, version, assetsFS, opts...)

	queryManager := am.NewQueryManager(assetsFS, engine)
	templateManager := am.NewTemplateManager(assetsFS)

	repo := sqlite.NewHermesRepo(queryManager)
	migrator := am.NewMigrator(assetsFS, engine)
	seeder := am.NewSeeder(assetsFS, engine)
	fileServer := am.NewFileServer(assetsFS)

	app.MountFileServer("/", fileServer)

	// Auth feature
	authService := auth.NewService(repo)
	authWebHandler := auth.NewWebHandler(templateManager, fm, authService)
	authWebRouter := auth.NewWebRouter(authWebHandler)
	authSeeder := auth.NewSeeder(assetsFS, engine, repo)

	app.MountWeb("/auth", authWebRouter)

	// SSG feature
	ssgService := ssg.NewService(repo)
	ssgWebHandler := ssg.NewWebHandler(templateManager, fm, ssgService)
	ssgWebRouter := ssg.NewWebRouter(ssgWebHandler, fm.Middlewares())
	ssgSeeder := ssg.NewSeeder(assetsFS, engine, repo)
	app.MountWeb("/ssg", ssgWebRouter)

	// Add deps
	app.Add(migrator)
	app.Add(seeder)
	app.Add(fm)
	app.Add(fileServer)
	app.Add(queryManager)
	app.Add(templateManager)
	app.Add(repo)
	app.Add(authService)
	app.Add(authWebHandler)
	app.Add(authWebRouter)
	app.Add(authSeeder)
	app.Add(ssgService)
	app.Add(ssgWebHandler)
	app.Add(ssgWebRouter)
	app.Add(ssgSeeder)

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
