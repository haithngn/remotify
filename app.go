package main

import (
	"context"
	"net/http"
	"notification-deployer/internal/domain/usecases"
	"notification-deployer/internal/domain/values"
	"notification-deployer/internal/pkg/logging"
	"os"
	"path/filepath"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/wailsapp/mimetype"
	"github.com/wailsapp/wails/v2/pkg/menu"
	"github.com/wailsapp/wails/v2/pkg/menu/keys"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// App struct
type App struct {
	ctx context.Context
	db  *gorm.DB
}

// NewApp creates a new App application struct
func NewApp() *App {
	a := &App{}

	a.Launch()

	return a
}

func (a *App) Launch() {
	//Create a hidden app directory
	hiddenDir, err := usecases.MakeAppDirIfNotExist()
	if err != nil {
		log.Error("Failed to create app directory: ", err)
		return
	}

	err = logging.InitLogger(log.DebugLevel, hiddenDir)
	if err != nil {
		log.Error("Failed to initialize logger: ", err)
	}

	logLevel := logger.Silent
	if os.Getenv("REMOTIFY_LOG_LEVEL") == "debug" {
		logLevel = logger.Info
	}

	//Connect database
	sqlDB, err := gorm.Open(sqlite.Open(hiddenDir+"/remotify.db"), &gorm.Config{
		Logger: logger.New(
			log.StandardLogger(),
			logger.Config{
				SlowThreshold:             time.Second,
				LogLevel:                  logLevel,
				IgnoreRecordNotFoundError: false,
				ParameterizedQueries:      true,
				Colorful:                  false,
			},
		),
	})
	if err != nil {
		log.Error("Failed to connect database")
	}

	//Migrate database schema
	err = usecases.Migration(sqlDB)
	if err != nil {
		log.Error("Failed to migrate database: ", err)
	}

	a.db = sqlDB
}

// startup is called when the app starts. The context is saved,
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) domReady(ctx context.Context) {
	runtime.OnFileDrop(ctx, a.handleDroppedFile)
}

func (a *App) GetSettings() string {
	return usecases.GetSettings(a.db)
}

func (a *App) GetRecentAPNSMessages() string {
	return usecases.GetRecentAPNSMessages(a.db)
}

func (a *App) GetRecentFCMMessages() string {
	return usecases.GetRecentFCMMessages(a.db)
}

func (a *App) UploadAPNSToken() {
	f, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				Pattern: "*.p8",
			},
		},
	})

	if err != nil {
		return
	}

	usecases.UpdateAPNSToken(f, a.db)
}

func (a *App) UploadAPNSCertificate() string {
	f, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				Pattern: "*.p12",
			},
		},
	})

	if err != nil {
		return values.Failed(err, http.StatusBadRequest)
	}

	return usecases.UpdateAPNSCertificate(f, a.db)
}

func (a *App) UploadServiceAccount() string {
	f, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Filters: []runtime.FileFilter{
			{
				Pattern: "*.json",
			},
		},
	})

	if err != nil {
		return values.Failed(err, http.StatusBadRequest)
	}

	return usecases.UpdateServiceAccount(f, a.db)
}

func (a *App) SetKeyID(id string) string {
	return usecases.SetAppleDeveloperKeyID(id, a.db)
}

func (a *App) SetTeamID(id string) string {
	return usecases.SetAppleDeveloperTeamID(id, a.db)
}

func (a *App) SetDecryptPassword(password string) string {
	return usecases.SetAppleCertificateDecryptPassword(password, a.db)
}

func (a *App) SetLegacyAPNSEnvironment(mode values.APNSEnvironment) string {
	return usecases.SetLegacyAPNSEnvironment(mode, a.db)
}

func (a *App) SetJWTAPNSEnvironment(mode values.APNSEnvironment) string {
	return usecases.SetJWTAPNSEnvironment(mode, a.db)
}

func (a *App) SetAPNSType(mode values.APNSType) string {
	return usecases.SetAPNSType(mode, a.db)
}

func (a *App) SendAPNS(deviceToken, bundleId, payload, apnsID, collapseID, expiredAt string, priority values.APNSPriority, pushType values.APNSPushType, toSave bool) string {
	usecases.EmitLog(` {
		"deviceToken": "`+deviceToken+`",
		"bundleId": "`+bundleId+`",
		"payload": "`+payload+`",
		"apnsID": "`+apnsID+`",
		"collapseID": "`+collapseID+`",
		"expiredAt": "`+expiredAt+`",
		"priority": "`+strconv.Itoa(int(priority))+`",
		"pushType": "`+string(pushType)+`",
		"toSave": "`+strconv.FormatBool(toSave)+`"
	}`, usecases.LogDirectionSent, a.ctx)

	var expiredString *string
	if expiredAt != "" {
		expiredString = &expiredAt
	}
	var apnsIDString *string
	if apnsID != "" {
		apnsIDString = &apnsID
	}

	var collapseIDString *string
	if collapseID != "" {
		collapseIDString = &collapseID
	}

	result := usecases.SendAPNS(usecases.APNSMessageParams{
		DeviceToken: deviceToken,
		BundleID:    bundleId,
		Payload:     payload,
		APNSID:      apnsIDString,
		CollapseID:  collapseIDString,
		ExpiredAt:   expiredString,
		Priority:    priority,
		PushType:    pushType,
		ToSave:      toSave,
	}, a.db)

	usecases.EmitLog(result, usecases.LogDirectionReceived, a.ctx)

	return result
}

func (a *App) SendFCM(registrationToken string, deviceType values.FCMDeviceType, payloadData string, save bool) string {
	usecases.EmitLog(` {
		"registrationToken": "`+registrationToken+`",
		"deviceType": "`+string(deviceType)+`",
		"payloadData": "`+payloadData+`",
		"toSave": "`+strconv.FormatBool(save)+`"
	}`, usecases.LogDirectionSent, a.ctx)

	result := usecases.SendFCM(usecases.SendFCMMessageParams{
		DeviceToken: registrationToken,
		DeviceType:  deviceType,
		PayloadData: payloadData,
		ToSave:      save,
	}, a.db)

	usecases.EmitLog(result, usecases.LogDirectionReceived, a.ctx)

	return result
}

func (a *App) RemoveRecentAPNSMessage(id uint) interface{} {
	usecases.EmitLog("Removing APNS message...", usecases.LogDirectionSent, a.ctx)
	return usecases.RemoveRecentAPNSMessage(id, a.db)
}

func (a *App) RemoveRecentFCMMessage(id uint) interface{} {
	usecases.EmitLog("Removing FCM message...", usecases.LogDirectionSent, a.ctx)
	return usecases.RemoveRecentFCMMessage(id, a.db)
}

func (a *App) UpdateAPNSMessageNote(note string, id uint) interface{} {
	usecases.EmitLog("Updating APNS message note...", usecases.LogDirectionSent, a.ctx)
	return usecases.SaveAPNSMessageNote(note, id, a.db)
}

func (a *App) UpdateFCMMessageNote(note string, id uint) interface{} {
	usecases.EmitLog("Updating FCM message note...", usecases.LogDirectionSent, a.ctx)
	return usecases.SaveFCMMessageNote(note, id, a.db)
}

func (a *App) ToggleLightTheme() interface{} {
	usecases.EmitLog("Toggling light theme...", usecases.LogDirectionSent, a.ctx)
	usecases.ToggleThemeMode(values.ThemeModeLight, a.db)
	runtime.EventsEmit(a.ctx, "onChangeLightMode", "{}")

	return values.Succeed(nil)
}

func (a *App) ToggleDarkTheme() interface{} {
	usecases.EmitLog("Toggling dark theme...", usecases.LogDirectionSent, a.ctx)
	_ = usecases.ToggleThemeMode(values.ThemeModeDark, a.db)
	runtime.EventsEmit(a.ctx, "onChangeDarkMode", "{}")

	return values.Succeed(nil)
}

// App Menu
func (a *App) getMenus() *menu.Menu {
	AppMenu := menu.NewMenu()

	mainMenu := AppMenu.AddSubmenu("Remotify")
	mainMenu.AddText("About", nil, func(_ *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "onOpenAboutWindow", "{"+
			"\"app_version\": \"1.0.2\","+
			"\"copy_right\":\"© 2024 Hai Nguyen\","+
			"}")
	})
	mainMenu.AddSeparator()
	mainMenu.AddText("Preference", keys.CmdOrCtrl(","), func(_ *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "onOpenSettingWindow", "{}")
	})
	mainMenu.AddSeparator()
	mainMenu.AddText("Quit", keys.CmdOrCtrl("q"), func(_ *menu.CallbackData) {
		runtime.Quit(a.ctx)
	})

	AppMenu.Append(menu.EditMenu())

	//Selection
	selectionMenu := AppMenu.AddSubmenu("Selection")
	selectionMenu.AddText("Reset Payload", keys.CmdOrCtrl("r"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "onResetPayload", "{}")
	})

	viewMenu := AppMenu.AddSubmenu("View")
	viewMenu.AddText("Dark Mode", nil, func(_ *menu.CallbackData) {
		_ = a.ToggleDarkTheme()
	})
	viewMenu.AddText("Light Mode", nil, func(_ *menu.CallbackData) {
		_ = a.ToggleLightTheme()
	})

	goMenu := AppMenu.AddSubmenu("Go")
	goMenu.AddText("APNS", keys.CmdOrCtrl("1"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "onOpenAPNS", "{}")
	})
	goMenu.AddText("FCM", keys.CmdOrCtrl("2"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "onOpenFCM", "{}")
	})

	goMenu.AddText("Show Logs", keys.CmdOrCtrl("3"), func(_ *menu.CallbackData) {
		runtime.EventsEmit(a.ctx, "onOpenLogs", "{}")
	})

	return AppMenu
}

// RevealFileLocation locates the file and opens it in the default application
func (a *App) RevealFileLocation(file string) string {
	return usecases.RevealFileLocation(file)
}

func (a *App) handleDroppedFile(x, y int, files []string) {
	//Validate length of files to avoid panics
	if len(files) == 0 {
		return
	}

	//Validate length of files to avoid multiple files drop
	if len(files) > 1 {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Oops!",
			Message: "Multiple files are not supported",
		})
		return
	}

	file := files[0]

	//Detect file mime type
	mtype, err := mimetype.DetectFile(file)

	if err != nil {
		runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
			Type:    runtime.InfoDialog,
			Title:   "Oops!",
			Message: "Failed to read file",
		})
		return
	}

	switch mtype.String() {
	case "application/json":
		_ = usecases.UpdateServiceAccount(file, a.db)
	case "application/pkcs8":
		_ = usecases.UpdateAPNSToken(file, a.db)
	case "application/x-pkcs12":
		_ = usecases.UpdateAPNSCertificate(file, a.db)

	default:
		// Try with file extension
		ext := filepath.Ext(file)
		switch ext {
		case ".json":
			_ = usecases.UpdateServiceAccount(file, a.db)
		case ".p8":
			_ = usecases.UpdateAPNSToken(file, a.db)
		case ".p12":
			_ = usecases.UpdateAPNSCertificate(file, a.db)
		default:
			runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
				Type:    runtime.InfoDialog,
				Title:   "Oops!",
				Message: "Unsupported file type",
			})
			return
		}
	}

	runtime.EventsEmit(a.ctx, "onOpenSettingWindow", "{}")
}
