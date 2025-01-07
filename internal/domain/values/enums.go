package values

type MobilePlatform string

const (
	MobilePlatformIOS     MobilePlatform = MobilePlatform("ios")
	MobilePlatformAndroid MobilePlatform = MobilePlatform("android")
)

type APNSType string

const (
	APNSTypeTokenBase APNSType = APNSType("token")
	APNSTypeLegacy    APNSType = APNSType("legacy")
)

type APNSEnvironment string

const (
	APNSEnvironmentSandbox    APNSEnvironment = APNSEnvironment("sandbox")
	APNSEnvironmentProduction APNSEnvironment = APNSEnvironment("production")
)

type SettingType string

const (
	SettingTypeServiceAccount                  SettingType = SettingType("service_account_file_path")
	SettingTypeAPNSJWTFilePath                             = SettingType("apns_jwt_file_path")
	SettingTypeAppleDeveloperTeamID                        = SettingType("apple_developer_team_id")
	SettingTypeAppleDeveloperKeyID                         = SettingType("apple_developer_key_id")
	SettingTypeAppleCertificateFilePath                    = SettingType("apple_certificate_file_path")
	SettingTypeAppleCertificateDecryptPassword             = SettingType("apple_certificate_decrypt_password")
	SettingTypeAPNSType                                    = SettingType("apns_type")
	SettingTypeAPNSLegacyEnvironment                       = SettingType("apns_mode")
	SettingTypeJWTAPNSEnvironment                          = SettingType("jwt_apns_mode")
	SettingTypeHistoryLength                               = SettingType("history_length")
	SettingTypeThemeMode                                   = SettingType("theme_mode")
)

type ThemeMode string

const (
	ThemeModeDark  ThemeMode = ThemeMode("dark")
	ThemeModeLight ThemeMode = ThemeMode("light")
)

type FCMDeviceType string

const (
	FCMDeviceTypeAndroid FCMDeviceType = FCMDeviceType("android")
	FCMDeviceTypeIos     FCMDeviceType = FCMDeviceType("ios")
	FCMDeviceTypeWeb     FCMDeviceType = FCMDeviceType("web")
)

var FCMDeviceTypes = []struct {
	Value  FCMDeviceType
	TSName string
}{
	{FCMDeviceTypeAndroid, "Android"},
	{FCMDeviceTypeIos, "iOS"},
	{FCMDeviceTypeWeb, "Web"},
}

// APNS Enums
type APNSPushType string

const (
	APNSPushTypeAlert        APNSPushType = APNSPushType("alert")
	APNSPushTypeBackground   APNSPushType = APNSPushType("background")
	APNSPushTypeLocation     APNSPushType = APNSPushType("location")
	APNSPushTypeVoIP         APNSPushType = APNSPushType("voip")
	APNSPushTypeComplication APNSPushType = APNSPushType("complication")
	APNSPushTypeFileProvider APNSPushType = APNSPushType("fileprovider")
	APNSPushTypeMDM          APNSPushType = APNSPushType("mdm")
)

var APNSPushTypes = []struct {
	Value  APNSPushType
	TSName string
}{
	{APNSPushTypeAlert, "Alert"},
	{APNSPushTypeBackground, "Background"},
	{APNSPushTypeLocation, "Location"},
	{APNSPushTypeVoIP, "VoIP"},
	{APNSPushTypeComplication, "Complication"},
	{APNSPushTypeFileProvider, "FileProvider"},
	{APNSPushTypeMDM, "MDM"},
}

type APNSPriority int

const (
	APNSPriorityImmediately APNSPriority = APNSPriority(10)
	APNSPriorityThrottled   APNSPriority = APNSPriority(5)
)

var APNSPriorities = []struct {
	Value  APNSPriority
	TSName string
}{
	{APNSPriorityImmediately, "Immediately"},
	{APNSPriorityThrottled, "Throttled"},
}

//Samples

type PayloadTemplate string

const (
	PayloadTemplateDefaultAPNS PayloadTemplate = PayloadTemplate(`{
  "aps": {
    "alert": {
      "title": "Notification Title",
      "subtitle": "Notification Subtitle",
      "body": "Notification Message Body"
    },
    "sound": "default",
    "badge": 19,
    "category": "CATEGORY_IDENTIFIER_SHOULD_BE_UPPERCASED",
    "thread-id": "message_thread_id",
    "mutable-content": 1
  },
  "custom_data": {
    "key": "value",
    "nested_object": {
      "key": "value"
    }
  }
}`)

	PayloadTemplateDefaultFCMAndroid PayloadTemplate = PayloadTemplate(`{
  "notification": {
    "title": "Urgent Update",
    "body": "Please update your app to the latest version immediately.",
    "image": "https://example.com/update.jpg"
  },
  "android": {
    "notification": {
      "icon": "https://example.com/icon.png"
    }
  },
  "data": {
    "updateVersion": "2.0.1",
    "customKey1": "customValue1",
    "customKey2": "customValue2"
  }
}`)
	PayloadTemplateDefaultFCMWeb PayloadTemplate = PayloadTemplate(`{
  "notification": {
    "title": "Urgent Update",
    "body": "Please update your app to the latest version immediately.",
    "icon": "https://example.com/icon.png",
    "image": "https://example.com/update.jpg",
    "actions": [
      {
        "action": "open_update",
        "title": "Update Now"
      },
      {
        "action": "remind_later",
        "title": "Remind Me Later"
      }
    ]
  },
  "webpush": {
    "notification": {
      "icon": "https://example.com/icon.png"
    }
  },
  "data": {
    "updateVersion": "2.0.1"
  }
}`)

	PayloadTemplateDefaultFCMiOS PayloadTemplate = PayloadTemplate(`{
  "aps": {
    "alert": {
      "title": "Notification Title",
      "subtitle": "Notification Subtitle",
      "body": "Notification Message Body"
    },
    "sound": "default",
    "badge": 19,
    "category": "CATEGORY_IDENTIFIER_SHOULD_BE_UPPERCASED",
    "thread-id": "message_thread_id",
    "mutable-content": 1
  },
  "custom_data": {
    "key": "value",
    "nested_object": {
      "key": "value"
    }
  }
}`)
)

var PayloadTemplates = []struct {
	Value  PayloadTemplate
	TSName string
}{
	{PayloadTemplateDefaultAPNS, "DefaultAPNS"},
	{PayloadTemplateDefaultFCMAndroid, "DefaultFCMAndroid"},
	{PayloadTemplateDefaultFCMWeb, "DefaultFCMWeb"},
	{PayloadTemplateDefaultFCMiOS, "DefaultFCMiOS"},
}

type Bundle string

const (
	AppName Bundle = "Remotify"
	Version Bundle = "1.0.2"
)

var AppBundle = []struct {
	Value  Bundle
	TSName string
}{
	{AppName, "AppName"},
	{Version, "Version"},
}
