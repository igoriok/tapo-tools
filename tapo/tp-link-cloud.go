package tapo

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/tls"
	"encoding/base64"
	"encoding/hex"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/google/uuid"
)

const (
	accessKey    = "4d11b6b9d5ea4d19a829adbb9714b057"
	accessSecret = "6ed7d97f3e73467f8a5bab90b577ba4c"
)

type TpLinkCloudClient struct {
	Client *resty.Client
}

func NewTpLinkCloudClient(baseURL string) *TpLinkCloudClient {

	client := resty.New()

	client.SetBaseURL(baseURL)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})
	client.OnBeforeRequest(onBeforeRequest)

	return &TpLinkCloudClient{
		Client: client,
	}
}

func onBeforeRequest(c *resty.Client, r *resty.Request) error {

	body := r.Body

	if body != nil {

		now := "9999999999" // time.Now().String()
		nonce := uuid.New().String()

		bodyBytes, err := c.JSONMarshal(body)

		if err != nil {
			return err
		}

		hash := md5.Sum(bodyBytes)
		hashString := base64.StdEncoding.EncodeToString(hash[:])

		hmac := hmac.New(sha1.New, []byte(accessSecret))

		hmac.Write([]byte(hashString))
		hmac.Write([]byte("\n"))
		hmac.Write([]byte(now))
		hmac.Write([]byte("\n"))
		hmac.Write([]byte(nonce))
		hmac.Write([]byte("\n"))
		hmac.Write([]byte(r.URL))

		signature := hex.EncodeToString(hmac.Sum(nil))

		header := fmt.Sprintf(
			"Timestamp=%s, Nonce=%s, AccessKey=%s, Signature=%s",
			now,
			nonce,
			accessKey,
			signature,
		)

		r.SetHeader("Content-MD5", hashString)
		r.SetHeader("X-Authorization", header)
		r.SetHeader("Content-Type", "application/json; charset=UTF-8")
		r.SetHeader("User-Agent", "okhttp/3.14.9")
		r.SetBody(bodyBytes)
	}

	return nil
}

func (c *TpLinkCloudClient) HelloCloud(params *GenericTapoParams, req *HelloCloudRequest) (*GenericTapoResponse[HelloCloudResult], error) {

	resp, err := c.Client.R().
		SetQueryParams(getQueryParams(params)).
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&GenericTapoResponse[HelloCloudResult]{}).
		Post("/api/v2/common/helloCloud")

	return resp.Result().(*GenericTapoResponse[HelloCloudResult]), err
}

func (c *TpLinkCloudClient) GetAppServiceUrl(params *GenericTapoParams, req *GetAppServiceUrlRequest) (*GenericTapoResponse[GetAppServiceUrlResult], error) {

	resp, err := c.Client.R().
		SetQueryParams(getQueryParams(params)).
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&GenericTapoResponse[GetAppServiceUrlResult]{}).
		Post("/api/v2/common/getAppServiceUrl")

	return resp.Result().(*GenericTapoResponse[GetAppServiceUrlResult]), err
}

func (c *TpLinkCloudClient) AccountLogin(params *GenericTapoParams, req *AccountLoginRequest) (*GenericTapoResponse[AccountLoginResult], error) {

	resp, err := c.Client.R().
		SetQueryParams(getQueryParams(params)).
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&GenericTapoResponse[AccountLoginResult]{}).
		Post("/api/v2/account/login")

	return resp.Result().(*GenericTapoResponse[AccountLoginResult]), err
}

func (c *TpLinkCloudClient) AccountLogout(params *GenericTapoParams, req *AccountLogoutRequest) (*GenericTapoResponse[EmptyTapoResult], error) {

	resp, err := c.Client.R().
		SetQueryParams(getQueryParams(params)).
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&GenericTapoResponse[EmptyTapoResult]{}).
		Post("/api/v2/account/logout")

	return resp.Result().(*GenericTapoResponse[EmptyTapoResult]), err
}

func (c *TpLinkCloudClient) GetMFAFeatureStatus(params *GenericTapoParams, req *GetMFAFeatureStatusRequest) (*GenericTapoResponse[GetMFAFeatureStatusResult], error) {

	resp, err := c.Client.R().
		SetQueryParams(getQueryParams(params)).
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&GenericTapoResponse[GetMFAFeatureStatusResult]{}).
		Post("/api/v2/account/getMFAFeatureStatus")

	return resp.Result().(*GenericTapoResponse[GetMFAFeatureStatusResult]), err
}

func (c *TpLinkCloudClient) GetPushVC4TerminalMFA(params *GenericTapoParams, req *GetPushVC4TerminalMFARequest) (*GenericTapoResponse[EmptyTapoResult], error) {

	resp, err := c.Client.R().
		SetQueryParams(getQueryParams(params)).
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&GenericTapoResponse[EmptyTapoResult]{}).
		Post("/api/v2/account/getPushVC4TerminalMFA")

	return resp.Result().(*GenericTapoResponse[EmptyTapoResult]), err
}

func (c *TpLinkCloudClient) GetEmailVC4TerminalMFA(params *GenericTapoParams, req *GetEmailVC4TerminalMFARequest) (*GenericTapoResponse[EmptyTapoResult], error) {

	resp, err := c.Client.R().
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&GenericTapoResponse[EmptyTapoResult]{}).
		Post("/api/v2/account/getEmailVC4TerminalMFA")

	return resp.Result().(*GenericTapoResponse[EmptyTapoResult]), err
}

func (c *TpLinkCloudClient) CheckMFACodeAndLogin(params *GenericTapoParams, req *CheckMFACodeAndLoginRequest) (*GenericTapoResponse[CheckMFACodeAndLoginResult], error) {

	resp, err := c.Client.R().
		SetQueryParams(getQueryParams(params)).
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&GenericTapoResponse[CheckMFACodeAndLoginResult]{}).
		Post("/api/v2/account/checkMFACodeAndLogin")

	return resp.Result().(*GenericTapoResponse[CheckMFACodeAndLoginResult]), err
}

func getQueryParams(params *GenericTapoParams) map[string]string {

	queryParams := map[string]string{
		"appName":  params.AppName,
		"appVer":   params.AppVer,
		"netType":  params.NetType,
		"termID":   params.TermID,
		"termName": params.TermName,
		"termMeta": params.TermMeta,
		"ospf":     params.OSPF,
		"brand":    params.Brand,
		"model":    params.Model,
		"locale":   params.Locale,
		"token":    params.Token,
	}

	if params.Token != "" {
		queryParams["token"] = params.Token
	}

	return queryParams
}

type HelloCloudRequest struct {
	AppPackageName string `json:"appPackageName"`
	AppType        string `json:"appType"`
	TcspVer        string `json:"tcspVer"`
	TerminalUUID   string `json:"terminalUUID"`
}

type HelloCloudResult struct {
	TcspStatus int `json:"tcspStatus"`
}

type GetAppServiceUrlRequest struct {
	ServiceIds []string `json:"serviceIds"`
}

type GetAppServiceUrlResult struct {
	RegionCode  string            `json:"regionCode"`
	ServiceUrls map[string]string `json:"serviceUrls"`
}

type AccountLoginRequest struct {
	AppType            string `json:"appType"`
	AppVersion         string `json:"appVersion"`
	CloudUserName      string `json:"cloudUserName"`
	CloudPassword      string `json:"cloudPassword"`
	Platform           string `json:"platform"`
	RefreshTokenNeeded bool   `json:"refreshTokenNeeded"`
	TerminalMeta       string `json:"terminalMeta"`
	TerminalName       string `json:"terminalName"`
	TerminalUUID       string `json:"terminalUUID"`
}

type AccountLogoutRequest struct {
	CloudUserName string `json:"cloudUserName"`
}

type AccountLoginResult struct {
	AccountId    string `json:"accountId"`
	AppServerUrl string `json:"appServerUrl"`
	Email        string `json:"email"`
	Nickname     string `json:"nickname"`
	RegTime      string `json:"regTime"`
	RegionCode   string `json:"regionCode"`
	RiskDetected int    `json:"riskDetected"`
	Token        string `json:"token"`

	ErrorCode string `json:"errorCode"`
	ErrorMsg  string `json:"errorMsg"`

	MFAEmail          string    `json:"mfaEmail"`
	MFAProcessId      string    `json:"MFAProcessId"`
	SupportedMFATypes []MFAType `json:"supportedMFATypes"`
}

type GetMFAFeatureStatusRequest struct {
}

type GetMFAFeatureStatusResult struct {
	MFAFeatureEnabled int       `json:"MFAFeatureEnabled"`
	Email             string    `json:"email"`
	MFAEmail          string    `json:"mfaEmail"`
	SupportedMFATypes []MFAType `json:"supportedMFATypes"`
	TerminalBound     bool      `json:"terminalBound"`
}

type GetPushVC4TerminalMFARequest struct {
	AppType       string `json:"appType"`
	CloudPassword string `json:"cloudPassword"`
	CloudUserName string `json:"cloudUserName"`
	TerminalUUID  string `json:"terminalUUID"`
}

type GetEmailVC4TerminalMFARequest struct {
	AppType       string `json:"appType"`
	CloudUserName string `json:"cloudUserName"`
	CloudPassword string `json:"cloudPassword"`
	TerminalUUID  string `json:"terminalUUID"`
}

type CheckMFACodeAndLoginRequest struct {
	MFAProcessId        string  `json:"MFAProcessId"`
	MFAType             MFAType `json:"MFAType"`
	AppType             string  `json:"appType"`
	CloudUserName       string  `json:"cloudUserName"`
	Code                string  `json:"code"`
	TerminalBindEnabled bool    `json:"terminalBindEnabled"`
}

type CheckMFACodeAndLoginResult struct {
	AccountId    string `json:"accountId"`
	AppServerUrl string `json:"appServerUrl"`
	Email        string `json:"email"`
	ErrorCode    string `json:"errorCode"`
	Nickname     string `json:"nickname"`
	RegTime      string `json:"regTime"`
	RegionCode   string `json:"regionCode"`
	Token        string `json:"token"`
}

type GenericTapoParams struct {
	AppName  string `json:"appName"`
	AppVer   string `json:"appVer"`
	NetType  string `json:"netType"`
	TermID   string `json:"termID"`
	TermName string `json:"termName"`
	TermMeta string `json:"termMeta"`
	OSPF     string `json:"ospf"`
	Brand    string `json:"brand"`
	Model    string `json:"model"`
	Locale   string `json:"locale"`
	Token    string `json:"token"`
}

type GenericTapoResponse[R any] struct {
	ErrorCode int    `json:"error_code"`
	Message   string `json:"msg"`
	Result    R      `json:"result"`
}

type EmptyTapoResult struct {
}

type MFAType int

const (
	MfaTypeNone MFAType = iota
	MFATypePush
	MFATypeEmail
)
