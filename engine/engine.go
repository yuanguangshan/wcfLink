package engine

import (
	"context"
	"log/slog"

	internalapp "wcfLink/internal/app"
	"wcfLink/internal/config"
	"wcfLink/internal/httpapi"
	"wcfLink/internal/model"
)

type Config = config.Config
type LoginSession = model.LoginSession
type Account = model.Account
type Event = model.Event
type LogEntry = model.LogEntry
type Settings = model.Settings

type Engine struct {
	app *internalapp.App
}

func LoadConfig() Config {
	return config.Load()
}

func New(ctx context.Context, cfg Config, logger *slog.Logger) (*Engine, error) {
	app, err := internalapp.New(ctx, cfg, logger)
	if err != nil {
		return nil, err
	}
	return &Engine{app: app}, nil
}

func (e *Engine) Run(ctx context.Context) error {
	return e.app.Run(ctx)
}

func (e *Engine) StartBackground(ctx context.Context) error {
	return e.app.StartBackground(ctx)
}

func (e *Engine) Shutdown() error {
	return e.app.Shutdown()
}

func (e *Engine) StartLogin(ctx context.Context, baseURL string) (LoginSession, error) {
	return e.app.StartLogin(ctx, baseURL)
}

func (e *Engine) GetLoginStatus(ctx context.Context, sessionID string) (LoginSession, error) {
	return e.app.GetLoginStatus(ctx, sessionID)
}

func (e *Engine) GetLoginSession(ctx context.Context, sessionID string) (LoginSession, error) {
	return e.app.GetLoginSession(ctx, sessionID)
}

func (e *Engine) GetLoginQRCodePNG(ctx context.Context, sessionID string) ([]byte, error) {
	session, err := e.app.GetLoginSession(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return httpapi.GenerateQRCodePNG(session.QRCodeURL)
}

func (e *Engine) ListAccounts(ctx context.Context) ([]Account, error) {
	return e.app.ListAccounts(ctx)
}

func (e *Engine) ListEvents(ctx context.Context, afterID int64, limit int) ([]Event, error) {
	return e.app.ListEvents(ctx, afterID, limit)
}

func (e *Engine) GetSettings(ctx context.Context) (Settings, error) {
	return e.app.GetSettings(ctx)
}

func (e *Engine) UpdateSettings(ctx context.Context, settings Settings) (Settings, error) {
	return e.app.UpdateSettings(ctx, settings)
}

func (e *Engine) SendText(ctx context.Context, accountID, toUserID, text, contextToken string) error {
	return e.app.SendText(ctx, accountID, toUserID, text, contextToken)
}

func (e *Engine) SendMedia(ctx context.Context, accountID, toUserID, mediaType, filePath, text, contextToken string) error {
	return e.app.SendMedia(ctx, accountID, toUserID, mediaType, filePath, text, contextToken)
}

func (e *Engine) LogoutAccount(ctx context.Context, accountID string) error {
	return e.app.LogoutAccount(ctx, accountID)
}
