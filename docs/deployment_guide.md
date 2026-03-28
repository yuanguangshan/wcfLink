# wcfLink 服务器部署指南

## 一、部署架构

```
┌─────────────────┐     HTTP API      ┌─────────────────┐
│  Flask 应用     │ ◀───────────────▶ │   wcfLink       │
│  (你的业务)     │   127.0.0.1:17890 │   (微信通道)    │
└─────────────────┘                   └─────────────────┘
                                              │
                                              ▼
                                      ┌─────────────────┐
                                      │   微信 iLink    │
                                      │   (官方服务)    │
                                      └─────────────────┘
```

## 二、服务器部署步骤

### 1. 编译 Linux 版本

```bash
# 在 Mac 上交叉编译 Linux amd64 版本
GOOS=linux GOARCH=amd64 go build -o ./bin/wcfLink-linux ./cmd/wcfLink

# 或 arm64 (如 AWS Graviton)
GOOS=linux GOARCH=arm64 go build -o ./bin/wcfLink-linux-arm64 ./cmd/wcfLink
```

### 2. 上传到服务器

```bash
# 创建目录
ssh user@your-server "mkdir -p /opt/wcfLink"

# 上传文件
scp ./bin/wcfLink-linux user@your-server:/opt/wcfLink/wcfLink
scp -r ./bin/data user@your-server:/opt/wcfLink/
```

### 3. 创建 systemd 服务

```bash
# /etc/systemd/system/wcflink.service
[Unit]
Description=wcfLink WeChat Bot Service
After=network.target

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/wcfLink
ExecStart=/opt/wcfLink/wcfLink
Restart=always
RestartSec=5
Environment=WCFLINK_LISTEN_ADDR=127.0.0.1:17890
Environment=WCFLINK_STATE_DIR=/opt/wcfLink/data
Environment=WCFLINK_LOG_LEVEL=INFO

[Install]
WantedBy=multi-user.target
```

```bash
# 启用并启动服务
sudo systemctl daemon-reload
sudo systemctl enable wcflink
sudo systemctl start wcflink

# 查看状态
sudo systemctl status wcflink
```

### 4. 配置 Nginx 反向代理 (可选，如需外网访问)

```nginx
# /etc/nginx/sites-available/wcflink
server {
    listen 80;
    server_name your-domain.com;

    # 建议添加认证
    auth_basic "wcfLink API";
    auth_basic_user_file /etc/nginx/.htpasswd;

    location / {
        proxy_pass http://127.0.0.1:17890;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_read_timeout 60s;
    }
}
```

## 三、Flask 集成示例

### 方法一：直接 HTTP 调用

```python
# wcflink_client.py
import requests

class WcfLinkClient:
    def __init__(self, base_url="http://127.0.0.1:17890"):
        self.base_url = base_url

    def get_accounts(self):
        """获取已登录账号列表"""
        resp = requests.get(f"{self.base_url}/api/accounts")
        return resp.json()

    def send_text(self, account_id, to_user_id, text):
        """发送文本消息"""
        resp = requests.post(
            f"{self.base_url}/api/messages/send-text",
            json={
                "account_id": account_id,
                "to_user_id": to_user_id,
                "text": text
            }
        )
        return resp.json()

    def send_image(self, account_id, to_user_id, file_path, text=""):
        """发送图片"""
        resp = requests.post(
            f"{self.base_url}/api/messages/send-media",
            json={
                "account_id": account_id,
                "to_user_id": to_user_id,
                "type": "image",
                "file_path": file_path,
                "text": text
            }
        )
        return resp.json()

    def send_file(self, account_id, to_user_id, file_path, text=""):
        """发送文件"""
        resp = requests.post(
            f"{self.base_url}/api/messages/send-media",
            json={
                "account_id": account_id,
                "to_user_id": to_user_id,
                "type": "file",
                "file_path": file_path,
                "text": text
            }
        )
        return resp.json()

    def get_events(self, after_id=0, limit=100):
        """获取事件列表"""
        resp = requests.get(
            f"{self.base_url}/api/events",
            params={"after_id": after_id, "limit": limit}
        )
        return resp.json()

    def set_webhook(self, webhook_url):
        """设置 Webhook URL (接收消息回调)"""
        resp = requests.post(
            f"{self.base_url}/api/settings",
            json={
                "listen_addr": "127.0.0.1:17890",
                "webhook_url": webhook_url
            }
        )
        return resp.json()


# 使用示例
if __name__ == "__main__":
    client = WcfLinkClient()

    # 获取账号
    accounts = client.get_accounts()
    print(f"已登录账号: {accounts}")

    if accounts.get("items"):
        account_id = accounts["items"][0]["account_id"]

        # 发送消息
        result = client.send_text(
            account_id=account_id,
            to_user_id="xxx@im.wechat",  # 对方微信 ID
            text="Hello from Flask!"
        )
        print(f"发送结果: {result}")
```

### 方法二：Flask 蓝图集成

```python
# app.py
from flask import Flask, request, jsonify
import requests

app = Flask(__name__)

WCFLINK_URL = "http://127.0.0.1:17890"

@app.route("/api/wechat/send", methods=["POST"])
def send_wechat():
    """发送微信消息的 API 端点"""
    data = request.json

    # 参数校验
    to_user = data.get("to_user")
    message = data.get("message")
    if not to_user or not message:
        return jsonify({"error": "缺少 to_user 或 message"}), 400

    # 获取第一个可用账号
    accounts = requests.get(f"{WCFLINK_URL}/api/accounts").json()
    if not accounts.get("items"):
        return jsonify({"error": "没有已登录的微信账号"}), 500

    account_id = accounts["items"][0]["account_id"]

    # 发送消息
    resp = requests.post(
        f"{WCFLINK_URL}/api/messages/send-text",
        json={
            "account_id": account_id,
            "to_user_id": to_user,
            "text": message
        }
    )

    return jsonify(resp.json())


@app.route("/api/wechat/webhook", methods=["POST"])
def wechat_webhook():
    """接收 wcfLink 推送的微信消息"""
    data = request.json
    print(f"收到微信消息: {data}")

    # 处理消息逻辑
    from_user = data.get("from_user_id")
    text = data.get("body_text")

    # TODO: 你的业务逻辑
    # 例如：自动回复、转发到其他系统等

    return jsonify({"ok": True})


@app.route("/api/wechat/events", methods=["GET"])
def get_events():
    """获取微信事件列表"""
    after_id = request.args.get("after_id", 0, type=int)
    limit = request.args.get("limit", 100, type=int)

    resp = requests.get(
        f"{WCFLINK_URL}/api/events",
        params={"after_id": after_id, "limit": limit}
    )
    return jsonify(resp.json())


if __name__ == "__main__":
    app.run(host="0.0.0.0", port=5000)
```

### 方法三：使用 Webhook 接收消息

```python
# webhook_server.py
from flask import Flask, request, jsonify
import requests

app = Flask(__name__)

WCFLINK_URL = "http://127.0.0.1:17890"

@app.route("/wechat/callback", methods=["POST"])
def wechat_callback():
    """wcfLink Webhook 回调地址"""
    data = request.json

    event_type = data.get("event_type")
    from_user = data.get("from_user_id")
    text = data.get("body_text")
    account_id = data.get("account_id")

    print(f"[{event_type}] {from_user}: {text}")

    # 示例：自动回复
    if text and "hello" in text.lower():
        requests.post(
            f"{WCFLINK_URL}/api/messages/send-text",
            json={
                "account_id": account_id,
                "to_user_id": from_user,
                "text": "你好！我是机器人，收到你的消息了。"
            }
        )

    return jsonify({"ok": True})


def setup_webhook():
    """设置 wcfLink 的 Webhook URL"""
    # 假设你的 Flask 服务运行在 5000 端口
    webhook_url = "http://your-server:5000/wechat/callback"

    resp = requests.post(
        f"{WCFLINK_URL}/api/settings",
        json={
            "listen_addr": "127.0.0.1:17890",
            "webhook_url": webhook_url
        }
    )
    print(f"Webhook 设置结果: {resp.json()}")


if __name__ == "__main__":
    # 先设置 webhook
    setup_webhook()

    # 启动 Flask 服务
    app.run(host="0.0.0.0", port=5000)
```

## 四、完整部署脚本

```bash
#!/bin/bash
# deploy.sh - 一键部署脚本

SERVER="user@your-server"
REMOTE_DIR="/opt/wcfLink"

echo "=== 编译 Linux 版本 ==="
GOOS=linux GOARCH=amd64 go build -o ./bin/wcfLink-linux ./cmd/wcfLink

echo "=== 创建远程目录 ==="
ssh $SERVER "sudo mkdir -p $REMOTE_DIR && sudo chown $USER:$USER $REMOTE_DIR"

echo "=== 上传文件 ==="
scp ./bin/wcfLink-linux $SERVER:$REMOTE_DIR/wcfLink
ssh $SERVER "chmod +x $REMOTE_DIR/wcfLink"

echo "=== 创建 systemd 服务 ==="
ssh $SERVER "sudo tee /etc/systemd/system/wcflink.service" << 'EOF'
[Unit]
Description=wcfLink WeChat Bot Service
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/wcfLink
ExecStart=/opt/wcfLink/wcfLink
Restart=always
RestartSec=5
Environment=WCFLINK_LISTEN_ADDR=127.0.0.1:17890
Environment=WCFLINK_STATE_DIR=/opt/wcfLink/data
Environment=WCFLINK_LOG_LEVEL=INFO

[Install]
WantedBy=multi-user.target
EOF

echo "=== 启动服务 ==="
ssh $SERVER "sudo systemctl daemon-reload && sudo systemctl enable wcflink && sudo systemctl start wcflink"

echo "=== 检查状态 ==="
ssh $SERVER "sudo systemctl status wcflink --no-pager"

echo "=== 部署完成 ==="
echo "API 地址: http://your-server:17890"
echo "请记得扫码登录微信账号！"
```

## 五、首次登录流程

部署后需要扫码登录微信：

```bash
# 1. 发起登录
curl -X POST http://127.0.0.1:17890/api/accounts/login/start

# 2. 下载二维码
curl -o qrcode.png "http://127.0.0.1:17890/api/accounts/login/qr?session_id=login_xxx"

# 3. 用微信扫码

# 4. 检查登录状态
curl "http://127.0.0.1:17890/api/accounts/login/status?session_id=login_xxx"
```

## 六、自动保活方案（解决 24 小时断联问题)

**问题**： context_token 约 24 小时过期，过期后无法主动发消息给用户。

**解决方案**： 使用保活脚本，定期检查即将过期的对话，发送提醒消息。
### 1. 保活脚本

已在 [scripts/keepalive.py](scripts/keepalive.py) 创建了自动保活脚本。
### 2. 配置 crontab
```bash
# 编辑 crontab
crontab -e

# 添加以下行（每小时检查一次）
0 * * * * python3 /opt/wcfLink/scripts/keepalive.py --wcflink-url http://127.0.0.1:17890 >> /var/log/wcflink_keepalive.log 2>&1
```
### 3. 保活消息示例
用户会收到类似这样的消息：
> "您的会话即将过期，请回复任意内容保持连接。"

## 七、注意事项
1. **安全性**：
   - wcfLink 默认无认证，建议只监听 127.0.0.1
   - 如需外网访问，务必添加认证
2. **数据备份**：
   - 定期备份 `data/wcfLink.db` 数据库
3. **监控**：
   ```bash
   # 健康检查
   curl http://127.0.0.1:17890/health/live
   ```
