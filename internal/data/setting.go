package data

type APNS struct {
	Token  string `json:"token"`
	KeyID  string `json:"key_id"`
	TeamID string `json:"team_id"`
}

type FCM struct {
	ServiceAccount string `json:"service_account"`
}

type DevelopmentMode string

const (
	DevelopmentMode_Sandbox    DevelopmentMode = DevelopmentMode("sandbox")
	DevelopmentMode_Production DevelopmentMode = DevelopmentMode("production")
)

type Platform string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"
)

type LegacyAPNS struct {
	Certificate     string          `json:"certificate"`
	DecryptPassword string          `json:"decrypt_password"`
	DevelopmentMode DevelopmentMode `json:"development_mode"`
}
