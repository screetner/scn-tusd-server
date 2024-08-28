package main

import (
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	functionHook "scn-tusd-server/hooks"
	"scn-tusd-server/services"
	"strings"

	"github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/hooks"
)

var stdout = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)
var stderr = log.New(os.Stderr, "", log.LstdFlags|log.Lmicroseconds)

func main() {
	// Credits to tusd
	// Link: https://github.com/tus/tusd/blob/main/cmd/tusd/cli/composer.go

	composer, err := services.CreateComposer()
	if err != nil {
		stderr.Fatalf(err.Error())
	}

	err = godotenv.Load()
	if err != nil {
		stderr.Fatalf("unable to load .env file: %v", err)
	}

	basePath := os.Getenv("BASE_PATH")
	if basePath == "" {
		basePath = "/"
	}

	config := handler.Config{
		BasePath:              basePath,
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	}

	enabledHooks := []hooks.HookType{"pre-create", "pre-finish"}

	tusHandler, err := createHandler(config, enabledHooks)
	if err != nil {
		stderr.Fatalf("unable to listen: %s", err)
	}

	go func() {
		for {
			event := <-tusHandler.CompleteUploads
			stdout.Printf("upload %s finished\n", event.Upload.ID)
		}
	}()

	http.Handle(basePath, http.StripPrefix(basePath, tusHandler))
	http.Handle(basePath+"/", http.StripPrefix(basePath+"/", tusHandler))
	http.HandleFunc("/", services.DisplayGreeting)

	stdout.Printf("Tusd server is hosting at http://localhost:8080/files 𓀐\n")

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		stderr.Fatalf("unable to listen: %s", err)
	}
}

func createHandler(config handler.Config, enabledHooks []hooks.HookType) (tusHandler *handler.Handler, err error) {
	hookHandler := functionHook.GetHookHandler(&config)
	if hookHandler != nil {
		tusHandler, err = hooks.NewHandlerWithHooks(&config, hookHandler, enabledHooks)

		var enabledHooksString []string
		for _, h := range enabledHooks {
			enabledHooksString = append(enabledHooksString, string(h))
		}

		stdout.Printf("Enabled hook events: %s", strings.Join(enabledHooksString, ", "))

	} else {
		tusHandler, err = handler.NewHandler(config)
	}

	return tusHandler, err
}
