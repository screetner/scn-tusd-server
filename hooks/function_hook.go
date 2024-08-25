package hooks

import (
	"github.com/tus/tusd/v2/pkg/hooks"
	hook_services "scn-tusd-server/hooks/hook-services"
)

type FunctionHook struct {
	handlerImpl hooks.HookHandler
}

func (h *FunctionHook) Setup() error {
	stdout.Printf("Running hook handler with built-in functions")
	return nil
}

func (h *FunctionHook) InvokeHook(req hooks.HookRequest) (res hooks.HookResponse, err error) {
	switch req.Type {
	case hooks.HookPreCreate:
		res, err = hook_services.PreCreateHookHandler(req)
		break
	case hooks.HookPreFinish:
		res, err = hook_services.PreFinishHookHandler(req)
		break
	}

	return res, nil
}
