package resource

import (
	"fmt"
	"log"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"

	"github.com/davidh-cyberark/conjur-sdk-go/conjur"
	"github.com/davidh-cyberark/privilegeaccessmanager-sdk-go/pam"
)

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	err := AddSafe(req, currentModel)
	if err != nil {
		return handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Error adding safe",
			ResourceModel:   currentModel,
		}, fmt.Errorf("failed to create safe: %s", err.Error())
	}

	if currentModel.SafeResourceId == nil {
		return handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "No safe id was generated",
			ResourceModel:   currentModel,
		}, fmt.Errorf("failed to create safe")
	}

	// Construct a new handler.ProgressEvent and return it
	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Create complete",
		ResourceModel:   currentModel,
	}

	return response, nil
}

// Read handles the Read event from the Cloudformation service.
func Read(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Read complete",
		ResourceModel:   currentModel,
	}
	return response, nil
}

// Update handles the Update event from the Cloudformation service.
func Update(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Update complete",
		ResourceModel:   currentModel,
	}
	return response, nil
}

// Delete handles the Delete event from the Cloudformation service.
func Delete(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "Delete complete",
		ResourceModel:   currentModel,
	}
	return response, nil
}

// List handles the List event from the Cloudformation service.
func List(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         "List complete",
		ResourceModel:   currentModel,
	}
	return response, nil
}

func CreateConjurClientFromModel(req handler.Request, model *Model) (*conjur.Client, error) {
	c, err := Configuration(req)
	if err != nil {
		return nil, err
	}
	notset := ""
	if c.ConjurAPIURL == nil {
		notset = fmt.Sprintf("%s  ConjurAPIURL", notset)
	}
	if c.ConjurAuthenticatorProperties.ConjurIdentity == nil {
		notset = fmt.Sprintf("%s  ConjurIdentity", notset)
	}
	if c.ConjurAuthenticatorProperties.ConjurAuthenticator == nil {
		notset = fmt.Sprintf("%s  ConjurAuthenticator", notset)
	}
	if c.ConjurAuthenticatorProperties.ConjurAccount == nil {
		notset = fmt.Sprintf("%s  ConjurAccount", notset)
	}
	if c.ConjurAuthenticatorProperties.ConjurIAMRoleARN == nil {
		notset = fmt.Sprintf("%s  ConjurIAMRoleARN", notset)
	}
	if c.ConjurAuthenticatorProperties.ConjurAWSRegion == nil {
		notset = fmt.Sprintf("%s  ConjurAWSRegion", notset)
	}

	if len(notset) > 0 {
		msg := strings.Replace(strings.Trim(notset, " "), "  ", ", ", -1)
		return nil, fmt.Errorf("conjur client properties not set: %s", msg)
	}

	currentRoleCreds, err := req.Session.Config.Credentials.Get()
	if err != nil {
		return nil, fmt.Errorf("unable to get session credentials: %s", err.Error())
	}

	key := currentRoleCreds.AccessKeyID
	sec := currentRoleCreds.SecretAccessKey
	stk := currentRoleCreds.SessionToken

	// use these creds to configigure sts client which will
	// be used to assume role of assumedrole.arn
	rolecreds := conjur.NewAWSProviderCredentials(
		conjur.WithAWSProviderCredentialsAccessKey(key),
		conjur.WithAWSProviderCredentialsAccessSecret(sec),
		conjur.WithAWSProviderCredentialsSessionToken(stk))
	assumedrolecreds := conjur.NewAWSProviderCredentials(
		conjur.WithAWSProviderCredentialsArn(*c.ConjurAuthenticatorProperties.ConjurIAMRoleARN),
	)

	awsprovider := conjur.NewAWSProvider(
		conjur.WithRegion(*c.ConjurAuthenticatorProperties.ConjurAWSRegion),
		conjur.WithAWSProviderRoleCredentials(rolecreds),
		conjur.WithAWSProviderAssumedRoleCredentials(assumedrolecreds))

	client := conjur.NewClient(*c.ConjurAPIURL,
		conjur.WithAccount(*c.ConjurAuthenticatorProperties.ConjurAccount),
		conjur.WithIdentity(*c.ConjurAuthenticatorProperties.ConjurIdentity),
		conjur.WithAuthenticator(*c.ConjurAuthenticatorProperties.ConjurAuthenticator),
		conjur.WithAwsProvider(&awsprovider))
	return client, nil
}

func CreatePAMClientFromModel(req handler.Request, model *Model) (*pam.Client, error) {
	c, err := Configuration(req)
	if err != nil {
		return nil, err
	}

	log.Printf("Entering Create Handler")
	conjclient, conjclienterr := CreateConjurClientFromModel(req, model)
	if conjclienterr != nil {
		return nil, fmt.Errorf("failed to create conjur client: %s", conjclienterr.Error())
	}
	if conjclient == nil {
		return nil, fmt.Errorf("failed to create conjur client")
	}

	notset := ""
	if c.PAMAccessProperties.IDTenantURL == nil {
		notset = fmt.Sprintf("%s  IDTenantURL", notset)
	}

	// Conjur keys used for fetching PAM creds from Conjur
	if c.PAMAccessProperties.PcloudURLKey == nil {
		notset = fmt.Sprintf("%s  PAMPcloudURLKey", notset)
	}
	if c.PAMAccessProperties.UserKey == nil {
		notset = fmt.Sprintf("%s  PAMUserKey", notset)
	}
	if c.PAMAccessProperties.PassKey == nil {
		notset = fmt.Sprintf("%s  PAMPassKey", notset)
	}

	if len(notset) > 0 {
		msg := strings.Replace(strings.Trim(notset, " "), "  ", ", ", -1)
		return nil, fmt.Errorf("pam client properties not set: %s", msg)
	}

	pamconf := pam.Config{
		IdTenantUrl: *c.PAMAccessProperties.IDTenantURL,
	}

	val, err := conjclient.FetchSecret(*c.PAMAccessProperties.PcloudURLKey)
	pamconf.PcloudUrl = string(val)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PAMPcloudURLKey from Conjur: %s", err.Error())
	}
	val, err = conjclient.FetchSecret(*c.PAMAccessProperties.UserKey)
	pamconf.User = string(val)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PAMUserKey from Conjur: %s", err.Error())
	}
	val, err = conjclient.FetchSecret(*c.PAMAccessProperties.PassKey)
	pamconf.Pass = string(val)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PAMPassKey from Conjur: %s", err.Error())
	}

	pamclient := pam.NewClient(pamconf.PcloudUrl, &pamconf)
	return pamclient, nil
}

func AddSafe(req handler.Request, model *Model) error {
	pamclient, err := CreatePAMClientFromModel(req, model)
	if err != nil {
		return err
	}
	err = pamclient.RefreshSession()
	if err != nil {
		return err
	}

	if model.NewSafeName == nil || len(*model.NewSafeName) == 0 {
		return fmt.Errorf("error NewSafeName is not set")
	}
	newsaferequest := pam.PostAddSafeRequest{
		SafeName: *model.NewSafeName,
	}
	newsafe, respcode, err := pamclient.AddSafe(newsaferequest)
	if err != nil {
		return fmt.Errorf("error, failed to add safe: %s", err.Error())
	}
	if respcode >= 300 {
		return fmt.Errorf("error, call to priv cloud returned non success code: %d", respcode)
	}

	if len(newsafe.SafeURLID) == 0 {
		return fmt.Errorf("no safe url id was set in the response")
	}
	resourceid := fmt.Sprintf("PCLOUDURL=%s|SAFEURLID=%s", pamclient.Config.PcloudUrl, newsafe.SafeURLID)

	// set the primaryIdentifier
	model.SafeResourceId = &resourceid
	return nil
}
