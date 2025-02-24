# Go REST API

このリポジトリは、GoとEchoフレームワークを使用したRESTful APIのサンプルアプリケーションです。

## 機能

- ユーザー認証（サインアップ/ログイン/ログアウト）
- JWTを使用したセッション管理
- タスクの作成/読み取り/更新/削除（CRUD）操作
- CSRF対策
- CORSサポート

## 技術スタック

- Go 1.21.13
- Echo (Webフレームワーク)
- GORM (ORM)
- PostgreSQL (データベース)
- JWT (認証)
- Docker (開発環境)

## セットアップ

1. リポジトリをクローン
```bash
git clone <repository-url>
cd go-rest-api
```

2. 環境変数の設定
`.env`ファイルを作成し、以下の環境変数を設定してください：
```
POSTGRES_USER=udemy
POSTGRES_PW=udemy
POSTGRES_DB=udemy
POSTGRES_PORT=5434
POSTGRES_HOST=localhost
GO_ENV=dev
API_DOMAIN=localhost
FE_URL=http://localhost:3000
SECRET=your-secret-key
```

3. データベースの起動
```bash
docker-compose up -d
```

4. マイグレーションの実行
```bash
go run migrate/migrate.go
```

5. アプリケーションの起動
```bash
 GO_ENV=dev go run main.go
```

## API エンドポイント

### 認証関連

- `POST /signup` - ユーザー登録
  - リクエストボディ: `{ "email": "user@example.com", "password": "password" }`

- `POST /login` - ログイン
  - リクエストボディ: `{ "email": "user@example.com", "password": "password" }`
  - 成功時: JWTトークンがCookieにセット

- `POST /logout` - ログアウト
  - Cookieからトークンを削除

- `GET /csrf` - CSRFトークンの取得
  - CSRFトークンを返却

### タスク関連 (要認証)

- `GET /tasks` - 全タスクの取得
- `GET /tasks/:taskId` - 特定のタスクの取得
- `POST /tasks` - タスクの作成
  - リクエストボディ: `{ "title": "タスクのタイトル" }`
- `PUT /tasks/:taskId` - タスクの更新
  - リクエストボディ: `{ "title": "新しいタイトル" }`
- `DELETE /tasks/:taskId` - タスクの削除

## 認証フロー

1. ユーザー登録（/signup）
   - メールアドレスとパスワードでアカウントを作成
   - パスワードはbcryptでハッシュ化して保存

2. ログイン（/login）
   - メールアドレスとパスワードで認証
   - 認証成功時、JWTトークンを生成
   - トークンはHttpOnlyクッキーとして保存（XSS対策）

3. 保護されたエンドポイントへのアクセス
   - リクエスト時にクッキーから自動的にJWTトークンを送信
   - ミドルウェアでトークンを検証
   - トークンが有効な場合のみアクセスを許可

4. CSRF対策
   - CSRFトークンを/csrfエンドポイントで取得
   - 状態を変更するリクエスト（POST/PUT/DELETE）時にトークンを送信
   - ミドルウェアでトークンを検証

## セキュリティ対策

1. パスワードハッシュ化
   - bcryptを使用してパスワードをハッシュ化
   - 生のパスワードは保存しない

2. JWTトークン
   - HttpOnlyクッキーで保存（JavaScriptからのアクセス防止）
   - 有効期限を12時間に設定

3. CORS設定
   - 許可されたオリジンからのリクエストのみ受け付け
   - クレデンシャル（Cookie）を含むリクエストを許可

4. CSRF対策
   - トークンベースの保護
   - 安全なCookie設定

## データベース設計

### Usersテーブル
- id (PRIMARY KEY)
- email (UNIQUE)
- password (ハッシュ化)
- created_at
- updated_at

### Tasksテーブル
- id (PRIMARY KEY)
- title
- user_id (FOREIGN KEY)
- created_at
- updated_at

## エラーハンドリング

- バリデーションエラー: 400 Bad Request
- 認証エラー: 401 Unauthorized
- 権限エラー: 403 Forbidden
- リソース未検出: 404 Not Found
- サーバーエラー: 500 Internal Server Error
