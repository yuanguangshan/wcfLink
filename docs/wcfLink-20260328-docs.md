# Project Documentation

- **Generated at:** 2026-03-28 09:13:59
- **Root Dir:** `.`
- **File Count:** 17
- **Total Lines:** 3026
- **Total Size:** 84.58 KB

<a name="toc"></a>
## 📂 扫描目录
- [.gitignore](#file-.gitignore) (13 lines, 0.10 KB)
- [cmd/wcfLink/main.go](#file-cmd-wcflink-main.go) (46 lines, 0.90 KB)
- [docs/test_report.md](#file-docs-test_report.md) (180 lines, 3.76 KB)
- [engine/doc.go](#file-engine-doc.go) (6 lines, 0.32 KB)
- [engine/engine.go](#file-engine-engine.go) (100 lines, 2.73 KB)
- [go.mod](#file-go.mod) (20 lines, 0.58 KB)
- [internal/app/app.go](#file-internal-app-app.go) (651 lines, 19.21 KB)
- [internal/app/runtime.go](#file-internal-app-runtime.go) (39 lines, 0.80 KB)
- [internal/config/config.go](#file-internal-config-config.go) (136 lines, 3.47 KB)
- [internal/httpapi/qr.go](#file-internal-httpapi-qr.go) (7 lines, 0.17 KB)
- [internal/httpapi/server.go](#file-internal-httpapi-server.go) (279 lines, 9.38 KB)
- [internal/ilink/client.go](#file-internal-ilink-client.go) (262 lines, 7.44 KB)
- [internal/ilink/media.go](#file-internal-ilink-media.go) (437 lines, 12.82 KB)
- [internal/model/types.go](#file-internal-model-types.go) (71 lines, 2.45 KB)
- [internal/store/store.go](#file-internal-store-store.go) (517 lines, 14.85 KB)
- [internal/worker/poller.go](#file-internal-worker-poller.go) (184 lines, 4.18 KB)
- [version/version.go](#file-version-version.go) (78 lines, 1.42 KB)

---

<a name="file-.gitignore"></a>
## 📄 .gitignore

````text
.*
!.gitignore
!.github/

bin/
dist/
data/
tmp/
build/
frontend/node_modules/
*.db
*.db-shm
*.db-wal

````

[⬆ 回到目录](#toc)

<a name="file-cmd-wcflink-main.go"></a>
## 📄 cmd/wcfLink/main.go

````go
package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/lich0821/wcfLink/internal/app"
	"github.com/lich0821/wcfLink/internal/config"
	coreversion "github.com/lich0821/wcfLink/version"
)

func main() {
	showVersion := flag.Bool("version", false, "print version and exit")
	flag.Parse()
	if *showVersion {
		fmt.Println(coreversion.String())
		return
	}

	cfg := config.Load()

	level := new(slog.LevelVar)
	level.Set(cfg.LogLevel())
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	a, err := app.New(ctx, cfg, logger)
	if err != nil {
		logger.Error("failed to initialize app", "err", err)
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		logger.Error("application exited with error", "err", err)
		os.Exit(1)
	}
}

````

[⬆ 回到目录](#toc)

<a name="file-docs-test_report.md"></a>
## 📄 docs/test_report.md

````markdown
# wcfLink 测试报告

**测试时间**: 2026-03-28
**测试环境**: macOS Darwin 21.6.0
**Go 版本**: 1.25.4

---

## 一、代码质量测试

| 测试项 | 状态 | 说明 |
|--------|------|------|
| `go vet ./...` | ✅ 通过 | 无代码问题 |
| `go build ./...` | ✅ 通过 | 编译成功 |
| `go mod verify` | ✅ 通过 | 依赖验证通过 |

---

## 二、二进制构建测试

```bash
go build -o ./bin/wcfLink ./cmd/wcfLink
./bin/wcfLink -version
```

**输出**:
```
v0.0.0-20260325125511-fb0999b81043 (fb0999b81043, built 2026-03-25T12:55:11Z)
```

状态: ✅ 成功

---

## 三、HTTP API 基础测试

### 3.1 健康检查

| 接口 | 状态 | 响应 |
|------|------|------|
| `GET /health/live` | ✅ | `{"ok":true}` |
| `GET /health/ready` | ✅ | `{"ok":true}` |
| `GET /api/version` | ✅ | 返回版本信息 |
| `GET /api/accounts` | ✅ | 正常响应 |
| `GET /api/settings` | ✅ | 返回配置信息 |
| `GET /api/events` | ✅ | 正常响应 |

---

## 四、登录流程测试

### 4.1 发起登录

```bash
POST /api/accounts/login/start
```

**响应**:
- session_id: `login_1774659812139187000`
- qr_code_url: `https://liteapp.weixin.qq.com/q/7GiQu1?qrcode=...`

### 4.2 扫码确认

| 字段 | 值 |
|------|-----|
| 状态 | `confirmed` ✅ |
| account_id | `93a4696a9846@im.bot` |
| ilink_user_id | `o9cq80xaZYxScOoGuSy-5fGnhz9o@im.wechat` |

### 4.3 账号持久化

| 字段 | 值 |
|------|-----|
| account_id | `93a4696a9846@im.bot` |
| ilink_user_id | `o9cq80xaZYxScOoGuSy-5fGnhz9o@im.wechat` |
| login_status | `connected` |
| enabled | `true` |

状态: ✅ 登录成功，账号已持久化

---

## 五、消息收发测试

### 5.1 接收消息 (inbound)

| 字段 | 值 |
|------|-----|
| 方向 | `inbound` |
| 类型 | `text` |
| 发送者 | `o9cq80xaZYxScOoGuSy-5fGnhz9o@im.wechat` |
| 内容 | "全部测试，这是来自微信的命令" |
| context_token | 已自动保存 ✅ |

状态: ✅ 通过

### 5.2 发送文本消息 (outbound)

```bash
POST /api/messages/send-text
{
  "account_id": "93a4696a9846@im.bot",
  "to_user_id": "o9cq80xaZYxScOoGuSy-5fGnhz9o@im.wechat",
  "text": "你好！这是来自 wcfLink 的测试消息"
}
```

响应: `{"ok":true}` ✅

### 5.3 发送图片 (outbound)

```bash
POST /api/messages/send-media
{
  "account_id": "93a4696a9846@im.bot",
  "to_user_id": "o9cq80xaZYxScOoGuSy-5fGnhz9o@im.wechat",
  "type": "image",
  "file_path": "/tmp/test_image.png",
  "text": "这是一张测试图片"
}
```

响应: `{"ok":true}` ✅

### 5.4 发送文件 (outbound)

```bash
POST /api/messages/send-media
{
  "account_id": "93a4696a9846@im.bot",
  "to_user_id": "o9cq80xaZYxScOoGuSy-5fGnhz9o@im.wechat",
  "type": "file",
  "file_path": "/tmp/test_file.txt",
  "text": "这是一个测试文件"
}
```

响应: `{"ok":true}` ✅

---

## 六、事件记录

| ID | 方向 | 类型 | 内容 |
|----|------|------|------|
| 1 | 📥 inbound | text | "全部测试，这是来自微信的命令" |
| 2 | 📤 outbound | text | "你好！这是来自 wcfLink 的测试消息" |
| 3 | 📤 outbound | image | "这是一张测试图片" + test_image.png |
| 4 | 📤 outbound | file | "这是一个测试文件" + test_file.txt |

所有事件已保存到 SQLite 数据库 ✅

---

## 七、测试总结

| 功能模块 | 状态 |
|----------|------|
| 代码编译 | ✅ 通过 |
| 模块验证 | ✅ 通过 |
| 二进制构建 | ✅ 通过 |
| HTTP 服务启动 | ✅ 通过 |
| 扫码登录 | ✅ 通过 |
| 账号持久化 | ✅ 通过 |
| 长轮询接收消息 | ✅ 通过 |
| 发送文本消息 | ✅ 通过 |
| 发送图片 | ✅ 通过 |
| 发送文件 | ✅ 通过 |
| 事件存储 | ✅ 通过 |
| context_token 管理 | ✅ 通过 |

---

## 结论

**wcfLink 项目所有核心功能测试通过！** 🎉

---

*报告生成时间: 2026-03-28 09:08*

````

[⬆ 回到目录](#toc)

<a name="file-engine-doc.go"></a>
## 📄 engine/doc.go

````go
// Package engine exposes the reusable wcfLink service layer for Go applications.
//
// The current implementation is a thin public wrapper over the existing internal
// app/service runtime so GUI, HTTP, and future CLI integrations can depend on a
// stable package boundary before the deeper refactor is complete.
package engine

````

[⬆ 回到目录](#toc)

<a name="file-engine-engine.go"></a>
## 📄 engine/engine.go

````go
package engine

import (
	"context"
	"log/slog"

	internalapp "github.com/lich0821/wcfLink/internal/app"
	"github.com/lich0821/wcfLink/internal/config"
	"github.com/lich0821/wcfLink/internal/httpapi"
	"github.com/lich0821/wcfLink/internal/model"
	coreversion "github.com/lich0821/wcfLink/version"
)

type Config = config.Config
type LoginSession = model.LoginSession
type Account = model.Account
type Event = model.Event
type LogEntry = model.LogEntry
type Settings = model.Settings
type VersionInfo = coreversion.Info

type Engine struct {
	app *internalapp.App
}

func LoadConfig() Config {
	return config.Load()
}

func CurrentVersion() VersionInfo {
	return coreversion.Current()
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

````

[⬆ 回到目录](#toc)

<a name="file-go.mod"></a>
## 📄 go.mod

````text
module github.com/lich0821/wcfLink

go 1.25.4

require (
	github.com/skip2/go-qrcode v0.0.0-20200617195104-da1b6568686e
	modernc.org/sqlite v1.47.0
)

require (
	github.com/dustin/go-humanize v1.0.1 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/ncruces/go-strftime v1.0.0 // indirect
	github.com/remyoudompheng/bigfft v0.0.0-20230129092748-24d4a6f8daec // indirect
	golang.org/x/sys v0.42.0 // indirect
	modernc.org/libc v1.70.0 // indirect
	modernc.org/mathutil v1.7.1 // indirect
	modernc.org/memory v1.11.0 // indirect
)

````

[⬆ 回到目录](#toc)

<a name="file-internal-app-app.go"></a>
## 📄 internal/app/app.go

````go
package app

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lich0821/wcfLink/internal/config"
	"github.com/lich0821/wcfLink/internal/httpapi"
	"github.com/lich0821/wcfLink/internal/ilink"
	"github.com/lich0821/wcfLink/internal/model"
	"github.com/lich0821/wcfLink/internal/store"
	"github.com/lich0821/wcfLink/internal/worker"
)

type App struct {
	cfg     config.Config
	logger  *slog.Logger
	store   *store.Store
	client  *ilink.Client
	pollers *worker.PollerManager
	server  *http.Server
	runtime *runtimeState
	svc     *service
}

func New(ctx context.Context, cfg config.Config, logger *slog.Logger) (*App, error) {
	st, err := store.New(ctx, cfg.DBPath)
	if err != nil {
		return nil, err
	}
	client := ilink.NewClient(cfg.ChannelVersion, cfg.PollTimeout+10*time.Second)
	runtime := newRuntimeState(st, cfg)
	svc := &service{
		cfg:     cfg,
		logger:  logger,
		store:   st,
		client:  client,
		runtime: runtime,
	}
	pollers := worker.NewPollerManager(st, client, logger, svc.HandleInboundMessage)
	svc.pollers = pollers
	api := httpapi.NewServer(&service{
		cfg:     cfg,
		logger:  logger,
		store:   st,
		client:  client,
		pollers: pollers,
		runtime: runtime,
	}, logger)

	return &App{
		cfg:     cfg,
		logger:  logger,
		store:   st,
		client:  client,
		pollers: pollers,
		runtime: runtime,
		svc:     svc,
		server: &http.Server{
			Addr:              cfg.ListenAddr,
			Handler:           api.Handler(),
			ReadHeaderTimeout: 10 * time.Second,
		},
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	if err := a.StartBackground(ctx); err != nil {
		return err
	}

	<-ctx.Done()
	a.logger.Info("shutdown requested")
	return a.Shutdown()
}

func (a *App) StartBackground(ctx context.Context) error {
	if err := a.store.Ping(ctx); err != nil {
		return err
	}
	if err := a.pollers.StartEnabledAccounts(ctx); err != nil {
		return err
	}

	go func() {
		a.logger.Info("http server listening", "addr", a.cfg.ListenAddr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("http server stopped with error", "err", err)
		}
	}()
	return nil
}

func (a *App) Shutdown() error {
	a.pollers.StopAll()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = a.server.Shutdown(shutdownCtx)
	return a.store.Close()
}

func (a *App) StartLogin(ctx context.Context, baseURL string) (model.LoginSession, error) {
	return a.svc.StartLogin(ctx, baseURL)
}

func (a *App) GetLoginStatus(ctx context.Context, sessionID string) (model.LoginSession, error) {
	return a.svc.GetLoginStatus(ctx, sessionID)
}

func (a *App) GetLoginSession(ctx context.Context, sessionID string) (model.LoginSession, error) {
	return a.svc.GetLoginSession(ctx, sessionID)
}

func (a *App) ListAccounts(ctx context.Context) ([]model.Account, error) {
	return a.svc.ListAccounts(ctx)
}

func (a *App) ListEvents(ctx context.Context, afterID int64, limit int) ([]model.Event, error) {
	return a.svc.ListEvents(ctx, afterID, limit)
}

func (a *App) GetSettings(ctx context.Context) (model.Settings, error) {
	return a.svc.GetSettings(ctx)
}

func (a *App) UpdateSettings(ctx context.Context, settings model.Settings) (model.Settings, error) {
	return a.svc.UpdateSettings(ctx, settings)
}

func (a *App) SendText(ctx context.Context, accountID, toUserID, text, contextToken string) error {
	return a.svc.SendText(ctx, accountID, toUserID, text, contextToken)
}

func (a *App) SendMedia(ctx context.Context, accountID, toUserID, mediaType, filePath, text, contextToken string) error {
	return a.svc.SendMedia(ctx, accountID, toUserID, mediaType, filePath, text, contextToken)
}

func (a *App) LogoutAccount(ctx context.Context, accountID string) error {
	return a.svc.LogoutAccount(ctx, accountID)
}

type service struct {
	cfg     config.Config
	logger  *slog.Logger
	store   *store.Store
	client  *ilink.Client
	pollers *worker.PollerManager
	runtime *runtimeState
}

func (s *service) StartLogin(ctx context.Context, baseURL string) (model.LoginSession, error) {
	if baseURL == "" {
		baseURL = s.cfg.DefaultBaseURL
	}
	qr, err := s.client.FetchQRCode(ctx, baseURL)
	if err != nil {
		return model.LoginSession{}, err
	}
	now := time.Now().UTC()
	session := model.LoginSession{
		SessionID: fmt.Sprintf("login_%d", now.UnixNano()),
		BaseURL:   baseURL,
		QRCode:    qr.QRCode,
		QRCodeURL: qr.QRCodeURL,
		Status:    "wait",
		StartedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateLoginSession(ctx, session); err != nil {
		return model.LoginSession{}, err
	}
	return session, nil
}

func (s *service) GetLoginStatus(ctx context.Context, sessionID string) (model.LoginSession, error) {
	session, err := s.store.GetLoginSession(ctx, sessionID)
	if err != nil {
		return model.LoginSession{}, err
	}
	if session.Status == "confirmed" {
		return session, nil
	}
	status, err := s.client.FetchQRCodeStatus(ctx, session.BaseURL, session.QRCode)
	if err != nil {
		_ = s.store.UpdateLoginSessionStatus(context.Background(), sessionID, session.Status, err.Error())
		return model.LoginSession{}, err
	}
	if status.Status == "" {
		return session, nil
	}
	if status.Status == "confirmed" {
		if err := s.store.CompleteLoginSession(ctx, sessionID, status); err != nil {
			return model.LoginSession{}, err
		}
		account, err := s.store.GetAccount(ctx, status.AccountID)
		if err == nil {
			s.pollers.StartAccount(context.Background(), account)
		}
		return s.store.GetLoginSession(ctx, sessionID)
	}
	if err := s.store.UpdateLoginSessionStatus(ctx, sessionID, status.Status, ""); err != nil {
		return model.LoginSession{}, err
	}
	return s.store.GetLoginSession(ctx, sessionID)
}

func (s *service) GetLoginSession(ctx context.Context, sessionID string) (model.LoginSession, error) {
	return s.store.GetLoginSession(ctx, sessionID)
}

func (s *service) ListAccounts(ctx context.Context) ([]model.Account, error) {
	return s.store.ListAccounts(ctx)
}

func (s *service) ListEvents(ctx context.Context, afterID int64, limit int) ([]model.Event, error) {
	return s.store.ListEvents(ctx, afterID, limit)
}

func (s *service) ListLogs(ctx context.Context, afterID int64, limit int) ([]model.LogEntry, error) {
	return s.store.ListLogs(ctx, afterID, limit)
}

func (s *service) GetSettings(ctx context.Context) (model.Settings, error) {
	_ = ctx
	return s.runtime.Settings(), nil
}

func (s *service) UpdateSettings(ctx context.Context, settings model.Settings) (model.Settings, error) {
	if err := config.SaveFileSettings(s.cfg.SettingsPath, config.FileSettings{
		ListenAddr: settings.ListenAddr,
		WebhookURL: settings.WebhookURL,
	}); err != nil {
		return model.Settings{}, err
	}
	if err := s.runtime.UpdateSettings(ctx, settings); err != nil {
		return model.Settings{}, err
	}
	return s.runtime.Settings(), nil
}

func (s *service) LogoutAccount(ctx context.Context, accountID string) error {
	s.pollers.StopAccount(accountID)
	if err := s.store.DeleteAccount(ctx, accountID); err != nil {
		return err
	}
	_ = s.store.AddLog(context.Background(), "INFO", "account disconnected locally", "account", fmt.Sprintf(`{"account_id":%q}`, accountID))
	return nil
}

func (s *service) SendText(ctx context.Context, accountID, toUserID, text, contextToken string) error {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}
	contextToken, err = s.resolveContextToken(ctx, accountID, toUserID, contextToken)
	if err != nil {
		return err
	}
	if strings.TrimSpace(contextToken) == "" {
		return errors.New("context token not found for this user; current text sending only supports replying to users who have already sent a message")
	}
	if err := s.client.SendTextMessage(ctx, account.BaseURL, account.Token, toUserID, text, contextToken); err != nil {
		_ = s.store.AddLog(context.Background(), "ERROR", "outbound send failed", "message", fmt.Sprintf(`{"account_id":%q,"to_user_id":%q,"err":%q}`, accountID, toUserID, err.Error()))
		return err
	}
	raw := fmt.Sprintf(`{"to_user_id":%q,"text":%q,"context_token":%q}`, toUserID, text, contextToken)
	if err := s.store.CreateOutboundEvent(ctx, accountID, "text", toUserID, contextToken, text, "", "", "", raw); err != nil {
		return err
	}
	_ = s.store.AddLog(context.Background(), "INFO", "outbound text sent", "message", raw)
	return nil
}

func (s *service) SendMedia(ctx context.Context, accountID, toUserID, mediaType, filePath, text, contextToken string) error {
	account, err := s.store.GetAccount(ctx, accountID)
	if err != nil {
		return err
	}
	contextToken, err = s.resolveContextToken(ctx, accountID, toUserID, contextToken)
	if err != nil {
		return err
	}
	if strings.TrimSpace(contextToken) == "" {
		return errors.New("context token not found for this user; media sending only supports replying to users who have already sent a message")
	}

	normalizedType, uploadType, err := normalizeMediaSendType(mediaType, filePath)
	if err != nil {
		return err
	}
	uploaded, err := s.client.UploadLocalMedia(ctx, s.cfg.CDNBaseURL, account.BaseURL, account.Token, toUserID, filePath, uploadType)
	if err != nil {
		_ = s.store.AddLog(context.Background(), "ERROR", "media upload failed", "message", fmt.Sprintf(`{"account_id":%q,"to_user_id":%q,"file_path":%q,"err":%q}`, accountID, toUserID, filePath, err.Error()))
		return err
	}

	fileName := filepath.Base(filePath)
	switch normalizedType {
	case "image":
		err = s.client.SendImageMessage(ctx, account.BaseURL, account.Token, toUserID, contextToken, text, uploaded)
	case "video":
		err = s.client.SendVideoMessage(ctx, account.BaseURL, account.Token, toUserID, contextToken, text, uploaded)
	case "file":
		err = s.client.SendFileMessage(ctx, account.BaseURL, account.Token, toUserID, contextToken, text, fileName, uploaded)
	case "voice":
		err = s.client.SendVoiceMessage(ctx, account.BaseURL, account.Token, toUserID, contextToken, text, detectVoiceEncodeType(filePath), uploaded)
	default:
		err = fmt.Errorf("unsupported media type %q", normalizedType)
	}
	if err != nil {
		_ = s.store.AddLog(context.Background(), "ERROR", "outbound media send failed", "message", fmt.Sprintf(`{"account_id":%q,"to_user_id":%q,"file_path":%q,"media_type":%q,"err":%q}`, accountID, toUserID, filePath, normalizedType, err.Error()))
		if isAudioFilePath(filePath) {
			return fmt.Errorf("%s发送失败：\n%s", strings.TrimPrefix(strings.ToLower(filepath.Ext(filePath)), "."), err.Error())
		}
		return err
	}

	mimeType := detectOutboundMIME(normalizedType, filePath)
	raw := fmt.Sprintf(`{"to_user_id":%q,"file_path":%q,"media_type":%q,"text":%q,"context_token":%q}`, toUserID, filePath, normalizedType, text, contextToken)
	if err := s.store.CreateOutboundEvent(ctx, accountID, normalizedType, toUserID, contextToken, text, filePath, fileName, mimeType, raw); err != nil {
		return err
	}
	_ = s.store.AddLog(context.Background(), "INFO", "outbound media sent", "message", raw)
	return nil
}

func (s *service) HandleInboundMessage(ctx context.Context, account model.Account, msg ilink.WeixinMessage) error {
	mediaPath, mediaFileName, mediaMimeType := "", "", ""
	if mediaItem, ok := firstInboundMediaItem(msg); ok {
		mediaBytes, suggestedFileName, mimeType, err := s.client.DownloadMessageMedia(ctx, s.cfg.CDNBaseURL, mediaItem)
		if err != nil {
			s.logger.Warn("download inbound media failed", "account_id", account.AccountID, "message_id", msg.MessageID, "err", err)
			_ = s.store.AddLog(context.Background(), "ERROR", "download inbound media failed", "media", fmt.Sprintf(`{"account_id":%q,"message_id":%d,"err":%q}`, account.AccountID, msg.MessageID, err.Error()))
		} else {
			mediaPath, mediaFileName, mediaMimeType, err = s.saveInboundMedia(account.AccountID, msg.MessageID, msg.FromUserID, suggestedFileName, mimeType, mediaBytes)
			if err != nil {
				s.logger.Warn("persist inbound media failed", "account_id", account.AccountID, "message_id", msg.MessageID, "err", err)
				_ = s.store.AddLog(context.Background(), "ERROR", "persist inbound media failed", "media", fmt.Sprintf(`{"account_id":%q,"message_id":%d,"err":%q}`, account.AccountID, msg.MessageID, err.Error()))
				mediaPath, mediaFileName, mediaMimeType = "", "", ""
			}
		}
	}

	if err := s.store.SaveInboundMessage(ctx, account.AccountID, msg, mediaPath, mediaFileName, mediaMimeType); err != nil {
		return err
	}

	payload, err := json.Marshal(map[string]any{
		"account_id":      account.AccountID,
		"base_url":        account.BaseURL,
		"event_type":      detectEventType(msg),
		"body_text":       extractBodyText(msg),
		"from_user_id":    msg.FromUserID,
		"to_user_id":      msg.ToUserID,
		"message_id":      msg.MessageID,
		"context_token":   msg.ContextToken,
		"media_path":      mediaPath,
		"media_file_name": mediaFileName,
		"media_mime_type": mediaMimeType,
		"raw_message":     msg,
		"received_at":     time.Now().UTC(),
	})
	if err != nil {
		return err
	}

	settings := s.runtime.Settings()
	if settings.WebhookURL == "" {
		text := "[non-text]"
		for _, item := range msg.ItemList {
			if item.Type == 1 && item.TextItem != nil && item.TextItem.Text != "" {
				text = item.TextItem.Text
				break
			}
			if item.Type == 3 && item.VoiceItem != nil && item.VoiceItem.Text != "" {
				text = item.VoiceItem.Text
				break
			}
		}
		return s.store.AddLog(ctx, "INFO", fmt.Sprintf("inbound message from %s: %s", msg.FromUserID, text), "inbound", string(payload))
	}

	go s.deliverWebhook(settings.WebhookURL, payload)
	return s.store.AddLog(ctx, "INFO", "inbound message queued for webhook", "webhook", string(payload))
}

func (s *service) resolveContextToken(ctx context.Context, accountID, toUserID, contextToken string) (string, error) {
	if strings.TrimSpace(contextToken) != "" {
		return contextToken, nil
	}
	peerCtx, err := s.store.GetPeerContext(ctx, accountID, toUserID)
	if err == nil {
		return peerCtx.ContextToken, nil
	}
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	return "", err
}

func (s *service) saveInboundMedia(accountID string, messageID int64, fromUserID, fileName, mimeType string, data []byte) (string, string, string, error) {
	if err := os.MkdirAll(s.cfg.MediaDir, 0o755); err != nil {
		return "", "", "", err
	}
	now := time.Now()
	dir := filepath.Join(s.cfg.MediaDir, sanitizePathSegment(accountID), now.Format("2006"), now.Format("01"), now.Format("02"))
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", "", "", err
	}

	safeName := sanitizeFileName(fileName)
	if safeName == "" {
		safeName = "media"
	}
	base := strings.TrimSuffix(safeName, filepath.Ext(safeName))
	ext := filepath.Ext(safeName)
	if ext == "" {
		ext = extensionForMIME(mimeType)
	}
	if ext == "" {
		ext = ".bin"
	}

	prefix := fmt.Sprintf("%d", messageID)
	if messageID == 0 {
		prefix = fmt.Sprintf("%d", now.UnixNano())
	}
	if fromUserID != "" {
		prefix += "_" + sanitizePathSegment(fromUserID)
	}
	finalName := prefix + "_" + base + ext
	fullPath := filepath.Join(dir, finalName)
	if err := os.WriteFile(fullPath, data, 0o644); err != nil {
		return "", "", "", err
	}
	return fullPath, finalName, mimeType, nil
}

func normalizeMediaSendType(mediaType, filePath string) (string, int, error) {
	value := strings.ToLower(strings.TrimSpace(mediaType))
	if value == "" {
		switch strings.ToLower(filepath.Ext(filePath)) {
		case ".jpg", ".jpeg", ".png", ".gif", ".webp":
			value = "image"
		case ".mp4", ".mov", ".m4v":
			value = "video"
		default:
			value = "file"
		}
	}
	switch value {
	case "image":
		return value, ilink.UploadMediaTypeImage, nil
	case "video":
		return value, ilink.UploadMediaTypeVideo, nil
	case "file":
		return value, ilink.UploadMediaTypeFile, nil
	case "voice":
		return value, ilink.UploadMediaTypeVoice, nil
	default:
		return "", 0, fmt.Errorf("unsupported media type %q", mediaType)
	}
}

func firstInboundMediaItem(msg ilink.WeixinMessage) (ilink.MessageItem, bool) {
	for _, item := range msg.ItemList {
		switch item.Type {
		case 2, 3, 4, 5:
			return item, true
		}
	}
	return ilink.MessageItem{}, false
}

func detectOutboundMIME(mediaType, filePath string) string {
	switch mediaType {
	case "image":
		switch strings.ToLower(filepath.Ext(filePath)) {
		case ".png":
			return "image/png"
		case ".gif":
			return "image/gif"
		default:
			return "image/jpeg"
		}
	case "video":
		return "video/mp4"
	case "voice":
		switch strings.ToLower(filepath.Ext(filePath)) {
		case ".amr":
			return "audio/amr"
		case ".mp3":
			return "audio/mpeg"
		case ".ogg":
			return "audio/ogg"
		default:
			return "audio/silk"
		}
	default:
		return "application/octet-stream"
	}
}

func sanitizePathSegment(value string) string {
	replacer := strings.NewReplacer(
		"/", "_",
		"\\", "_",
		":", "_",
		"*", "_",
		"?", "_",
		"\"", "_",
		"<", "_",
		">", "_",
		"|", "_",
		"@", "_",
	)
	out := strings.TrimSpace(replacer.Replace(value))
	if out == "" {
		return "unknown"
	}
	return out
}

func sanitizeFileName(value string) string {
	value = filepath.Base(strings.TrimSpace(value))
	if value == "." || value == "/" || value == "" {
		return ""
	}
	return sanitizePathSegment(value)
}

func extensionForMIME(mimeType string) string {
	switch mimeType {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "video/mp4":
		return ".mp4"
	case "audio/silk":
		return ".silk"
	case "audio/amr":
		return ".amr"
	case "audio/mpeg":
		return ".mp3"
	case "audio/ogg":
		return ".ogg"
	case "application/pdf":
		return ".pdf"
	default:
		return ""
	}
}

func detectVoiceEncodeType(filePath string) int {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".amr":
		return 5
	case ".mp3":
		return 7
	case ".ogg":
		return 8
	default:
		return 6
	}
}

func isAudioFilePath(filePath string) bool {
	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".silk", ".amr", ".mp3", ".ogg", ".wav", ".m4a":
		return true
	default:
		return false
	}
}

func extractBodyText(msg ilink.WeixinMessage) string {
	for _, item := range msg.ItemList {
		switch item.Type {
		case 1:
			if item.TextItem != nil {
				return item.TextItem.Text
			}
		case 3:
			if item.VoiceItem != nil && item.VoiceItem.Text != "" {
				return item.VoiceItem.Text
			}
		case 2:
			return "[image]"
		case 4:
			if item.FileItem != nil && item.FileItem.FileName != "" {
				return "[file] " + item.FileItem.FileName
			}
			return "[file]"
		case 5:
			return "[video]"
		}
	}
	return ""
}

func detectEventType(msg ilink.WeixinMessage) string {
	for _, item := range msg.ItemList {
		switch item.Type {
		case 1:
			return "text"
		case 2:
			return "image"
		case 3:
			return "voice"
		case 4:
			return "file"
		case 5:
			return "video"
		}
	}
	return "unknown"
}

func (s *service) deliverWebhook(webhookURL string, payload []byte) {
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, webhookURL, bytes.NewReader(payload))
	if err != nil {
		_ = s.store.AddLog(context.Background(), "ERROR", "build webhook request failed", "webhook", err.Error())
		return
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		_ = s.store.AddLog(context.Background(), "ERROR", "webhook delivery failed", "webhook", err.Error())
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		_ = s.store.AddLog(context.Background(), "ERROR", fmt.Sprintf("webhook delivery failed with status %d", resp.StatusCode), "webhook", "")
		return
	}
	_ = s.store.AddLog(context.Background(), "INFO", "webhook delivered", "webhook", "")
}

````

[⬆ 回到目录](#toc)

<a name="file-internal-app-runtime.go"></a>
## 📄 internal/app/runtime.go

````go
package app

import (
	"context"
	"sync"

	"github.com/lich0821/wcfLink/internal/config"
	"github.com/lich0821/wcfLink/internal/model"
	"github.com/lich0821/wcfLink/internal/store"
)

type runtimeState struct {
	mu       sync.RWMutex
	settings model.Settings
	store    *store.Store
}

func newRuntimeState(st *store.Store, cfg config.Config) *runtimeState {
	return &runtimeState{
		settings: model.Settings{
			ListenAddr: cfg.ListenAddr,
			WebhookURL: cfg.WebhookURL,
		},
		store: st,
	}
}

func (r *runtimeState) Settings() model.Settings {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.settings
}

func (r *runtimeState) UpdateSettings(ctx context.Context, settings model.Settings) error {
	r.mu.Lock()
	r.settings = settings
	r.mu.Unlock()
	return r.store.AddLog(ctx, "INFO", "settings updated", "settings", "")
}

````

[⬆ 回到目录](#toc)

<a name="file-internal-config-config.go"></a>
## 📄 internal/config/config.go

````go
package config

import (
	"encoding/json"
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const (
	defaultListenAddr     = "127.0.0.1:17890"
	defaultBaseURL        = "https://ilinkai.weixin.qq.com"
	defaultCDNBaseURL     = "https://novac2c.cdn.weixin.qq.com/c2c"
	defaultChannelVersion = "2.0.1"
	defaultPollTimeout    = 35 * time.Second
)

type Config struct {
	ListenAddr     string
	StateDir       string
	MediaDir       string
	DBPath         string
	SettingsPath   string
	DefaultBaseURL string
	CDNBaseURL     string
	ChannelVersion string
	PollTimeout    time.Duration
	LogLevelText   string
	OpenBrowser    bool
	WebhookURL     string
}

func Load() Config {
	stateDir := envOrDefault("WCFLINK_STATE_DIR", defaultStateDir())
	mediaDir := envOrDefault("WCFLINK_MEDIA_DIR", filepath.Join(stateDir, "media"))
	dbPath := envOrDefault("WCFLINK_DB_PATH", filepath.Join(stateDir, "wcfLink.db"))
	settingsPath := filepath.Join(stateDir, "settings.json")
	fileSettings := loadFileSettings(settingsPath)
	return Config{
		ListenAddr:     envOrDefault("WCFLINK_LISTEN_ADDR", valueOrDefault(fileSettings.ListenAddr, defaultListenAddr)),
		StateDir:       stateDir,
		MediaDir:       mediaDir,
		DBPath:         dbPath,
		SettingsPath:   settingsPath,
		DefaultBaseURL: envOrDefault("WCFLINK_BASE_URL", defaultBaseURL),
		CDNBaseURL:     envOrDefault("WCFLINK_CDN_BASE_URL", defaultCDNBaseURL),
		ChannelVersion: envOrDefault("WCFLINK_CHANNEL_VERSION", defaultChannelVersion),
		PollTimeout:    envDurationOrDefault("WCFLINK_POLL_TIMEOUT", defaultPollTimeout),
		LogLevelText:   envOrDefault("WCFLINK_LOG_LEVEL", "INFO"),
		OpenBrowser:    envBoolOrDefault("WCFLINK_OPEN_BROWSER", false),
		WebhookURL:     envOrDefault("WCFLINK_WEBHOOK_URL", fileSettings.WebhookURL),
	}
}

type FileSettings struct {
	ListenAddr string `json:"listen_addr"`
	WebhookURL string `json:"webhook_url"`
}

func SaveFileSettings(path string, settings FileSettings) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func (c Config) LogLevel() slog.Level {
	switch c.LogLevelText {
	case "DEBUG", "debug":
		return slog.LevelDebug
	case "WARN", "warn", "WARNING", "warning":
		return slog.LevelWarn
	case "ERROR", "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func defaultStateDir() string {
	exePath, err := os.Executable()
	if err == nil && exePath != "" {
		return filepath.Join(filepath.Dir(exePath), "data")
	}
	return filepath.Join(".", "data")
}

func envOrDefault(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envDurationOrDefault(key string, fallback time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func envBoolOrDefault(key string, fallback bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func loadFileSettings(path string) FileSettings {
	data, err := os.ReadFile(path)
	if err != nil {
		return FileSettings{}
	}
	var settings FileSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return FileSettings{}
	}
	return settings
}

func valueOrDefault(value, fallback string) string {
	if value != "" {
		return value
	}
	return fallback
}

````

[⬆ 回到目录](#toc)

<a name="file-internal-httpapi-qr.go"></a>
## 📄 internal/httpapi/qr.go

````go
package httpapi

import qrcode "github.com/skip2/go-qrcode"

func GenerateQRCodePNG(content string) ([]byte, error) {
	return qrcode.Encode(content, qrcode.Medium, 320)
}

````

[⬆ 回到目录](#toc)

<a name="file-internal-httpapi-server.go"></a>
## 📄 internal/httpapi/server.go

````go
package httpapi

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/lich0821/wcfLink/internal/model"
	coreversion "github.com/lich0821/wcfLink/version"
)

type Service interface {
	StartLogin(ctx context.Context, baseURL string) (model.LoginSession, error)
	GetLoginSession(ctx context.Context, sessionID string) (model.LoginSession, error)
	GetLoginStatus(ctx context.Context, sessionID string) (model.LoginSession, error)
	ListAccounts(ctx context.Context) ([]model.Account, error)
	ListEvents(ctx context.Context, afterID int64, limit int) ([]model.Event, error)
	ListLogs(ctx context.Context, afterID int64, limit int) ([]model.LogEntry, error)
	GetSettings(ctx context.Context) (model.Settings, error)
	UpdateSettings(ctx context.Context, settings model.Settings) (model.Settings, error)
	SendText(ctx context.Context, accountID, toUserID, text, contextToken string) error
	SendMedia(ctx context.Context, accountID, toUserID, mediaType, filePath, text, contextToken string) error
}

type Server struct {
	logger  *slog.Logger
	service Service
}

func NewServer(service Service, logger *slog.Logger) *Server {
	return &Server{logger: logger, service: service}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health/live", s.handleLive)
	mux.HandleFunc("GET /health/ready", s.handleReady)
	mux.HandleFunc("GET /api/version", s.handleVersion)
	mux.HandleFunc("POST /api/accounts/login/start", s.handleLoginStart)
	mux.HandleFunc("GET /api/accounts/login/status", s.handleLoginStatus)
	mux.HandleFunc("GET /api/accounts/login/qr", s.handleLoginQR)
	mux.HandleFunc("GET /api/accounts", s.handleAccounts)
	mux.HandleFunc("GET /api/events", s.handleEvents)
	mux.HandleFunc("GET /api/logs", s.handleLogs)
	mux.HandleFunc("GET /api/settings", s.handleGetSettings)
	mux.HandleFunc("POST /api/settings", s.handleUpdateSettings)
	mux.HandleFunc("POST /api/messages/send-text", s.handleSendText)
	mux.HandleFunc("POST /api/messages/send-media", s.handleSendMedia)
	return withJSONContentType(mux)
}

func (s *Server) handleLive(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"timestamp": time.Now().UTC(),
		"version":   coreversion.Current(),
	})
}

func (s *Server) handleReady(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"ok":        true,
		"timestamp": time.Now().UTC(),
		"version":   coreversion.Current(),
	})
}

func (s *Server) handleVersion(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, coreversion.Current())
}

func (s *Server) handleLoginStart(w http.ResponseWriter, r *http.Request) {
	var req struct {
		BaseURL string `json:"base_url"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil && !errors.Is(err, http.ErrBodyNotAllowed) && err.Error() != "EOF" {
		writeError(w, http.StatusBadRequest, err)
		return
	}
	session, err := s.service.StartLogin(r.Context(), strings.TrimSpace(req.BaseURL))
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (s *Server) handleLoginStatus(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
	if sessionID == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "session_id is required"})
		return
	}
	session, err := s.service.GetLoginStatus(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			writeJSON(w, http.StatusNotFound, map[string]any{"error": "login session not found"})
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, session)
}

func (s *Server) handleLoginQR(w http.ResponseWriter, r *http.Request) {
	sessionID := strings.TrimSpace(r.URL.Query().Get("session_id"))
	if sessionID == "" {
		http.Error(w, "session_id is required", http.StatusBadRequest)
		return
	}
	session, err := s.service.GetLoginSession(r.Context(), sessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "login session not found", http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if session.QRCodeURL == "" {
		http.Error(w, "qr code url is empty", http.StatusBadRequest)
		return
	}
	png, err := GenerateQRCodePNG(session.QRCodeURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "image/png")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(png)
}

func (s *Server) handleAccounts(w http.ResponseWriter, r *http.Request) {
	accounts, err := s.service.ListAccounts(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": accounts})
}

func (s *Server) handleEvents(w http.ResponseWriter, r *http.Request) {
	afterID, _ := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("after_id")), 10, 64)
	limit, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("limit")))
	events, err := s.service.ListEvents(r.Context(), afterID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": events})
}

func (s *Server) handleLogs(w http.ResponseWriter, r *http.Request) {
	afterID, _ := strconv.ParseInt(strings.TrimSpace(r.URL.Query().Get("after_id")), 10, 64)
	limit, _ := strconv.Atoi(strings.TrimSpace(r.URL.Query().Get("limit")))
	items, err := s.service.ListLogs(r.Context(), afterID, limit)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"items": items})
}

func (s *Server) handleGetSettings(w http.ResponseWriter, r *http.Request) {
	settings, err := s.service.GetSettings(r.Context())
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, settings)
}

func (s *Server) handleUpdateSettings(w http.ResponseWriter, r *http.Request) {
	var settings model.Settings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json body"})
		return
	}
	if strings.TrimSpace(settings.ListenAddr) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "listen_addr is required"})
		return
	}
	out, err := s.service.UpdateSettings(r.Context(), settings)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"settings":       out,
		"restart_needed": true,
	})
}

func (s *Server) handleSendText(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID    string `json:"account_id"`
		ToUserID     string `json:"to_user_id"`
		Text         string `json:"text"`
		ContextToken string `json:"context_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json body"})
		return
	}
	if strings.TrimSpace(req.AccountID) == "" || strings.TrimSpace(req.ToUserID) == "" || strings.TrimSpace(req.Text) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "account_id, to_user_id and text are required"})
		return
	}
	if err := s.service.SendText(r.Context(), req.AccountID, req.ToUserID, req.Text, req.ContextToken); err != nil {
		if isContextTokenMissingError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (s *Server) handleSendMedia(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountID    string `json:"account_id"`
		ToUserID     string `json:"to_user_id"`
		Type         string `json:"type"`
		FilePath     string `json:"file_path"`
		Text         string `json:"text"`
		ContextToken string `json:"context_token"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "invalid json body"})
		return
	}
	if strings.TrimSpace(req.AccountID) == "" || strings.TrimSpace(req.ToUserID) == "" || strings.TrimSpace(req.FilePath) == "" {
		writeJSON(w, http.StatusBadRequest, map[string]any{"error": "account_id, to_user_id and file_path are required"})
		return
	}
	if err := s.service.SendMedia(r.Context(), req.AccountID, req.ToUserID, req.Type, req.FilePath, req.Text, req.ContextToken); err != nil {
		if isContextTokenMissingError(err) {
			writeJSON(w, http.StatusBadRequest, map[string]any{"error": err.Error()})
			return
		}
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func isContextTokenMissingError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(strings.ToLower(err.Error()), "context token not found")
}

func withJSONContentType(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.WriteHeader(status)
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	_ = enc.Encode(payload)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{"error": err.Error()})
}

````

[⬆ 回到目录](#toc)

<a name="file-internal-ilink-client.go"></a>
## 📄 internal/ilink/client.go

````go
package ilink

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Client struct {
	httpClient     *http.Client
	channelVersion string
}

func NewClient(channelVersion string, timeout time.Duration) *Client {
	return &Client{
		httpClient:     &http.Client{Timeout: timeout},
		channelVersion: channelVersion,
	}
}

type QRCodeResponse struct {
	QRCode    string `json:"qrcode"`
	QRCodeURL string `json:"qrcode_img_content"`
}

type QRStatusResponse struct {
	Status      string `json:"status"`
	BotToken    string `json:"bot_token"`
	AccountID   string `json:"ilink_bot_id"`
	BaseURL     string `json:"baseurl"`
	ILinkUserID string `json:"ilink_user_id"`
}

type TextItem struct {
	Text string `json:"text,omitempty"`
}

type VoiceItem struct {
	Media          CDNMedia `json:"media,omitempty"`
	Text       string `json:"text,omitempty"`
	EncodeType int    `json:"encode_type,omitempty"`
	Playtime   int    `json:"playtime,omitempty"`
}

type FileItem struct {
	Media    CDNMedia `json:"media,omitempty"`
	FileName string `json:"file_name,omitempty"`
	MD5      string `json:"md5,omitempty"`
	Len      string `json:"len,omitempty"`
}

type CDNMedia struct {
	EncryptQueryParam string `json:"encrypt_query_param,omitempty"`
	AESKey            string `json:"aes_key,omitempty"`
	EncryptType       int    `json:"encrypt_type,omitempty"`
}

type ImageItem struct {
	Media      CDNMedia `json:"media,omitempty"`
	ThumbMedia CDNMedia `json:"thumb_media,omitempty"`
	AESKey     string   `json:"aeskey,omitempty"`
	MidSize    int      `json:"mid_size,omitempty"`
}

type VideoItem struct {
	Media      CDNMedia `json:"media,omitempty"`
	ThumbMedia CDNMedia `json:"thumb_media,omitempty"`
	VideoSize  int      `json:"video_size,omitempty"`
}

type MessageItem struct {
	Type      int        `json:"type,omitempty"`
	TextItem  *TextItem  `json:"text_item,omitempty"`
	VoiceItem *VoiceItem `json:"voice_item,omitempty"`
	FileItem  *FileItem  `json:"file_item,omitempty"`
	ImageItem *ImageItem `json:"image_item,omitempty"`
	VideoItem *VideoItem `json:"video_item,omitempty"`
}

type WeixinMessage struct {
	Seq          int64         `json:"seq,omitempty"`
	MessageID    int64         `json:"message_id,omitempty"`
	FromUserID   string        `json:"from_user_id,omitempty"`
	ToUserID     string        `json:"to_user_id,omitempty"`
	ClientID     string        `json:"client_id,omitempty"`
	CreateTimeMS int64         `json:"create_time_ms,omitempty"`
	MessageType  int           `json:"message_type,omitempty"`
	MessageState int           `json:"message_state,omitempty"`
	ItemList     []MessageItem `json:"item_list,omitempty"`
	ContextToken string        `json:"context_token,omitempty"`
}

type GetUpdatesResponse struct {
	Ret                int             `json:"ret,omitempty"`
	ErrCode            int             `json:"errcode,omitempty"`
	ErrMsg             string          `json:"errmsg,omitempty"`
	Messages           []WeixinMessage `json:"msgs,omitempty"`
	GetUpdatesBuf      string          `json:"get_updates_buf,omitempty"`
	LongPollingTimeout int             `json:"longpolling_timeout_ms,omitempty"`
}

type SendMessageResponse struct {
	Ret     int    `json:"ret,omitempty"`
	ErrCode int    `json:"errcode,omitempty"`
	ErrMsg  string `json:"errmsg,omitempty"`
}

const (
	UploadMediaTypeImage = 1
	UploadMediaTypeVideo = 2
	UploadMediaTypeFile  = 3
	UploadMediaTypeVoice = 4
)

func (c *Client) FetchQRCode(ctx context.Context, baseURL string) (QRCodeResponse, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/") + "/ilink/bot/get_bot_qrcode?bot_type=3")
	if err != nil {
		return QRCodeResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return QRCodeResponse{}, err
	}
	var out QRCodeResponse
	if err := c.doJSON(req, "", nil, &out); err != nil {
		return QRCodeResponse{}, err
	}
	return out, nil
}

func (c *Client) FetchQRCodeStatus(ctx context.Context, baseURL, qrCode string) (QRStatusResponse, error) {
	u, err := url.Parse(strings.TrimRight(baseURL, "/") + "/ilink/bot/get_qrcode_status?qrcode=" + url.QueryEscape(qrCode))
	if err != nil {
		return QRStatusResponse{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return QRStatusResponse{}, err
	}
	req.Header.Set("iLink-App-ClientVersion", "1")
	var out QRStatusResponse
	if err := c.doJSON(req, "", nil, &out); err != nil {
		return QRStatusResponse{}, err
	}
	return out, nil
}

func (c *Client) GetUpdates(ctx context.Context, baseURL, token, getUpdatesBuf string) (GetUpdatesResponse, error) {
	body := map[string]any{
		"get_updates_buf": getUpdatesBuf,
		"base_info": map[string]any{
			"channel_version": c.channelVersion,
		},
	}
	var out GetUpdatesResponse
	if err := c.postJSON(ctx, strings.TrimRight(baseURL, "/")+"/ilink/bot/getupdates", token, body, &out); err != nil {
		return GetUpdatesResponse{}, err
	}
	return out, nil
}

func (c *Client) SendTextMessage(ctx context.Context, baseURL, token, toUserID, text, contextToken string) error {
	msg := map[string]any{
		"from_user_id":  "",
		"to_user_id":    toUserID,
		"client_id":     fmt.Sprintf("wcfLink-%d", time.Now().UnixNano()),
		"message_type":  2,
		"message_state": 2,
		"item_list": []map[string]any{
			{
				"type": 1,
				"text_item": map[string]any{
					"text": text,
				},
			},
		},
	}
	if strings.TrimSpace(contextToken) != "" {
		msg["context_token"] = contextToken
	}

	body := map[string]any{
		"msg": msg,
		"base_info": map[string]any{
			"channel_version": c.channelVersion,
		},
	}
	var out SendMessageResponse
	if err := c.postJSON(ctx, strings.TrimRight(baseURL, "/")+"/ilink/bot/sendmessage", token, body, &out); err != nil {
		return err
	}
	if out.ErrCode != 0 || out.Ret != 0 {
		errText := out.ErrMsg
		if strings.TrimSpace(errText) == "" {
			errText = "sendmessage returned non-zero status"
		}
		return fmt.Errorf("%s (ret=%d errcode=%d)", errText, out.Ret, out.ErrCode)
	}
	return nil
}

func (c *Client) postJSON(ctx context.Context, endpoint, token string, body any, out any) error {
	payload, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	return c.doJSON(req, token, payload, out)
}

func (c *Client) doJSON(req *http.Request, token string, payload []byte, out any) error {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("AuthorizationType", "ilink_bot_token")
	req.Header.Set("X-WECHAT-UIN", randomWechatUIN())
	if len(payload) > 0 {
		req.Header.Set("Content-Length", fmt.Sprintf("%d", len(payload)))
	}
	if strings.TrimSpace(token) != "" {
		req.Header.Set("Authorization", "Bearer "+strings.TrimSpace(token))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("ilink http %d: %s", resp.StatusCode, string(raw))
	}
	if out == nil || len(raw) == 0 {
		return nil
	}
	if err := json.Unmarshal(raw, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func randomWechatUIN() string {
	n, err := rand.Int(rand.Reader, big.NewInt(1<<32-1))
	if err != nil {
		return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d", time.Now().UnixNano())))
	}
	return base64.StdEncoding.EncodeToString([]byte(n.String()))
}

````

[⬆ 回到目录](#toc)

<a name="file-internal-ilink-media.go"></a>
## 📄 internal/ilink/media.go

````go
package ilink

import (
	"bytes"
	"context"
	"crypto/aes"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type GetUploadURLResponse struct {
	UploadParam      string `json:"upload_param,omitempty"`
	ThumbUploadParam string `json:"thumb_upload_param,omitempty"`
}

type UploadedMedia struct {
	DownloadEncryptedQueryParam string
	AESKeyHex                   string
	PlainSize                   int
	CipherSize                  int
}

func (c *Client) GetUploadURL(
	ctx context.Context,
	baseURL,
	token string,
	reqBody map[string]any,
) (GetUploadURLResponse, error) {
	reqBody["base_info"] = map[string]any{
		"channel_version": c.channelVersion,
	}
	var out GetUploadURLResponse
	if err := c.postJSON(ctx, strings.TrimRight(baseURL, "/")+"/ilink/bot/getuploadurl", token, reqBody, &out); err != nil {
		return GetUploadURLResponse{}, err
	}
	return out, nil
}

func (c *Client) UploadLocalMedia(
	ctx context.Context,
	cdnBaseURL,
	baseURL,
	token,
	toUserID,
	filePath string,
	mediaType int,
) (UploadedMedia, error) {
	plaintext, err := os.ReadFile(filePath)
	if err != nil {
		return UploadedMedia{}, err
	}
	rawMD5 := md5.Sum(plaintext)
	aesKey := make([]byte, 16)
	if _, err := rand.Read(aesKey); err != nil {
		return UploadedMedia{}, err
	}
	fileKeyBytes := make([]byte, 16)
	if _, err := rand.Read(fileKeyBytes); err != nil {
		return UploadedMedia{}, err
	}
	fileKey := hex.EncodeToString(fileKeyBytes)
	ciphertext, err := encryptAesECB(plaintext, aesKey)
	if err != nil {
		return UploadedMedia{}, err
	}

	uploadResp, err := c.GetUploadURL(ctx, baseURL, token, map[string]any{
		"filekey":       fileKey,
		"media_type":    mediaType,
		"to_user_id":    toUserID,
		"rawsize":       len(plaintext),
		"rawfilemd5":    hex.EncodeToString(rawMD5[:]),
		"filesize":      len(ciphertext),
		"no_need_thumb": true,
		"aeskey":        hex.EncodeToString(aesKey),
	})
	if err != nil {
		return UploadedMedia{}, err
	}
	if strings.TrimSpace(uploadResp.UploadParam) == "" {
		return UploadedMedia{}, fmt.Errorf("getuploadurl returned empty upload_param")
	}

	downloadParam, err := c.uploadCiphertextToCDN(ctx, cdnBaseURL, uploadResp.UploadParam, fileKey, ciphertext)
	if err != nil {
		return UploadedMedia{}, err
	}
	return UploadedMedia{
		DownloadEncryptedQueryParam: downloadParam,
		AESKeyHex:                   hex.EncodeToString(aesKey),
		PlainSize:                   len(plaintext),
		CipherSize:                  len(ciphertext),
	}, nil
}

func (c *Client) SendImageMessage(ctx context.Context, baseURL, token, toUserID, contextToken, text string, uploaded UploadedMedia) error {
	item := map[string]any{
		"type": 2,
		"image_item": map[string]any{
			"media": map[string]any{
				"encrypt_query_param": uploaded.DownloadEncryptedQueryParam,
				"aes_key":             base64.StdEncoding.EncodeToString([]byte(uploaded.AESKeyHex)),
				"encrypt_type":        1,
			},
			"mid_size": uploaded.CipherSize,
		},
	}
	return c.sendMediaItems(ctx, baseURL, token, toUserID, contextToken, text, item)
}

func (c *Client) SendVideoMessage(ctx context.Context, baseURL, token, toUserID, contextToken, text string, uploaded UploadedMedia) error {
	item := map[string]any{
		"type": 5,
		"video_item": map[string]any{
			"media": map[string]any{
				"encrypt_query_param": uploaded.DownloadEncryptedQueryParam,
				"aes_key":             base64.StdEncoding.EncodeToString([]byte(uploaded.AESKeyHex)),
				"encrypt_type":        1,
			},
			"video_size": uploaded.CipherSize,
		},
	}
	return c.sendMediaItems(ctx, baseURL, token, toUserID, contextToken, text, item)
}

func (c *Client) SendFileMessage(ctx context.Context, baseURL, token, toUserID, contextToken, text, fileName string, uploaded UploadedMedia) error {
	item := map[string]any{
		"type": 4,
		"file_item": map[string]any{
			"media": map[string]any{
				"encrypt_query_param": uploaded.DownloadEncryptedQueryParam,
				"aes_key":             base64.StdEncoding.EncodeToString([]byte(uploaded.AESKeyHex)),
				"encrypt_type":        1,
			},
			"file_name": fileName,
			"len":       fmt.Sprintf("%d", uploaded.PlainSize),
		},
	}
	return c.sendMediaItems(ctx, baseURL, token, toUserID, contextToken, text, item)
}

func (c *Client) SendVoiceMessage(ctx context.Context, baseURL, token, toUserID, contextToken, text string, encodeType int, uploaded UploadedMedia) error {
	item := map[string]any{
		"type": 3,
		"voice_item": map[string]any{
			"media": map[string]any{
				"encrypt_query_param": uploaded.DownloadEncryptedQueryParam,
				"aes_key":             base64.StdEncoding.EncodeToString([]byte(uploaded.AESKeyHex)),
				"encrypt_type":        1,
			},
			"encode_type": encodeType,
			"text":        "",
		},
	}
	return c.sendMediaItems(ctx, baseURL, token, toUserID, contextToken, text, item)
}

func (c *Client) sendMediaItems(ctx context.Context, baseURL, token, toUserID, contextToken, text string, mediaItem map[string]any) error {
	items := make([]map[string]any, 0, 2)
	if strings.TrimSpace(text) != "" {
		items = append(items, map[string]any{
			"type": 1,
			"text_item": map[string]any{
				"text": text,
			},
		})
	}
	items = append(items, mediaItem)

	for _, item := range items {
		msg := map[string]any{
			"from_user_id":  "",
			"to_user_id":    toUserID,
			"client_id":     fmt.Sprintf("wcfLink-%d", time.Now().UnixNano()),
			"message_type":  2,
			"message_state": 2,
			"item_list":     []map[string]any{item},
		}
		if strings.TrimSpace(contextToken) != "" {
			msg["context_token"] = contextToken
		}
		var out SendMessageResponse
		if err := c.postJSON(ctx, strings.TrimRight(baseURL, "/")+"/ilink/bot/sendmessage", token, map[string]any{
			"msg": msg,
			"base_info": map[string]any{
				"channel_version": c.channelVersion,
			},
		}, &out); err != nil {
			return err
		}
		if out.ErrCode != 0 || out.Ret != 0 {
			errText := out.ErrMsg
			if strings.TrimSpace(errText) == "" {
				errText = "sendmessage returned non-zero status"
			}
			return fmt.Errorf("%s (ret=%d errcode=%d)", errText, out.Ret, out.ErrCode)
		}
	}
	return nil
}

func (c *Client) DownloadMessageMedia(ctx context.Context, cdnBaseURL string, item MessageItem) ([]byte, string, string, error) {
	switch item.Type {
	case 2:
		if item.ImageItem == nil || strings.TrimSpace(item.ImageItem.Media.EncryptQueryParam) == "" {
			return nil, "", "", fmt.Errorf("image media is missing")
		}
		aesKey := item.ImageItem.Media.AESKey
		if item.ImageItem.AESKey != "" {
			aesKey = base64.StdEncoding.EncodeToString([]byte(item.ImageItem.AESKey))
		}
		buf, err := c.downloadCDNMedia(ctx, cdnBaseURL, item.ImageItem.Media.EncryptQueryParam, aesKey)
		if err != nil {
			return nil, "", "", err
		}
		mime := detectMIME(buf, ".jpg")
		return buf, "image"+extensionFromMIME(mime, ".jpg"), mime, nil
	case 3:
		if item.VoiceItem == nil || strings.TrimSpace(item.VoiceItem.Media.EncryptQueryParam) == "" || strings.TrimSpace(item.VoiceItem.Media.AESKey) == "" {
			return nil, "", "", fmt.Errorf("voice media is missing")
		}
		buf, err := c.downloadCDNMedia(ctx, cdnBaseURL, item.VoiceItem.Media.EncryptQueryParam, item.VoiceItem.Media.AESKey)
		if err != nil {
			return nil, "", "", err
		}
		return buf, "voice.silk", "audio/silk", nil
	case 4:
		if item.FileItem == nil || strings.TrimSpace(item.FileItem.Media.EncryptQueryParam) == "" || strings.TrimSpace(item.FileItem.Media.AESKey) == "" {
			return nil, "", "", fmt.Errorf("file media is missing")
		}
		buf, err := c.downloadCDNMedia(ctx, cdnBaseURL, item.FileItem.Media.EncryptQueryParam, item.FileItem.Media.AESKey)
		if err != nil {
			return nil, "", "", err
		}
		fileName := item.FileItem.FileName
		if strings.TrimSpace(fileName) == "" {
			fileName = "file.bin"
		}
		return buf, fileName, detectMIME(buf, filepath.Ext(fileName)), nil
	case 5:
		if item.VideoItem == nil || strings.TrimSpace(item.VideoItem.Media.EncryptQueryParam) == "" || strings.TrimSpace(item.VideoItem.Media.AESKey) == "" {
			return nil, "", "", fmt.Errorf("video media is missing")
		}
		buf, err := c.downloadCDNMedia(ctx, cdnBaseURL, item.VideoItem.Media.EncryptQueryParam, item.VideoItem.Media.AESKey)
		if err != nil {
			return nil, "", "", err
		}
		return buf, "video.mp4", "video/mp4", nil
	default:
		return nil, "", "", fmt.Errorf("unsupported media item type %d", item.Type)
	}
}

func (c *Client) uploadCiphertextToCDN(ctx context.Context, cdnBaseURL, uploadParam, fileKey string, ciphertext []byte) (string, error) {
	endpoint := fmt.Sprintf("%s/upload?encrypted_query_param=%s&filekey=%s",
		strings.TrimRight(cdnBaseURL, "/"),
		url.QueryEscape(uploadParam),
		url.QueryEscape(fileKey),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(ciphertext))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/octet-stream")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("cdn upload http %d: %s", resp.StatusCode, string(body))
	}
	downloadParam := resp.Header.Get("x-encrypted-param")
	if strings.TrimSpace(downloadParam) == "" {
		return "", fmt.Errorf("cdn upload response missing x-encrypted-param")
	}
	return downloadParam, nil
}

func (c *Client) downloadCDNMedia(ctx context.Context, cdnBaseURL, encryptedQueryParam, aesKeyBase64 string) ([]byte, error) {
	endpoint := fmt.Sprintf("%s/download?encrypted_query_param=%s",
		strings.TrimRight(cdnBaseURL, "/"),
		url.QueryEscape(encryptedQueryParam),
	)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	raw, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, fmt.Errorf("cdn download http %d: %s", resp.StatusCode, string(raw))
	}
	if strings.TrimSpace(aesKeyBase64) == "" {
		return raw, nil
	}
	key, err := parseAESKey(aesKeyBase64)
	if err != nil {
		return nil, err
	}
	return decryptAesECB(raw, key)
}

func encryptAesECB(plaintext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	padded := pkcs7Pad(plaintext, block.BlockSize())
	out := make([]byte, len(padded))
	for start := 0; start < len(padded); start += block.BlockSize() {
		block.Encrypt(out[start:start+block.BlockSize()], padded[start:start+block.BlockSize()])
	}
	return out, nil
}

func decryptAesECB(ciphertext, key []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	if len(ciphertext)%block.BlockSize() != 0 {
		return nil, fmt.Errorf("ciphertext size %d is not a multiple of block size", len(ciphertext))
	}
	out := make([]byte, len(ciphertext))
	for start := 0; start < len(ciphertext); start += block.BlockSize() {
		block.Decrypt(out[start:start+block.BlockSize()], ciphertext[start:start+block.BlockSize()])
	}
	return pkcs7Unpad(out, block.BlockSize())
}

func parseAESKey(aesKeyBase64 string) ([]byte, error) {
	decoded, err := base64.StdEncoding.DecodeString(aesKeyBase64)
	if err != nil {
		return nil, err
	}
	if len(decoded) == 16 {
		return decoded, nil
	}
	if len(decoded) == 32 && isHexASCII(decoded) {
		return hex.DecodeString(string(decoded))
	}
	return nil, fmt.Errorf("unexpected aes_key length %d", len(decoded))
}

func isHexASCII(value []byte) bool {
	for _, b := range value {
		switch {
		case b >= '0' && b <= '9':
		case b >= 'a' && b <= 'f':
		case b >= 'A' && b <= 'F':
		default:
			return false
		}
	}
	return true
}

func pkcs7Pad(src []byte, blockSize int) []byte {
	padding := blockSize - len(src)%blockSize
	out := make([]byte, len(src)+padding)
	copy(out, src)
	for i := len(src); i < len(out); i++ {
		out[i] = byte(padding)
	}
	return out
}

func pkcs7Unpad(src []byte, blockSize int) ([]byte, error) {
	if len(src) == 0 || len(src)%blockSize != 0 {
		return nil, fmt.Errorf("invalid padded buffer size")
	}
	padding := int(src[len(src)-1])
	if padding == 0 || padding > blockSize || padding > len(src) {
		return nil, fmt.Errorf("invalid PKCS7 padding")
	}
	for _, b := range src[len(src)-padding:] {
		if int(b) != padding {
			return nil, fmt.Errorf("invalid PKCS7 padding")
		}
	}
	return src[:len(src)-padding], nil
}

func detectMIME(buf []byte, fallbackExt string) string {
	contentType := http.DetectContentType(buf)
	if contentType == "application/octet-stream" {
		switch strings.ToLower(fallbackExt) {
		case ".jpg", ".jpeg":
			return "image/jpeg"
		case ".png":
			return "image/png"
		case ".gif":
			return "image/gif"
		case ".mp4":
			return "video/mp4"
		case ".pdf":
			return "application/pdf"
		}
	}
	return contentType
}

func extensionFromMIME(mime, fallback string) string {
	switch mime {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/gif":
		return ".gif"
	case "video/mp4":
		return ".mp4"
	}
	return fallback
}

````

[⬆ 回到目录](#toc)

<a name="file-internal-model-types.go"></a>
## 📄 internal/model/types.go

````go
package model

import "time"

type LoginSession struct {
	SessionID   string     `json:"session_id"`
	BaseURL     string     `json:"base_url"`
	QRCode      string     `json:"-"`
	QRCodeURL   string     `json:"qr_code_url"`
	Status      string     `json:"status"`
	AccountID   string     `json:"account_id,omitempty"`
	ILinkUserID string     `json:"ilink_user_id,omitempty"`
	BotToken    string     `json:"-"`
	Error       string     `json:"error,omitempty"`
	StartedAt   time.Time  `json:"started_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty"`
}

type Account struct {
	AccountID     string     `json:"account_id"`
	BaseURL       string     `json:"base_url"`
	Token         string     `json:"-"`
	ILinkUserID   string     `json:"ilink_user_id,omitempty"`
	Enabled       bool       `json:"enabled"`
	LoginStatus   string     `json:"login_status"`
	LastError     string     `json:"last_error,omitempty"`
	GetUpdatesBuf string     `json:"-"`
	LastPollAt    *time.Time `json:"last_poll_at,omitempty"`
	LastInboundAt *time.Time `json:"last_inbound_at,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
}

type PeerContext struct {
	AccountID    string    `json:"account_id"`
	PeerUserID   string    `json:"peer_user_id"`
	ContextToken string    `json:"context_token"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type Event struct {
	ID           int64     `json:"id"`
	AccountID    string    `json:"account_id"`
	Direction    string    `json:"direction"`
	EventType    string    `json:"event_type"`
	FromUserID   string    `json:"from_user_id,omitempty"`
	ToUserID     string    `json:"to_user_id,omitempty"`
	MessageID    int64     `json:"message_id,omitempty"`
	ContextToken string    `json:"context_token,omitempty"`
	BodyText     string    `json:"body_text,omitempty"`
	MediaPath    string    `json:"media_path,omitempty"`
	MediaFileName string   `json:"media_file_name,omitempty"`
	MediaMimeType string   `json:"media_mime_type,omitempty"`
	RawJSON      string    `json:"raw_json"`
	CreatedAt    time.Time `json:"created_at"`
}

type LogEntry struct {
	ID        int64     `json:"id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Source    string    `json:"source"`
	MetaJSON  string    `json:"meta_json,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type Settings struct {
	ListenAddr string `json:"listen_addr"`
	WebhookURL string `json:"webhook_url"`
}

````

[⬆ 回到目录](#toc)

<a name="file-internal-store-store.go"></a>
## 📄 internal/store/store.go

````go
package store

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"github.com/lich0821/wcfLink/internal/ilink"
	"github.com/lich0821/wcfLink/internal/model"
)

type Store struct {
	db *sql.DB
}

func New(ctx context.Context, dbPath string) (*Store, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)

	s := &Store{db: db}
	if err := s.migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

func (s *Store) Ping(ctx context.Context) error {
	return s.db.PingContext(ctx)
}

func (s *Store) CreateLoginSession(ctx context.Context, session model.LoginSession) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO login_sessions (
  session_id, base_url, qr_code, qr_code_url, status, error, started_at, updated_at
) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		session.SessionID, session.BaseURL, session.QRCode, session.QRCodeURL, session.Status,
		session.Error, session.StartedAt.UTC(), session.UpdatedAt.UTC(),
	)
	return err
}

func (s *Store) GetLoginSession(ctx context.Context, sessionID string) (model.LoginSession, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT session_id, base_url, qr_code, qr_code_url, status, account_id, ilink_user_id, bot_token,
       error, started_at, updated_at, completed_at
FROM login_sessions
WHERE session_id = ?`, sessionID)
	var session model.LoginSession
	var completedAt sql.NullTime
	err := row.Scan(
		&session.SessionID, &session.BaseURL, &session.QRCode, &session.QRCodeURL, &session.Status,
		&session.AccountID, &session.ILinkUserID, &session.BotToken, &session.Error,
		&session.StartedAt, &session.UpdatedAt, &completedAt,
	)
	if err != nil {
		return model.LoginSession{}, err
	}
	if completedAt.Valid {
		session.CompletedAt = &completedAt.Time
	}
	return session, nil
}

func (s *Store) UpdateLoginSessionStatus(ctx context.Context, sessionID, status, errorText string) error {
	_, err := s.db.ExecContext(ctx, `
UPDATE login_sessions
SET status = ?, error = ?, updated_at = ?
WHERE session_id = ?`, status, errorText, time.Now().UTC(), sessionID)
	return err
}

func (s *Store) CompleteLoginSession(ctx context.Context, sessionID string, status ilink.QRStatusResponse) error {
	now := time.Now().UTC()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
UPDATE login_sessions
SET status = ?, account_id = ?, ilink_user_id = ?, bot_token = ?, base_url = ?, updated_at = ?, completed_at = ?
WHERE session_id = ?`,
		status.Status, status.AccountID, status.ILinkUserID, status.BotToken, status.BaseURL, now, now, sessionID,
	)
	if err != nil {
		return err
	}

	baseURL := status.BaseURL
	if baseURL == "" {
		var fallback string
		if err := tx.QueryRowContext(ctx, `SELECT base_url FROM login_sessions WHERE session_id = ?`, sessionID).Scan(&fallback); err == nil {
			baseURL = fallback
		}
	}

	_, err = tx.ExecContext(ctx, `
INSERT INTO accounts (
  account_id, base_url, token, ilink_user_id, enabled, login_status, created_at, updated_at
) VALUES (?, ?, ?, ?, 1, 'connected', ?, ?)
ON CONFLICT(account_id) DO UPDATE SET
  base_url = excluded.base_url,
  token = excluded.token,
  ilink_user_id = excluded.ilink_user_id,
  enabled = 1,
  login_status = 'connected',
  last_error = '',
  updated_at = excluded.updated_at`,
		status.AccountID, baseURL, status.BotToken, status.ILinkUserID, now, now,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) ListAccounts(ctx context.Context) ([]model.Account, error) {
	rows, err := s.db.QueryContext(ctx, `
SELECT account_id, base_url, token, ilink_user_id, enabled, login_status, last_error,
       get_updates_buf, last_poll_at, last_inbound_at, created_at, updated_at
FROM accounts
ORDER BY created_at ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Account
	for rows.Next() {
		var item model.Account
		var enabled int
		var lastPollAt sql.NullTime
		var lastInboundAt sql.NullTime
		if err := rows.Scan(
			&item.AccountID, &item.BaseURL, &item.Token, &item.ILinkUserID, &enabled, &item.LoginStatus,
			&item.LastError, &item.GetUpdatesBuf, &lastPollAt, &lastInboundAt, &item.CreatedAt, &item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		item.Enabled = enabled == 1
		if lastPollAt.Valid {
			item.LastPollAt = &lastPollAt.Time
		}
		if lastInboundAt.Valid {
			item.LastInboundAt = &lastInboundAt.Time
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) GetAccount(ctx context.Context, accountID string) (model.Account, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT account_id, base_url, token, ilink_user_id, enabled, login_status, last_error,
       get_updates_buf, last_poll_at, last_inbound_at, created_at, updated_at
FROM accounts
WHERE account_id = ?`, accountID)
	var item model.Account
	var enabled int
	var lastPollAt sql.NullTime
	var lastInboundAt sql.NullTime
	err := row.Scan(
		&item.AccountID, &item.BaseURL, &item.Token, &item.ILinkUserID, &enabled, &item.LoginStatus,
		&item.LastError, &item.GetUpdatesBuf, &lastPollAt, &lastInboundAt, &item.CreatedAt, &item.UpdatedAt,
	)
	if err != nil {
		return model.Account{}, err
	}
	item.Enabled = enabled == 1
	if lastPollAt.Valid {
		item.LastPollAt = &lastPollAt.Time
	}
	if lastInboundAt.Valid {
		item.LastInboundAt = &lastInboundAt.Time
	}
	return item, nil
}

func (s *Store) DeleteAccount(ctx context.Context, accountID string) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	statements := []string{
		`DELETE FROM accounts WHERE account_id = ?`,
		`DELETE FROM peer_contexts WHERE account_id = ?`,
		`DELETE FROM login_sessions WHERE account_id = ?`,
	}
	for _, stmt := range statements {
		if _, err := tx.ExecContext(ctx, stmt, accountID); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (s *Store) UpdateAccountPollState(ctx context.Context, accountID, getUpdatesBuf, loginStatus, lastError string) error {
	_, err := s.db.ExecContext(ctx, `
UPDATE accounts
SET get_updates_buf = ?, login_status = ?, last_error = ?, last_poll_at = ?, updated_at = ?
WHERE account_id = ?`, getUpdatesBuf, loginStatus, lastError, time.Now().UTC(), time.Now().UTC(), accountID)
	return err
}

func (s *Store) TouchAccountInbound(ctx context.Context, accountID string) error {
	now := time.Now().UTC()
	_, err := s.db.ExecContext(ctx, `
UPDATE accounts
SET last_inbound_at = ?, updated_at = ?
WHERE account_id = ?`, now, now, accountID)
	return err
}

func (s *Store) SaveInboundMessage(ctx context.Context, accountID string, msg ilink.WeixinMessage, mediaPath, mediaFileName, mediaMimeType string) error {
	raw, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	bodyText := extractBodyText(msg)
	now := time.Now().UTC()

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, `
INSERT OR IGNORE INTO events (
  account_id, direction, event_type, from_user_id, to_user_id, message_id, context_token, body_text, media_path, media_file_name, media_mime_type, raw_json, created_at
) VALUES (?, 'inbound', ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		accountID, detectEventType(msg), msg.FromUserID, msg.ToUserID, msg.MessageID, msg.ContextToken, bodyText, mediaPath, mediaFileName, mediaMimeType, string(raw), now,
	)
	if err != nil {
		return err
	}

	if stringsNotEmpty(msg.FromUserID, msg.ContextToken) {
		_, err = tx.ExecContext(ctx, `
INSERT INTO peer_contexts (account_id, peer_user_id, context_token, updated_at)
VALUES (?, ?, ?, ?)
ON CONFLICT(account_id, peer_user_id) DO UPDATE SET
  context_token = excluded.context_token,
  updated_at = excluded.updated_at`, accountID, msg.FromUserID, msg.ContextToken, now)
		if err != nil {
			return err
		}
	}

	_, err = tx.ExecContext(ctx, `
UPDATE accounts
SET last_inbound_at = ?, updated_at = ?, last_error = '', login_status = 'connected'
WHERE account_id = ?`, now, now, accountID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Store) GetPeerContext(ctx context.Context, accountID, peerUserID string) (model.PeerContext, error) {
	row := s.db.QueryRowContext(ctx, `
SELECT account_id, peer_user_id, context_token, updated_at
FROM peer_contexts
WHERE account_id = ? AND peer_user_id = ?`, accountID, peerUserID)
	var item model.PeerContext
	if err := row.Scan(&item.AccountID, &item.PeerUserID, &item.ContextToken, &item.UpdatedAt); err != nil {
		return model.PeerContext{}, err
	}
	return item, nil
}

func (s *Store) CreateOutboundEvent(ctx context.Context, accountID, eventType, toUserID, contextToken, bodyText, mediaPath, mediaFileName, mediaMimeType, rawJSON string) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO events (
  account_id, direction, event_type, to_user_id, context_token, body_text, media_path, media_file_name, media_mime_type, raw_json, created_at
) VALUES (?, 'outbound', ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		accountID, eventType, toUserID, contextToken, bodyText, mediaPath, mediaFileName, mediaMimeType, rawJSON, time.Now().UTC(),
	)
	return err
}

func (s *Store) AddLog(ctx context.Context, level, message, source, metaJSON string) error {
	_, err := s.db.ExecContext(ctx, `
INSERT INTO logs (level, message, source, meta_json, created_at)
VALUES (?, ?, ?, ?, ?)`, level, message, source, metaJSON, time.Now().UTC())
	return err
}

func (s *Store) ListLogs(ctx context.Context, afterID int64, limit int) ([]model.LogEntry, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT id, level, message, source, meta_json, created_at
FROM logs
WHERE id > ?
ORDER BY id ASC
LIMIT ?`, afterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.LogEntry
	for rows.Next() {
		var item model.LogEntry
		if err := rows.Scan(&item.ID, &item.Level, &item.Message, &item.Source, &item.MetaJSON, &item.CreatedAt); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) ListEvents(ctx context.Context, afterID int64, limit int) ([]model.Event, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	rows, err := s.db.QueryContext(ctx, `
SELECT id, account_id, direction, event_type, from_user_id, to_user_id, message_id, context_token, body_text, media_path, media_file_name, media_mime_type, raw_json, created_at
FROM events
WHERE id > ?
ORDER BY id ASC
LIMIT ?`, afterID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []model.Event
	for rows.Next() {
		var item model.Event
		if err := rows.Scan(
			&item.ID, &item.AccountID, &item.Direction, &item.EventType, &item.FromUserID, &item.ToUserID,
			&item.MessageID, &item.ContextToken, &item.BodyText, &item.MediaPath, &item.MediaFileName, &item.MediaMimeType, &item.RawJSON, &item.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (s *Store) migrate(ctx context.Context) error {
	stmts := []string{
		`PRAGMA journal_mode=WAL;`,
		`CREATE TABLE IF NOT EXISTS login_sessions (
			session_id TEXT PRIMARY KEY,
			base_url TEXT NOT NULL,
			qr_code TEXT NOT NULL,
			qr_code_url TEXT NOT NULL,
			status TEXT NOT NULL,
			account_id TEXT NOT NULL DEFAULT '',
			ilink_user_id TEXT NOT NULL DEFAULT '',
			bot_token TEXT NOT NULL DEFAULT '',
			error TEXT NOT NULL DEFAULT '',
			started_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			completed_at TIMESTAMP
		);`,
		`CREATE TABLE IF NOT EXISTS accounts (
			account_id TEXT PRIMARY KEY,
			base_url TEXT NOT NULL,
			token TEXT NOT NULL,
			ilink_user_id TEXT NOT NULL DEFAULT '',
			enabled INTEGER NOT NULL DEFAULT 1,
			login_status TEXT NOT NULL DEFAULT 'pending',
			last_error TEXT NOT NULL DEFAULT '',
			get_updates_buf TEXT NOT NULL DEFAULT '',
			last_poll_at TIMESTAMP,
			last_inbound_at TIMESTAMP,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS peer_contexts (
			account_id TEXT NOT NULL,
			peer_user_id TEXT NOT NULL,
			context_token TEXT NOT NULL,
			updated_at TIMESTAMP NOT NULL,
			PRIMARY KEY (account_id, peer_user_id)
		);`,
		`CREATE TABLE IF NOT EXISTS events (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			account_id TEXT NOT NULL,
			direction TEXT NOT NULL,
			event_type TEXT NOT NULL,
			from_user_id TEXT NOT NULL DEFAULT '',
			to_user_id TEXT NOT NULL DEFAULT '',
			message_id INTEGER NOT NULL DEFAULT 0,
			context_token TEXT NOT NULL DEFAULT '',
			body_text TEXT NOT NULL DEFAULT '',
			media_path TEXT NOT NULL DEFAULT '',
			media_file_name TEXT NOT NULL DEFAULT '',
			media_mime_type TEXT NOT NULL DEFAULT '',
			raw_json TEXT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);`,
		`CREATE UNIQUE INDEX IF NOT EXISTS idx_events_account_message_inbound
		 ON events(account_id, direction, message_id)
		 WHERE direction = 'inbound' AND message_id != 0;`,
		`CREATE TABLE IF NOT EXISTS logs (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			level TEXT NOT NULL,
			message TEXT NOT NULL,
			source TEXT NOT NULL,
			meta_json TEXT NOT NULL DEFAULT '',
			created_at TIMESTAMP NOT NULL
		);`,
	}

	for _, stmt := range stmts {
		if _, err := s.db.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("migrate: %w", err)
		}
	}
	for _, stmt := range []string{
		`ALTER TABLE events ADD COLUMN media_path TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE events ADD COLUMN media_file_name TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE events ADD COLUMN media_mime_type TEXT NOT NULL DEFAULT ''`,
	} {
		if err := s.execMigrationCompat(ctx, stmt); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) execMigrationCompat(ctx context.Context, stmt string) error {
	if _, err := s.db.ExecContext(ctx, stmt); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate column name") {
			return nil
		}
		return fmt.Errorf("migrate: %w", err)
	}
	return nil
}

func extractBodyText(msg ilink.WeixinMessage) string {
	for _, item := range msg.ItemList {
		switch item.Type {
		case 1:
			if item.TextItem != nil {
				return item.TextItem.Text
			}
		case 3:
			if item.VoiceItem != nil && item.VoiceItem.Text != "" {
				return item.VoiceItem.Text
			}
		case 2:
			return "[image]"
		case 4:
			if item.FileItem != nil && item.FileItem.FileName != "" {
				return "[file] " + item.FileItem.FileName
			}
			return "[file]"
		case 5:
			return "[video]"
		}
	}
	return ""
}

func detectEventType(msg ilink.WeixinMessage) string {
	for _, item := range msg.ItemList {
		switch item.Type {
		case 1:
			return "text"
		case 2:
			return "image"
		case 3:
			return "voice"
		case 4:
			return "file"
		case 5:
			return "video"
		}
	}
	return "unknown"
}

func stringsNotEmpty(values ...string) bool {
	for _, value := range values {
		if value == "" {
			return false
		}
	}
	return true
}

var ErrNotFound = errors.New("not found")

````

[⬆ 回到目录](#toc)

<a name="file-internal-worker-poller.go"></a>
## 📄 internal/worker/poller.go

````go
package worker

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/lich0821/wcfLink/internal/ilink"
	"github.com/lich0821/wcfLink/internal/model"
	"github.com/lich0821/wcfLink/internal/store"
)

type PollerManager struct {
	store  *store.Store
	client *ilink.Client
	logger *slog.Logger
	onMsg  func(context.Context, model.Account, ilink.WeixinMessage) error

	mu      sync.Mutex
	running map[string]context.CancelFunc
}

func NewPollerManager(
	st *store.Store,
	client *ilink.Client,
	logger *slog.Logger,
	onMsg func(context.Context, model.Account, ilink.WeixinMessage) error,
) *PollerManager {
	return &PollerManager{
		store:   st,
		client:  client,
		logger:  logger,
		onMsg:   onMsg,
		running: make(map[string]context.CancelFunc),
	}
}

func (m *PollerManager) StartEnabledAccounts(ctx context.Context) error {
	accounts, err := m.store.ListAccounts(ctx)
	if err != nil {
		return err
	}
	for _, account := range accounts {
		if account.Enabled && account.Token != "" {
			m.StartAccount(ctx, account)
		}
	}
	return nil
}

func (m *PollerManager) StartAccount(parent context.Context, account model.Account) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.running[account.AccountID]; exists {
		return
	}
	ctx, cancel := context.WithCancel(parent)
	m.running[account.AccountID] = cancel
	go m.run(ctx, account)
}

func (m *PollerManager) StopAll() {
	m.mu.Lock()
	defer m.mu.Unlock()
	for accountID, cancel := range m.running {
		cancel()
		delete(m.running, accountID)
	}
}

func (m *PollerManager) StopAccount(accountID string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	cancel, exists := m.running[accountID]
	if !exists {
		return
	}
	cancel()
	delete(m.running, accountID)
}

func (m *PollerManager) run(ctx context.Context, account model.Account) {
	log := m.logger.With("account_id", account.AccountID)
	log.Info("poller started")
	defer func() {
		m.mu.Lock()
		delete(m.running, account.AccountID)
		m.mu.Unlock()
		log.Info("poller stopped")
	}()

	current := account
	backoff := 2 * time.Second
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		resp, err := m.client.GetUpdates(ctx, current.BaseURL, current.Token, current.GetUpdatesBuf)
		if err != nil {
			log.Error("getupdates failed", "err", err)
			_ = m.store.UpdateAccountPollState(context.Background(), current.AccountID, current.GetUpdatesBuf, "error", err.Error())
			if !sleepWithContext(ctx, backoff) {
				return
			}
			if backoff < 30*time.Second {
				backoff *= 2
			}
			continue
		}
		backoff = 2 * time.Second

		if resp.GetUpdatesBuf != "" {
			current.GetUpdatesBuf = resp.GetUpdatesBuf
		}
		status := "connected"
		errText := ""
		if resp.ErrCode != 0 || resp.Ret != 0 {
			status = "warning"
			errText = resp.ErrMsg
			if errText == "" {
				errText = "getupdates returned non-zero status"
			}
		}
		_ = m.store.UpdateAccountPollState(context.Background(), current.AccountID, current.GetUpdatesBuf, status, errText)

		for _, msg := range resp.Messages {
			if msg.MessageType != 1 {
				continue
			}
			if m.onMsg == nil {
				if err := m.store.SaveInboundMessage(context.Background(), current.AccountID, msg, "", "", ""); err != nil {
					log.Error("save inbound event failed", "err", err, "message_id", msg.MessageID)
				}
				continue
			}
			if err := m.onMsg(context.Background(), current, msg); err != nil {
				log.Error("save inbound event failed", "err", err, "message_id", msg.MessageID)
			}
		}

		if refreshed, err := m.store.GetAccount(context.Background(), current.AccountID); err == nil {
			current = refreshed
		}

		delay := time.Second
		if resp.LongPollingTimeout > 0 {
			delay = 0
		}
		if delay > 0 && !sleepWithContext(ctx, delay) {
			return
		}
	}
}

func (m *PollerManager) LookupContextToken(ctx context.Context, accountID, peerUserID string) (string, error) {
	item, err := m.store.GetPeerContext(ctx, accountID, peerUserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return item.ContextToken, nil
}

func sleepWithContext(ctx context.Context, delay time.Duration) bool {
	if delay <= 0 {
		return true
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}

````

[⬆ 回到目录](#toc)

<a name="file-version-version.go"></a>
## 📄 version/version.go

````go
package version

import (
	"fmt"
	"runtime/debug"
	"strings"
)

type Info struct {
	Version   string `json:"version"`
	Commit    string `json:"commit,omitempty"`
	BuildTime string `json:"build_time,omitempty"`
	Modified  bool   `json:"modified,omitempty"`
}

func Current() Info {
	info := Info{
		Version: "devel",
	}

	buildInfo, ok := debug.ReadBuildInfo()
	if !ok {
		return info
	}

	if v := normalizeVersion(buildInfo.Main.Version); v != "" {
		info.Version = v
	}

	for _, setting := range buildInfo.Settings {
		switch setting.Key {
		case "vcs.revision":
			info.Commit = shortCommit(setting.Value)
		case "vcs.time":
			info.BuildTime = setting.Value
		case "vcs.modified":
			info.Modified = setting.Value == "true"
		}
	}

	return info
}

func String() string {
	info := Current()
	parts := []string{info.Version}

	meta := make([]string, 0, 3)
	if info.Commit != "" {
		meta = append(meta, info.Commit)
	}
	if info.BuildTime != "" {
		meta = append(meta, "built "+info.BuildTime)
	}
	if info.Modified {
		meta = append(meta, "dirty")
	}
	if len(meta) == 0 {
		return parts[0]
	}
	return fmt.Sprintf("%s (%s)", parts[0], strings.Join(meta, ", "))
}

func normalizeVersion(value string) string {
	value = strings.TrimSpace(value)
	if value == "" || value == "(devel)" {
		return "devel"
	}
	return value
}

func shortCommit(value string) string {
	value = strings.TrimSpace(value)
	if len(value) > 12 {
		return value[:12]
	}
	return value
}

````

[⬆ 回到目录](#toc)

---
### 📊 最终统计汇总
- **文件总数:** 17
- **代码总行数:** 3026
- **物理总大小:** 84.58 KB
