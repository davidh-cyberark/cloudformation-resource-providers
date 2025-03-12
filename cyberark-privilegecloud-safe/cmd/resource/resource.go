package resource

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"

	"github.com/davidh-cyberark/conjur-sdk-go/conjur"
	"github.com/davidh-cyberark/identityadmin-sdk-go/identity"
	"github.com/davidh-cyberark/privilegeaccessmanager-sdk-go/pam"
)

// Create handles the Create event from the Cloudformation service.
func Create(req handler.Request, prevModel *Model, currentModel *Model) (handler.ProgressEvent, error) {
	// 1. CreateSafe
	// 2. CreateRole
	// 3. AddSafeMember(safe,role)
	newsafe, pcloudurl, err := CreateSafe(req, currentModel)
	var errMsg error
	if err != nil || len(newsafe.SafeURLID) == 0 {
		if err != nil {
			errMsg = fmt.Errorf("failed to create safe: %s", err.Error())
		} else {
			errMsg = fmt.Errorf("failed to create safe: no safe id was set in the response")
		}
		return handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Error adding safe",
			ResourceModel:   currentModel,
		}, errMsg
	}
	resourceid := fmt.Sprintf("PCLOUDURL=%s;SAFEURLID=%s", pcloudurl, newsafe.SafeURLID)

	// set the primaryIdentifier
	currentModel.SafeResourceId = &resourceid
	undef := "ROLENAME=undefined;ROLEID=undefined"
	currentModel.RoleId = &undef

	newrole, err := CreateRole(req, currentModel)
	if err != nil {
		return handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Error creating role for safe",
			ResourceModel:   currentModel,
		}, fmt.Errorf("failed to create role for safe: %s", err.Error())
	}
	roleBodyErr := identity.ReturnErrorWhenBodySuccessIsFalse(newrole.JSON200)
	if roleBodyErr != nil {
		return handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Error creating role for safe",
			ResourceModel:   currentModel,
		}, fmt.Errorf("failed to create role for safe: %s", roleBodyErr.Error())
	}

	// set roleid
	// ROLENAME=myrolename;ROLEID=79d63960-3f1e-4b6c-8a5d-2f7b9e1a0b1c
	roleid := fmt.Sprintf("ROLENAME=%s;ROLEID=%s", *currentModel.NewSafeRole, *newrole.JSON200.Result.Rowkey)
	currentModel.RoleId = &roleid

	member, err := AddRoleToSafe(req, currentModel, newsafe, newrole)
	if err != nil {
		return handler.ProgressEvent{
			OperationStatus: handler.Failed,
			Message:         "Error adding role to safe",
			ResourceModel:   currentModel,
		}, fmt.Errorf("failed adding role to safe: %s", err.Error())
	}

	msg := fmt.Sprintf("Create safe complete.\nSafe Name: %s, SafeURLID: %s, Member Name: %s, MemberID: %s",
		member.SafeName,
		member.SafeURLID,
		member.MemberName,
		member.MemberID)

	response := handler.ProgressEvent{
		OperationStatus: handler.Success,
		Message:         msg,
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

	log.Printf("Entering Create PAM Client Handler")
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

func CreateIdentityClientFromModel(req handler.Request, model *Model) (*identity.Service, error) {
	c, err := Configuration(req)
	if err != nil {
		return nil, err
	}

	log.Printf("Entering Create Identity Client Handler")
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
	if c.PAMAccessProperties.UserKey == nil {
		notset = fmt.Sprintf("%s  PAMUserKey", notset)
	}
	if c.PAMAccessProperties.PassKey == nil {
		notset = fmt.Sprintf("%s  PAMPassKey", notset)
	}

	if len(notset) > 0 {
		msg := strings.Replace(strings.Trim(notset, " "), "  ", ", ", -1)
		return nil, fmt.Errorf("identity/pam client properties not set: %s", msg)
	}

	ctx := context.Background()

	val, err := conjclient.FetchSecret(*c.PAMAccessProperties.UserKey)
	identityUser := string(val)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PAMUserKey from Conjur: %s", err.Error())
	}
	val, err = conjclient.FetchSecret(*c.PAMAccessProperties.PassKey)
	identityPass := string(val)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch PAMPassKey from Conjur: %s", err.Error())
	}

	userCreds := &identity.UserCredentials{
		User: identityUser,
		Pass: identityPass,
	}
	// Implements the AuthenticationProvider interface
	var authnProvider identity.AuthenticationProvider = userCreds

	// Create a new identity service
	service, errService := identity.NewService(ctx, *c.PAMAccessProperties.IDTenantURL, identity.ServiceWithAuthnProvider(authnProvider))
	if errService != nil {
		return nil, fmt.Errorf("failed to create identity service: %s", errService)
	}
	if service == nil {
		return nil, fmt.Errorf("failed to create identity service")
	}
	return service, nil
}

func CreateSafe(req handler.Request, model *Model) (*pam.PostAddSafeResponse, string, error) {
	pamclient, err := CreatePAMClientFromModel(req, model)
	if err != nil {
		return nil, "", err
	}
	pcloudurl := pamclient.Config.PcloudUrl
	err = pamclient.RefreshSession()
	if err != nil {
		return nil, pcloudurl, err
	}

	if model.NewSafeName == nil || len(*model.NewSafeName) == 0 {
		return nil, pcloudurl, fmt.Errorf("error NewSafeName is not set")
	}

	// https://docs.cyberark.com/privilege-cloud-standard/latest/en/content/webservices/add%20safe.htm
	// see "safeName" description
	// Ex safename sanitize implementation:
	//    badchars := "\\/:*<>.|?\"%&+" // Bad chars: \ / : * < > . | ? “% & +
	//    safename := RemoveSpecificChars(*model.NewSafeName, badchars)

	newsaferequest := pam.PostAddSafeRequest{
		SafeName: *model.NewSafeName,
	}
	newsafe, respcode, err := pamclient.AddSafe(newsaferequest)
	if err != nil {
		return nil, pcloudurl, fmt.Errorf("error, failed to add safe: %s", err.Error())
	}
	if respcode >= 300 {
		return nil, pcloudurl, fmt.Errorf("error, call to priv cloud returned non success code: %d. msg: %s", respcode, newsafe.ErrorResponse.ErrorMessage)
	}
	return &newsafe, pcloudurl, nil
}

func CreateRole(req handler.Request, model *Model) (*identity.PostRolesStoreRoleResponse, error) {
	identityservice, err := CreateIdentityClientFromModel(req, model)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, identity.ServiceKey, identityservice)

	if model.NewSafeRole == nil || len(*model.NewSafeRole) == 0 {
		name := fmt.Sprintf("%s-safe-role", *model.NewSafeName)
		model.NewSafeRole = &name
	}

	desc := fmt.Sprintf("Role for safe %s.", *model.NewSafeName)

	reqCreateRole := identity.PostRolesStoreRoleJSONRequestBody{
		Name:        *model.NewSafeRole,
		Description: &desc,
	}

	headers := map[string]string{
		"Accept":       "application/json",
		"Content-Type": "application/json",
	}
	ctx = context.WithValue(ctx, identity.HeadersKey, headers)

	roleResp, err := identityservice.Client.PostRolesStoreRoleWithResponse(ctx, reqCreateRole,
		identity.AddRequestHeaders,
		identityservice.AuthnProvider.AuthenticateRequest,
	)
	if err != nil {
		return roleResp, err
	}
	respErr := identity.ReturnErrorWhenBodySuccessIsFalse(roleResp.JSON200)
	if respErr != nil {
		return roleResp, respErr
	}
	return roleResp, nil
}

func AddRoleToSafe(req handler.Request, model *Model, safe *pam.PostAddSafeResponse, role *identity.PostRolesStoreRoleResponse) (*pam.PostAddMemberResponse, error) {
	pamclient, err := CreatePAMClientFromModel(req, model)
	if err != nil {
		return nil, err
	}
	err = pamclient.RefreshSession()
	if err != nil {
		return nil, err
	}

	memberReq := pam.PostAddMemberRequest{
		MemberName: *model.NewSafeRole,
		MemberType: "Role",
		IsReadOnly: true, // member can not update itself
		Permissions: pam.Permissions{
			ListAccounts:                           true,
			AddAccounts:                            true,
			UpdateAccountContent:                   true,
			UpdateAccountProperties:                true,
			InitiateCPMAccountManagementOperations: true,
			AccessWithoutConfirmation:              true,
			ManageSafeMembers:                      true,
		},
	}
	addsaferoleresp, _, err := pamclient.AddSafeMember(memberReq, *model.NewSafeName)
	if err != nil {
		return nil, fmt.Errorf("error: could not add safe member: %s", err.Error())
	}
	return &addsaferoleresp, nil
}
