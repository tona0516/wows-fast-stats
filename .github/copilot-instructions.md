# copilot-instructions.md

このドキュメントは本リポジトリで GitHub Copilot / AI アシスタントを利用する際のガイドラインをまとめたものです。日本語で明示的に記述し、過度な自動生成による品質低下やライセンスリスクを避けます。

---
## 1. プロジェクト概要
- リポジトリ名: `wows-fast-stats`
- 主目的: World of Warships の対戦中のマッチの各種統計・効率指標を高速に取得/表示するデスクトップアプリ。
- 主構成: Go (`internal/`), Svelte + TypeScript フロントエンド (`frontend/`), 設定ファイル・ユーザーデータ (`user_data/`).

## 2. 技術スタック
- Go 1.x (`go.mod` 参照)
- Wails (デスクトップ統合)
- Svelte + TypeScript + Vite (フロントエンド)
- テスト: Go 標準 testing (`*_test.go`), Jest (フロントエンド)

## 3. コーディング一般方針
- 可読性 > 省行数。過度なチェーンやマジックナンバー禁止。
- エラー処理: `internal/apperr` のパターン・`error` をラップする既存方式踏襲。
- ログ: `internal/infra/logger.go` を経由。直接 `fmt.Println` しない。
- 同期/並行: 明確な競合がない限りチャネルよりもシンプルなロック/直列処理を優先。

## 4. 命名規則
- Go: 公開シンボルは `PascalCase`、非公開は先頭小文字。略語は一般的なもののみ (`ID`, `API`).
- ファイル名: 機能 + ドメイン (`player.go`, `stats_pattern.go`). テストは同名 + `_test`。
- 定数: バックエンドは`UpperCamelCase`、フロントエンドは`SNAKE_CASE`。

## 5. ディレクトリ指針
- `internal/data/`: 外部APIレスポンスのマッピングや計算ロジック。
- `internal/di/`: 依存注入設定。
- `internal/infra/`: 外部サービス接続 (Discord, GitHub, ローカルファイル等)。
- `internal/service/`: ビジネスロジック (必要なら階層化)。
- `internal/mock/`: `go.uber.org/mock`で自動生成されたモック。この配下は手動編集禁止。
- `frontend/src/`: Svelte コンポーネント・ストア・ユーティリティ。

## 6. 依存追加
- Go: 新規ライブラリ追加前に標準 + 既存ユーティリティで代替不可か検討。`go get` → `go mod tidy` 実行。
- Node: `frontend/package.json` に追加。不要依存は削除。ライセンス確認 (OSS: MIT/Apache/BSD を優先)。

## 7. テスト戦略
- 新規ロジックには原則ユニットテスト。難読な計算は例示的ケースを複数。
- Go テスト: `testing` + `go.uber.org/mock/gomock` + `github.com/stretchr/testify`  を利用。テーブル駆動テスト推奨。
- 外部 API 呼び出しはモック (`mock/` / インターフェース) を利用。
- ファイル操作は一時ディレクトリ (`t.TempDir()`) を使用。
- 並列実行可能な場合は `t.Parallel()` を利用。
- テスト名は日本語で正常系/異常系を明示。
- フロントは状態計算/ストアロジックを Jest でテスト。UI のみの小変更多発は Snapshot 多用しすぎない。

## 8. AI / Copilot 利用ルール
- 提案コードは丸ごと受け入れず、プロジェクト構造・命名整合を確認。
- ライセンス不明な長文は拒否。著名 OSS の特定的フレーズ大量含有なら再生成。
- 秘匿情報 (API キー・ユーザーデータ) をプロンプトに貼らない。
- モデル名を尋ねられた場合のみ "GPT-5" と回答。
- 自動生成のコメントは最小限。説明必要箇所のみ手書き。

## 9. ビルド/テスト
- `task dev` で開発サーバ起動 (フロント + バックエンド)。
- `task build` で本番ビルド生成。
- `task test` でフロントエンド/バックエンドテスト実行。

## 10. コンフィグ / 環境
- `di/config.go` 経由で集約も検討。
- `di/container.go` で依存注入設定。
- ユーザーデータ (`user_data/`): 開発中に個人情報追加しない。

## 11. 既存計算ロジック参照
- 新規指標追加時は `internal/data/rating.go`, `pr_factor.go`, `stats_pattern.go` 等の既存式/パターンを参照し整合性を確保。

## 12. 新バージョン検出
- `internal/data/new_version.go` 付近の処理を拡張する際は API レート/キャッシュを考慮。

## 13. ドキュメント更新
- 公開 API/構造変更時は `README.md` とここ (必要なら) を更新。
- 生成ツールの追加 (例: lint, format) を導入したらセクション追記。

## 14. 禁止事項
- 直接 `panic` 濫用。
- ランダムなグローバル状態追加。
- デバッグ用 `fmt.Println` の残置。
- 大量自動生成コメントや未使用コードのコミット。
- 自動生成コードの手動生成修正。
  - `model.ts`
    - `wails generate module` で再生成可能する。
  - `mock`配下のコード
    - `task gen-mock` で再生成可能する。

## 15. セキュリティ
- 例外ログに機密 (トークン等) を含めない。

## 16. 今後拡張時の AI 指示テンプレ
```
要件: <具体的機能>
制約: 既存スタイル尊重 / 最小差分 / テスト追加
出力: diff パッチのみ
注意: ライセンス安全, 秘匿情報非公開
```