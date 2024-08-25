package hook_services

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/tus/tusd/v2/pkg/hooks"
	"io"
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

		if res.ChangeFileInfo.Storage == nil {
			res.ChangeFileInfo.Storage = make(map[string]string)
		}

		directoryName := "records"
		id := directoryName + "/" + uid()

		res.ChangeFileInfo.ID = id
		log.Printf("PRE_CREATE: Uploading %s at : %s\n", filename, id)
	}

	return res, nil
}

// found in github.com/tus/tusd/v2/internal/hooks/uid
func uid() string {
	id := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, id)
	if err != nil {
		// This is probably an appropriate way to handle errors from our source
		// for random bits.
		panic(err)
	}
	return hex.EncodeToString(id)
}
