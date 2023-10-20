package tapo

import "time"

const (
	// Encryption method
	ENCRYPTION_METHOD_AES_128_CBC EncryptionMethod = "AES-128-CBC"

	// Event type
	EVENT_TYPE_PD     EventType = "PD"
	EVENT_TYPE_MOTION EventType = "MOTION"
)

type EncryptionMethod string
type EventType string
type DateTime time.Time

type GetVideosDevicesResponse struct {
	DeviceList []DeviceListItem `json:"deviceList"`
}

type DeviceListItem struct {
	Alias       string `json:"alias"`
	DeviceId    string `json:"deviceId"`
	DeviceMac   string `json:"deviceMac"`
	DeviceModel string `json:"deviceModel"`
	DeviceType  string `json:"deviceType"`
}

type ListActivitiesByDateRequest struct {
	DeviceId         string      `json:"deviceId"`
	StartTime        string      `json:"startTime"`
	EndTime          string      `json:"endTime"`
	Source           string      `json:"source"`
	EventTypeFilters []EventType `json:"eventTypeFilters"`
	Page             int         `json:"page"`
	PageSize         int         `json:"pageSize"`
}

type ListActivitiesByDateResponse struct {
	Listing  []Activity `json:"listing"`
	Page     int        `json:"page"`
	PageSize int        `json:"pageSize"`
	Total    int        `json:"total"`
}

type GetVideosListRequest struct {
	DeviceId  string `json:"deviceId"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
	Order     string `json:"order"`
	Page      int    `json:"page"`
	PageSize  int    `json:"pageSize"`
}

type GetVideosListResponse struct {
	Index    []VideoIndex `json:"index"`
	DeviceId string       `json:"deviceId"`
	Page     int          `json:"page"`
	PageSize int          `json:"pageSize"`
	Total    int          `json:"total"`
}

type VideoIndex struct {
	UUID           string       `json:"uuid"`
	CreatedTime    int64        `json:"createdTime"`
	EventLocalTime string       `json:"eventLocalTime"`
	EventTypeList  []EventType  `json:"eventTypeList"`
	Video          []IndexVideo `json:"video"`
	Image          []IndexImage `json:"image"`
}

type IndexVideo struct {
	Uri              string           `json:"uri"`
	Duration         int              `json:"duration"`
	M3U8             string           `json:"m3u8"`
	StartTimestamp   int64            `json:"startTimestamp"`
	UriExpiresAt     int64            `json:"uriExpiresAt"`
	EncryptionMethod EncryptionMethod `json:"encryptionMethod"`
	DecryptionInfo   DecryptionInfo   `json:"decryptionInfo"`
	Resolution       string           `json:"resolution"`
}

type IndexImage struct {
	Uri              string         `json:"uri"`
	Length           int64          `json:"length"`
	UriExpiresAt     int64          `json:"uriExpiresAt"`
	EncryptionMethod string         `json:"encryptionMethod"`
	DecryptionInfo   DecryptionInfo `json:"decryptionInfo"`
	EventTypeNames   []string       `json:"eventTypeNames"`
	EventTimestamp   int64          `json:"eventTimestamp"`
}

type Activity struct {
	Id         string   `json:"id"`
	Level      int      `json:"level"`
	Source     int      `json:"source"`
	CreatedOn  int64    `json:"createdOn"`
	UpdatedOn  int64    `json:"updatedOn"`
	AccountId  string   `json:"accountId"`
	Device     Device   `json:"device"`
	Event      Event    `json:"event"`
	SearchTags []string `json:"searchTags"`
}

type Device struct {
	DeviceAlias    string `json:"deviceAlias"`
	DeviceCategory string `json:"deviceCategory"`
	DeviceId       string `json:"deviceId"`
	DeviceModel    string `json:"deviceModel"`
	DeviceType     string `json:"deviceType"`
}

type Event struct {
	Data           EventData `json:"data"`
	EventLocalTime string    `json:"eventLocalTime"`
	Id             string    `json:"id"`
	Name           string    `json:"name"`
	Timestamp      int64     `json:"timestamp"`
	Type           string    `json:"type"`
}

type EventData struct {
	CloudStorage        CloudStorage `json:"cloudStorage"`
	CloudStorageEnabled bool         `json:"cloudStorageEnabled"`
	Snapshot            Snapshot     `json:"snapshot"`
	TypeUri             string       `json:"typeUri"`
	Video               Video        `json:"video"`
}

type CloudStorage struct {
	Enabled                  bool `json:"enabled"`
	RichNotificationsFeature bool `json:"richNotificationsFeature"`
}

type Snapshot struct {
	DecryptionInfo   DecryptionInfo `json:"decryptionInfo"`
	DeletedOn        int64          `json:"deletedOn"`
	EncryptionMethod string         `json:"encryptionMethod"`
	Size             int64          `json:"size"`
	Status           string         `json:"status"`
	Url              string         `json:"url"`
}

type Video struct {
	Duration         int    `json:"duration"`
	EncryptionMethod string `json:"encryptionMethod"`
	Resolution       string `json:"resolution"`
	Status           string `json:"status"`
	StreamUrl        string `json:"streamUrl"`
}

type DecryptionInfo struct {
	IV     string `json:"iv"`
	Key    string `json:"key"`
	KeyUri string `json:"keyUri"`
}
