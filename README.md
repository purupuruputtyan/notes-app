# notes-app

ユーザー管理 API と React フロントを同じリポジトリ（モノレポ）で開発するためのプロジェクトです。

## リポジトリ構成

| ディレクトリ | 役割 |
|--------------|------|
| `backend/` | Go 製 API（`cmd/api` がエントリ） |
| `frontend/` | React + TypeScript + Vite |
| `db/migrations/` | PostgreSQL 用マイグレーション（goose） |

フロントとバックは **別ディレクトリではあるが同一 Git リポジトリ** です。ルートの `compose.yaml` で両方と DB をまとめて起動できます。

## アーキテクチャ概要

### バックエンド（`backend/`）

**ハンドラ → ユースケース → リポジトリ** のレイヤー構成です。

```text
HTTP (cmd/api, handler/user)
    → usecase/user（バリデーション・パスワードハッシュ）
    → repository/user（sqlboiler で生成した models を利用）
    → PostgreSQL
```

- **`internal/domain/user`** … リポジトリのインターフェースとドメインエラー（`AppError`）
- **`internal/usecase/user`** … 入力検証・bcrypt・トランザクションに相当するユースケース
- **`internal/repository/user`** … DB 実装（sqlboiler）
- **`internal/models`** … sqlboiler 生成コード（手編集しない）

### フロントエンド（`frontend/`）

- **Vite** で開発サーバー・ビルド
- **React 19** 系 + TypeScript
- **Storybook**（`src/stories/`）と Vitest 連携あり

API のベース URL は **`VITE_API_BASE_URL`**（Docker Compose では `http://localhost:8080`）を想定しています。ブラウザからバックエンドへアクセスするため、開発時は CORS やプロキシの要否に注意してください。

## 前提条件

- Docker / Docker Compose
- ローカルで Go や Node を直接使う場合は、各ディレクトリの `README` や `go.mod` / `package.json` を参照

## 環境構築（Docker Compose 推奨）

リポジトリルートで次を実行します。

```bash
cd backend
make setup
```

`setup` はコンテナ起動（`-d`）→ **マイグレーション適用** → **sqlboiler でモデル生成** まで行います。

- API: `http://localhost:8080`
- ヘルス: `GET http://localhost:8080/healthz`
- フロント（Compose で起動した場合）: `http://localhost:5173`

手動で段階実行する(`make set up`を使わない)場合:

```bash
docker compose up -d --build
cd backend && make migrate-up && make generate-models
```

ログを流し読みする場合:

```bash
cd backend && make logs
```

停止:

```bash
cd backend && make down
```

## バックエンドの環境変数

`compose.yaml` の `backend` サービスと同じキーをローカル実行時にも設定します。

| 変数 | 説明 |
|------|------|
| `PORT` | HTTP 待受（例: `8080` または `:8080`。未設定時は `:8080`） |
| `DB_HOST` | PostgreSQL ホスト |
| `DB_PORT` | ポート番号 |
| `DB_USER` / `DB_PASSWORD` / `DB_NAME` | 接続資格情報 |
| `DB_SSLMODE` | 例: `disable`（開発） / 本番は適切な TLS モード |
| `DATABASE_URL` | **goose マイグレーション用**（`make migrate-up` 等） |

必須チェックは `cmd/api` 側で `DB_HOST` など（パスワード除く）に対して行っています。

## Makefile（`backend/Makefile`）の主なターゲット

| ターゲット | 内容 |
|------------|------|
| `make up` / `make down` | Compose 起動・停止 |
| `make migrate-create name=...` | 新規マイグレーション SQL の雛形作成（`name` 必須） |
| `make migrate-up` / `migrate-down` / `migrate-status` | goose の適用・戻し・状態確認 |
| `make generate-models` | sqlboiler で `internal/models` を再生成 |
| `make fmt` / `fmt-check` / `vet` / `test` | Go の整形・静的解析・テスト |

## sqlboiler

- 設定: `backend/sqlboiler.toml`
- **Docker 内の DB が起動した状態**で `make generate-models`（＝コンテナ内で `sqlboiler psql`）を実行します。
- 接続先は開発用です。**本番の認証情報を `sqlboiler.toml` にコミットしないでください。**

## セキュリティについて（開発用の限界）

- `compose.yaml` の DB パスワード等は **ローカル開発用の固定値**です。本番や共有クラスタでは使わず、シークレット管理に置き換えてください。
- `compose.yaml` 先頭と `sqlboiler.toml` 先頭に注意コメントを入れています。

## 関連ドキュメント

- フロントのテンプレート説明: `frontend/README.md`（Vite テンプレ由来の ESLint 説明など）
