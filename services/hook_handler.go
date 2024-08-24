package services

import (
	"github.com/tus/tusd/v2/pkg/handler"
	"github.com/tus/tusd/v2/pkg/hooks"
	"github.com/tus/tusd/v2/pkg/hooks/plugin"
	"log"
	"os"
)

var stdout = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

func GetHookHandler(config *handler.Config) hooks.HookHandler {
	pluginHookPath := "plugin/hook_plugin"

	stdout.Printf("Using '%s' to load plugin for hooks", pluginHookPath)

	return &plugin.PluginHook{
		Path: pluginHookPath,
	}
}
