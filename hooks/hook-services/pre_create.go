package hook_services

import (
	"github.com/tus/tusd/v2/pkg/hooks"
	"log"
)

func PreCreateHookHandler(req hooks.HookRequest) (res hooks.HookResponse, err error) {
	res.HTTPResponse.Header = make(map[string]string)

	if filename, ok := req.Event.Upload.MetaData["filename"]; !ok {
		res.RejectUpload = true
		res.HTTPResponse.StatusCode = 400
		res.HTTPResponse.Body = "no filename provided"
		res.HTTPResponse.Header["X-Some-Header"] = "yes"
	} else {
		log.Printf("PRE_CREATE: Creating new file: %s\n", filename)
	}

	return res, nil
}
