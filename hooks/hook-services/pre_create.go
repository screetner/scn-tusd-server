package hook_services

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/tus/tusd/v2/pkg/hooks"
	"io"
	"log"
)

func PreCreateHookHandler(req hooks.HookRequest) (res hooks.HookResponse, err error) {
	res.HTTPResponse.Header = make(map[string]string)

	fileName, fileNameErr := validateMetaDataField(&res, &req, "filename")
	if fileNameErr != nil {
		return res, fileNameErr
	}

	sessionName, sessionNameErr := validateMetaDataField(&res, &req, "sessionName")
	if sessionNameErr != nil {
		return res, sessionNameErr
	}

	uniqueSessionId := uid()

	if res.ChangeFileInfo.MetaData == nil {
		res.ChangeFileInfo.MetaData = make(map[string]string)
	}

	// TODO: validate user upload permission

	// TODO: decode organizationId from jwt token
	organizationId := "orgId1"

	// TODO: decode userId from jwt token
	userId := uid()

	directoryName := organizationId
	id := directoryName + "/records/" + sessionName + "_" + uniqueSessionId + "/" + fileName

	res.ChangeFileInfo.ID = id
	res.ChangeFileInfo.MetaData["recorderUserId"] = userId
	log.Printf("PRE_CREATE: Uploading %s at : %s\n", fileName, id)

	return res, nil
}

func validateMetaDataField(res *hooks.HookResponse, req *hooks.HookRequest, fieldName string) (string, error) {
	fieldValue, ok := req.Event.Upload.MetaData[fieldName]
	if !ok {
		errorMessage := fmt.Sprintf("no %s provided in the request metadata", fieldName)
		res.RejectUpload = true
		res.HTTPResponse.StatusCode = 400
		res.HTTPResponse.Body = errorMessage
		res.HTTPResponse.Header["X-Some-Header"] = "yes"
		return "", fmt.Errorf(errorMessage)
	}
	return fieldValue, nil
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
