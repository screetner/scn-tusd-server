package hook_services

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tus/tusd/v2/pkg/hooks"
	"io"
	"log"
	"scn-tusd-server/hooks/types"
	services "scn-tusd-server/services"
	"strings"
)

func PreCreateHookHandler(req hooks.HookRequest) (res hooks.HookResponse, err error) {
	res.HTTPResponse.Header = make(map[string]string)

	fileName, fileNameErr := validateMetaDataField(&res, &req, "filename")
	if fileNameErr != nil {
		return res, fileNameErr
	}

	// TODO: test this
	sessionCloudName, sessionCloudNameErr := validateMetaDataField(&res, &req, "sessionCloudName")
	if sessionCloudNameErr != nil {
		return res, sessionCloudNameErr
	}

	if res.ChangeFileInfo.MetaData == nil {
		res.ChangeFileInfo.MetaData = make(map[string]string)
	}

	headers := make(map[string]string)
	for key, values := range req.Event.HTTPRequest.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}
	body := bytes.NewBuffer([]byte(`{}`))

	backendResp, backedErr := services.GetAPIClient().Get("/tusd/user/info", headers, body)
	if backedErr != nil {
		// Handle the error
		log.Printf("Error validating: %v", err)
		return res, backedErr
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Printf("Error closing response body: %v", err)
		}
	}(backendResp.Body)

	var userInfo types.UserInfo // Use the UserInfo type from the types package
	if err := json.NewDecoder(backendResp.Body).Decode(&userInfo); err != nil {
		log.Printf("Error decoding JSON response: %v", err)
		return res, err
	}

	if backendResp.StatusCode == 403 {
		res.RejectUpload = true
		res.HTTPResponse.StatusCode = 403
		res.HTTPResponse.Body = "User's tusd token is invalid"
		return res, nil
	}

	// TODO: validate user upload permission
	hasPermission := userInfo.AbilityScope.Mobile.Access
	if !hasPermission {
		res.RejectUpload = true
		res.HTTPResponse.StatusCode = 403
		res.HTTPResponse.Body = "User does not have permission to upload"
		return res, nil
	}

	organizationId := userInfo.OrgId
	organizationName := userInfo.OrgName

	userId := userInfo.UserId
	username := userInfo.UserName

	organizationDirectory := fmt.Sprintf("%s_%s", organizationName, organizationId)
	id := fmt.Sprintf("%s/records/%s/%s", organizationDirectory, sessionCloudName, fileName)
	id = strings.ReplaceAll(id, " ", "_")

	res.ChangeFileInfo.ID = id
	res.ChangeFileInfo.MetaData["recorderUserId"] = userId
	res.ChangeFileInfo.MetaData["recorderUserName"] = username
	log.Printf("PRE_CREATE: Uploading %s at : %s\n", fileName, id)

	res.HTTPResponse.StatusCode = 200

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
