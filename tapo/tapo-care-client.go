package tapo

import (
	"crypto/tls"
	"fmt"

	"github.com/go-resty/resty/v2"
)

type TapoCareClient struct {
	Client *resty.Client
}

func NewTapoCareClient(baseURL string, token string, termID string) *TapoCareClient {

	appName := "TP-Link_Tapo_Android"
	appVersion := "3.0.536"

	client := resty.New()

	client.SetBaseURL(baseURL)
	client.SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true})

	client.SetHeader("Authorization", fmt.Sprintf("ut|%s", token))
	client.SetHeader("User-Agent", fmt.Sprintf("%s/%s(Pixel 7/;Android 14)", appName, appVersion))

	client.SetHeaderVerbatim("app-cid", fmt.Sprintf("app:%s:%s", appName, termID))
	client.SetHeaderVerbatim("x-app-name", appName)
	client.SetHeaderVerbatim("x-app-version", termID)
	client.SetHeaderVerbatim("x-ospf", "Android 14")
	client.SetHeaderVerbatim("x-net-type", "wifi")
	client.SetHeaderVerbatim("x-locale", "en_US")

	return &TapoCareClient{
		Client: client,
	}
}

func (c *TapoCareClient) GetVideosDevices() (*GetVideosDevicesResponse, error) {

	resp, err := c.Client.R().
		ForceContentType("application/json").
		SetResult(&GetVideosDevicesResponse{}).
		Get("/v2/videos/devices")

	return resp.Result().(*GetVideosDevicesResponse), err
}

func (c *TapoCareClient) ListActivitiesByDate(req *ListActivitiesByDateRequest) (*ListActivitiesByDateResponse, error) {

	resp, err := c.Client.R().
		SetBody(req).
		ForceContentType("application/json").
		SetResult(&ListActivitiesByDateResponse{}).
		Post("/v1/activities/listActivitiesByDate")

	return resp.Result().(*ListActivitiesByDateResponse), err
}

func (c *TapoCareClient) GetVideosList(req *GetVideosListRequest) (*GetVideosListResponse, error) {
	resp, err := c.Client.R().
		SetQueryParams(map[string]string{
			"deviceId":  req.DeviceId,
			"startTime": req.StartTime,
			"endTime":   req.EndTime,
			"order":     req.Order,
			"page":      fmt.Sprintf("%d", req.Page),
			"pageSize":  fmt.Sprintf("%d", req.PageSize),
		}).
		ForceContentType("application/json").
		SetResult(&GetVideosListResponse{}).
		Get("/v2/videos/list")

	return resp.Result().(*GetVideosListResponse), err
}
