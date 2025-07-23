# 技術設計書: Todoアプリ

## 1. アーキテクチャ

- **フロントエンド**: React (TypeScript) を使用したシングルページアプリケーション (SPA)。Vite をビルドツールとして利用する。
- **バックエンド**: Go (Ginフレームワーク) を使用した RESTful API。ユーザー認証とデータ永続化を担う。
- **データベース**: PostgreSQL を使用し、ユーザーとタスクの情報を格納する。
- **認証**: JSON Web Tokens (JWT) を用いたステートレス認証を実装する。
- **監視**: OpenTelemetry を導入し、トレース、メトリクス、ログを収集する。

## 2. フロントエンド設計

### 2.1. コンポーネント構成

- `App.tsx`: アプリケーションのルートコンポーネント。ルーティングを管理する。
- **Pages**:
    - `SignupPage.tsx`: 新規アカウント登録ページ。
    - `LoginPage.tsx`: ログインページ。
    - `TodoPage.tsx`: メインのTODOリストページ。
- **Components**:
    - `Navbar.tsx`: ナビゲーションバー。ログイン状態に応じて表示を切り替える。
    - `TodoList.tsx`: タスクのリストを表示するコンポーネント。
    - `TodoItem.tsx`: 個々のタスクを表示・操作するコンポーネント。
    - `AddTodoForm.tsx`: 新しいタスクを追加するフォーム。
    - `FilterButtons.tsx`: タスクをフィルタリングするためのボタン群。
    - `AuthForm.tsx`: サインアップとログインで共通利用するフォームコンポーネント。

### 2.2. 状態管理

- React Context API または Zustand を利用して、認証状態（ユーザートークン）とタスクデータを管理する。

### 2.3. ルーティング

- React Router を使用して、以下のルートを定義する。
    - `/signup`: 新規登録ページ
    - `/login`: ログインページ
    - `/`: TODOリストページ（要認証）

## 3. バックエンド設計

### 3.1. APIエンドポイント

- **Auth**
    - `POST /api/auth/signup`: ユーザー登録
    - `POST /api/auth/login`: ログイン
- **Tasks** (要認証)
    - `GET /api/tasks`: ログインユーザーの全タスクを取得
    - `POST /api/tasks`: 新規タスクを作成
    - `PUT /api/tasks/:id`: タスクを更新（完了状態の切り替え）
    - `DELETE /api/tasks/:id`: タスクを削除

### 3.2. データモデル (PostgreSQL)

- **users テーブル**
    - `id`: SERIAL PRIMARY KEY
    - `username`: VARCHAR(255) UNIQUE NOT NULL
    - `password_hash`: VARCHAR(255) NOT NULL
    - `created_at`: TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
- **tasks テーブル**
    - `id`: SERIAL PRIMARY KEY
    - `user_id`: INTEGER REFERENCES users(id) ON DELETE CASCADE
    - `title`: TEXT NOT NULL
    - `completed`: BOOLEAN DEFAULT FALSE
    - `created_at`: TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP

## 4. 認証フロー

1.  ユーザーがサインアップまたはログインする。
2.  成功すると、バックエンドはJWTを生成してフロントエンドに返す。
3.  フロントエンドはJWTをローカルストレージなどに保存する。
4.  以降のAPIリクエストでは、AuthorizationヘッダーにJWTを付与して送信する。
5.  バックエンドはJWTを検証し、リクエストされた操作を許可する。

## 5. 監視設計

- Goバックエンドアプリケーションに OpenTelemetry SDK を統合する。
- **構造化ロギング**（例: `slog`）を導入し、リクエストのコンテキストに紐づく **`trace_id`** を全てのログ出力に含める。これにより、特定のトレースに関連するログを簡単に検索・分析できるようにする。
- HTTPリクエストごとにトレース情報を生成し、各処理のパフォーマンスを可視化する。データベースへのクエリや外部API呼び出しなども計測対象とし、スパンとしてトレースに含める。
- 収集したトレース、メトリクス、ログは、Jaeger（トレース）、Prometheus（メトリクス）、Loki（ログ）などの互換性のあるバックエンドシステムに送信し、統合的な分析やアラートに活用する。