#!/usr/bin/env python3
"""
wcfLink 自动保活脚本

功能：定期检查所有活跃对话的 context_token，在即将过期前（如23小时）
      主动发送保活消息，让用户回复后刷新 token。

使用方法：
    python keepalive.py --wcflink-url http://127.0.0.1:17890

建议：通过 crontab 每小时运行一次
    0 * * * * python3 /opt/wcfLink/scripts/keepalive.py >> /var/log/wcflink_keepalive.log 2>&1
"""

import argparse
import requests
import sqlite3
import json
from datetime import datetime, timedelta
from pathlib import Path


class KeepAliveManager:
    def __init__(self, wcflink_url: str, db_path: str = None, expire_hours: int = 23):
        self.wcflink_url = wcflink_url.rstrip("/")
        self.expire_hours = expire_hours

        # 默认数据库路径
        if db_path is None:
            db_path = Path(__file__).parent.parent / "data" / "wcfLink.db"
        self.db_path = Path(db_path)

        # 保活消息模板
        self.keepalive_messages = [
            "【系统提示】您的会话将在1小时后过期，回复任意内容保持连接。",
            "您好，为确保服务正常，请回复此消息以保持会话活跃。",
            "【wcfLink】会话保活提醒，请回复确认。",
        ]

    def get_expiring_contexts(self) -> list:
        """获取即将过期的对话"""
        expire_threshold = datetime.utcnow() - timedelta(hours=self.expire_hours)

        conn = sqlite3.connect(self.db_path)
        cursor = conn.cursor()

        cursor.execute("""
            SELECT account_id, peer_user_id, context_token, updated_at
            FROM peer_contexts
            WHERE updated_at < ?
            ORDER BY updated_at ASC
        """, (expire_threshold.isoformat(),))

        results = cursor.fetchall()
        conn.close()

        return [
            {
                "account_id": row[0],
                "peer_user_id": row[1],
                "context_token": row[2],
                "updated_at": row[3]
            }
            for row in results
        ]

    def send_keepalive(self, account_id: str, peer_user_id: str, context_token: str) -> bool:
        """发送保活消息"""
        import random
        message = random.choice(self.keepalive_messages)

        try:
            resp = requests.post(
                f"{self.wcflink_url}/api/messages/send-text",
                json={
                    "account_id": account_id,
                    "to_user_id": peer_user_id,
                    "text": message,
                    "context_token": context_token
                },
                timeout=10
            )

            if resp.status_code == 200 and resp.json().get("ok"):
                print(f"[OK] {peer_user_id}: 保活消息已发送")
                return True
            else:
                print(f"[FAIL] {peer_user_id}: {resp.json()}")
                return False

        except Exception as e:
            print(f"[ERROR] {peer_user_id}: {e}")
            return False

    def run(self, dry_run: bool = False):
        """执行保活检查"""
        print(f"\n=== wcfLink 保活检查 {datetime.now().isoformat()} ===")
        print(f"数据库: {self.db_path}")
        print(f"过期阈值: {self.expire_hours} 小时")

        # 获取即将过期的对话
        contexts = self.get_expiring_contexts()
        print(f"发现 {len(contexts)} 个即将过期的对话")

        if not contexts:
            print("无需处理")
            return

        if dry_run:
            print("\n[DRY RUN] 以下对话将收到保活消息:")
            for ctx in contexts:
                print(f"  - {ctx['peer_user_id']} (上次更新: {ctx['updated_at']})")
            return

        # 发送保活消息
        success = 0
        failed = 0

        for ctx in contexts:
            if self.send_keepalive(
                ctx["account_id"],
                ctx["peer_user_id"],
                ctx["context_token"]
            ):
                success += 1
            else:
                failed += 1

        print(f"\n完成: 成功 {success}, 失败 {failed}")


def main():
    parser = argparse.ArgumentParser(description="wcfLink 自动保活")
    parser.add_argument(
        "--wcflink-url",
        default="http://127.0.0.1:17890",
        help="wcfLink API 地址"
    )
    parser.add_argument(
        "--db-path",
        default=None,
        help="SQLite 数据库路径"
    )
    parser.add_argument(
        "--expire-hours",
        type=int,
        default=23,
        help="多少小时未互动视为即将过期 (默认 23)"
    )
    parser.add_argument(
        "--dry-run",
        action="store_true",
        help="仅显示将处理的对话，不实际发送"
    )

    args = parser.parse_args()

    manager = KeepAliveManager(
        wcflink_url=args.wcflink_url,
        db_path=args.db_path,
        expire_hours=args.expire_hours
    )
    manager.run(dry_run=args.dry_run)


if __name__ == "__main__":
    main()
