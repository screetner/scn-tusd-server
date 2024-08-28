package hook_services

import (
	"github.com/tus/tusd/v2/pkg/hooks"
	"log"
)

func PreFinishHookHandler(req hooks.HookRequest) (res hooks.HookResponse, err error) {
	res.HTTPResponse.Header = make(map[string]string)

	id := req.Event.Upload.ID
	size := req.Event.Upload.Size
	storage := req.Event.Upload.Storage

	log.Printf("PRE_FINISH: Upload %s (%d bytes) is finished. Find the file at:\n", id, size)
	log.Println(storage)

	return res, nil
}
