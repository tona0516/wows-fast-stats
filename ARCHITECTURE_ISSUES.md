# ソフトウェアアーキテクチャにおける問題点

このドキュメントでは、`wows-fast-stats` プロジェクトにおけるソフトウェアアーキテクチャ上の問題点を詳細に分析し、改善提案を提示します。

## 概要

本リポジトリは World of Warships の統計表示アプリケーションで、Go バックエンドと Svelte フロントエンドで構成されています。アーキテクチャ検証ツール (`arch-go`) により基本的なレイヤー分離は守られていますが、いくつかの構造的な問題が存在します。

---

## 問題1: 肥大化したServiceクラス（God Object アンチパターン）

### 現状

**ファイル:** `backend/service/battle_fetcher.go` (631行)

`BattleFetcher` クラスが以下の複数の責務を持っています：

1. **データ取得**: 複数の外部API（Wargaming, Numbers, Unofficial Wargaming）からのデータ取得
2. **データ変換**: 取得したデータの変換・マッピング
3. **データ組み立て**: 最終的な `Battle` オブジェクトの組み立て
4. **並行処理の制御**: goroutine とチャネルによる並行処理
5. **状態管理**: `isFirstBattle`, `warship`, `allExpectedStats` などの内部状態
6. **エラーハンドリング**: 各API呼び出しのエラー処理
7. **イベント発行**: フロントエンドへのイベント通知

```go
type BattleFetcher struct {
	ctx            context.Context
	wargaming      repository.WargamingInterface
	uwargaming     repository.UnofficialWargamingInterface
	numbers        repository.NumbersInterface
	persistence    repository.PersistenceInterface
	logger         repository.LoggerInterface
	eventsEmitFunc eventEmitFunc

	isFirstBattle                      bool
	isNotifyExpectedStatsUnavaillalble bool  // Note: typo in source (Unavaillalble)
	warship                            data.Warships
	allExpectedStats                   data.ExpectedStats
	battleArenas                       data.WGBattleArenas
	battleTypes                        data.WGBattleTypes
}
```

### 問題点

- **単一責任原則（SRP）違反**: 1つのクラスが多すぎる責務を持つ
- **保守性の低下**: 変更が他の機能に影響を与えやすい
- **テスト困難**: 631行のクラスを包括的にテストするのは困難
- **再利用性の欠如**: データ取得ロジックを他の場所で再利用できない

### 改善提案

以下のように責務を分離：

```go
// 1. データ取得専門
type BattleDataFetcher struct {
    wargaming  repository.WargamingInterface
    uwargaming repository.UnofficialWargamingInterface
    numbers    repository.NumbersInterface
}

// 2. データキャッシュ管理
type BattleDataCache struct {
    warships         data.Warships
    expectedStats    data.ExpectedStats
    battleArenas     data.WGBattleArenas
    battleTypes      data.WGBattleTypes
    persistence      repository.PersistenceInterface
}

// 3. Battle組み立て専門
type BattleComposer struct {
    cache  *BattleDataCache
}

// 4. オーケストレーション（薄いレイヤー）
type BattleService struct {
    fetcher  *BattleDataFetcher
    cache    *BattleDataCache
    composer *BattleComposer
    logger   repository.LoggerInterface
    events   EventEmitter
}
```

**期待される効果:**
- 各クラスが100-200行程度に収まる
- 個別にテスト可能
- 責務が明確で理解しやすい

---

## 問題2: 手動依存性注入（DI）の複雑性

### 現状

**ファイル:** `dependency_container.go` (104行)

すべての依存関係を手動で構築・配線しています：

```go
func NewDependencyContainer(ctx context.Context, config Config) *DependencyContainer {
	alertDiscord := infra.NewDiscord(...)
	infoDiscord := infra.NewDiscord(...)
	persistence := infra.NewPersistence(...)
	ownIGN, _ := persistence.LoadOwnIGN()
	logger := infra.NewLogger(...)
	logger.SetOwnIGN(ownIGN)
	logger.Init(ctx)
	wargaming := infra.NewWargaming(...)
	// ... さらに続く
	configService := service.NewSetting(...)
	battleFetcher := service.NewBattleFetcher(...)
	// ...
}
```

### 問題点

- **拡張性の低さ**: 新しい依存関係追加時に大量の変更が必要
- **テスト困難**: 部分的なモック化が難しい
- **依存関係の可視化**: 依存グラフが見えにくい
- **初期化順序**: 手動で管理する必要があり、バグの温床

### 改善提案

**オプション1: Wire（Google製DIライブラリ）の導入**

```go
// +build wireinject

func InitializeApp(ctx context.Context, config Config) (*App, error) {
    wire.Build(
        // Infrastructure
        infra.NewDiscord,
        infra.NewPersistence,
        infra.NewLogger,
        infra.NewWargaming,
        // Services
        service.NewBattleService,
        service.NewConfig,
        // App
        NewApp,
    )
    return nil, nil
}
```

**オプション2: シンプルなBuilder パターン**

```go
type AppBuilder struct {
    config Config
    ctx    context.Context
}

func (b *AppBuilder) BuildInfra() *InfraContainer { ... }
func (b *AppBuilder) BuildServices(infra *InfraContainer) *ServiceContainer { ... }
func (b *AppBuilder) Build() (*App, error) { ... }
```

**期待される効果:**
- 依存関係が自動的に解決される
- テスト時の柔軟性向上
- コンパイル時の依存関係検証

---

## 問題3: ステートフルサービスと並行処理の安全性

### 現状

`BattleFetcher` が可変状態を持っています：

```go
type BattleFetcher struct {
    // ...
    isFirstBattle                      bool
    isNotifyExpectedStatsUnavaillalble bool  // Note: typo in source (Unavaillalble)
    warship                            data.Warships
    allExpectedStats                   data.ExpectedStats
    battleArenas                       data.WGBattleArenas
    battleTypes                        data.WGBattleTypes
}
```

`Invoke` メソッド内でこれらの状態を更新：

```go
func (b *BattleFetcher) Invoke(tempArenaInfo data.TempArenaInfo) {
    if b.isFirstBattle {
        // 初回のみキャッシュをロード
        b.warship = warshipResult.Value
        b.allExpectedStats = expectedStatsResult.Value
        // ...
    }
    b.isFirstBattle = false  // 状態変更
}
```

### 問題点

- **並行処理の安全性**: 複数のgoroutineから同時に呼ばれた場合、データ競合の可能性
- **予測不可能性**: 呼び出し順序に依存する動作
- **テスト困難**: 状態依存のテストが複雑
- **関数型原則違反**: 副作用が大きい

### 改善提案

**ステートレス化とキャッシュの分離:**

```go
// キャッシュを独立した構造体に
type BattleDataCache struct {
    mu               sync.RWMutex
    warships         data.Warships
    expectedStats    data.ExpectedStats
    battleArenas     data.WGBattleArenas
    battleTypes      data.WGBattleTypes
    initialized      bool
}

func (c *BattleDataCache) Get() (CachedData, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return CachedData{...}, c.initialized
}

func (c *BattleDataCache) Initialize(ctx context.Context, fetcher Fetcher) error {
    c.mu.Lock()
    defer c.mu.Unlock()
    if c.initialized {
        return nil
    }
    // 初期化処理
    c.initialized = true
    return nil
}

// BattleServiceはステートレスに
type BattleService struct {
    cache    *BattleDataCache  // 共有キャッシュ
    fetcher  BattleDataFetcher
    composer BattleComposer
}

func (s *BattleService) FetchBattle(ctx context.Context, info data.TempArenaInfo) (data.Battle, error) {
    // 状態を持たない純粋な処理
    cachedData, ok := s.cache.Get()
    if !ok {
        if err := s.cache.Initialize(ctx, s.fetcher); err != nil {
            return data.Battle{}, err
        }
        cachedData, _ = s.cache.Get()
    }
    
    // ステートレスな処理
    playerData := s.fetcher.FetchPlayerData(ctx, info)
    return s.composer.Compose(info, playerData, cachedData), nil
}
```

**期待される効果:**
- 並行処理に対して安全
- テストが容易（状態依存なし）
- 予測可能な動作

---

## 問題4: レイヤー間の関心の分離不足

### 現状

**`app.go`**: プレゼンテーション層とビジネスロジック層の中間層として動作

```go
func (a *App) SubscribeBattle() {
    // UIイベント処理
    if !a.container.battlePublisher.CanSubcribe() {
        return
    }
    // ...
    // ビジネスロジックの直接呼び出し
    for tempArenaInfo := range channel {
        a.container.battleService.Invoke(tempArenaInfo)
    }
}

func (a *App) SaveUserConfig(config data.UserConfig) error {
    return a.container.configService.SaveUserConfig(config)
}
```

### 問題点

- **レイヤー境界が曖昧**: UIとビジネスロジックが密結合
- **テスト困難**: Wails ランタイムに依存
- **再利用性の欠如**: CLI版やAPI版を作る場合に大規模な変更が必要

### 改善提案

**Application Service 層の導入:**

```go
// アプリケーション層（ユースケース）
package application

type BattleUseCase struct {
    battleService  *service.BattleService
    publisher      *service.BattlePublisher
    eventBus       EventBus
}

func (uc *BattleUseCase) StartBattleMonitoring(ctx context.Context) error {
    if !uc.publisher.CanSubscribe() {
        return ErrCannotSubscribe
    }
    
    channel := make(chan data.TempArenaInfo)
    go uc.publisher.Subscribe(ctx, channel)
    
    for info := range channel {
        battle, err := uc.battleService.FetchBattle(ctx, info)
        if err != nil {
            uc.eventBus.Emit(EventError, err)
            continue
        }
        uc.eventBus.Emit(EventBattleReady, battle)
    }
    return nil
}

// Wails UI Adapter（薄いアダプター層）
package ui

type WailsApp struct {
    ctx      context.Context
    useCase  *application.BattleUseCase
}

func (a *WailsApp) SubscribeBattle() {
    err := a.useCase.StartBattleMonitoring(a.ctx)
    if err != nil {
        runtime.EventsEmit(a.ctx, "error", err.Error())
    }
}
```

**期待される効果:**
- ビジネスロジックがフレームワークから独立
- 複数のUI（CLI、Web API、Desktop）をサポート可能
- ユニットテストが容易

---

## 問題5: エラーハンドリングの一貫性欠如

### 現状

複数のエラーパターンが混在：

```go
// パターン1: failure.New
return failure.New(apperr.InvalidInstallPath)

// パターン2: failure.Wrap
return failure.Wrap(err)

// パターン3: 無視
ownIGN, _ := persistence.LoadOwnIGN()

// パターン4: エラーコード文字列変換
b.eventsEmitFunc(b.ctx, EventErr, apperr.ToStringCode(err))

// パターン5: apperr.Unwrap
return apperr.Unwrap(err)
```

### 問題点

- **一貫性の欠如**: エラーハンドリング方法がバラバラ
- **コンテキスト情報の喪失**: ラップせずに無視されるエラー
- **デバッグ困難**: エラーの発生源が追跡しにくい

### 改善提案

**統一されたエラーハンドリング戦略:**

```go
// 1. エラー型の明確化
package apperr

type AppError struct {
    Code    ErrorCode
    Message string
    Cause   error
    Context map[string]interface{}
}

func (e *AppError) Error() string {
    return fmt.Sprintf("[%s] %s: %v", e.Code, e.Message, e.Cause)
}

func New(code ErrorCode, message string) *AppError {
    return &AppError{Code: code, Message: message}
}

func Wrap(err error, code ErrorCode, message string) *AppError {
    return &AppError{
        Code:    code,
        Message: message,
        Cause:   err,
    }
}

// 2. 明示的なエラーハンドリング方針
// - 無視してよいエラーにはコメント必須
// - ビジネスロジックのエラーは常にラップ
// - ログには常にコンテキスト情報を付与

// 例:
ownIGN, err := persistence.LoadOwnIGN()
if err != nil {
    // 初回起動時は存在しないため無視
    logger.Debug("Own IGN not found (expected on first run)", nil)
    ownIGN = ""
}
```

**期待される効果:**
- エラーのトレーサビリティ向上
- 一貫したエラーレスポンス
- デバッグの効率化

---

## 問題6: テストカバレッジの不足

### 現状

```bash
$ go test -cover ./...
wfs/backend/data       coverage: [low]
wfs/backend/service    coverage: 3 tests only
```

特に以下が未テスト：
- `BattleFetcher.compose` (120行の複雑なロジック)
- `BattleFetcher.Invoke` (メインフロー)
- 並行処理のエッジケース

### 問題点

- **リグレッションリスク**: 変更時の既存機能破壊
- **リファクタリング困難**: 安全にリファクタリングできない
- **ドキュメントの欠如**: テストが仕様書の役割を果たせない

### 改善提案

**1. テスト戦略の確立:**

```go
// テーブル駆動テスト
func TestBattleComposer_Compose(t *testing.T) {
    tests := []struct {
        name     string
        input    ComposerInput
        expected data.Battle
    }{
        {
            name: "通常戦闘_両チーム完全データ",
            input: ComposerInput{...},
            expected: data.Battle{...},
        },
        {
            name: "非表示プロフィール含む",
            input: ComposerInput{...},
            expected: data.Battle{...},
        },
        // ...
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            composer := NewBattleComposer()
            result := composer.Compose(tt.input)
            assert.Equal(t, tt.expected, result)
        })
    }
}

// 統合テスト
func TestBattleService_Integration(t *testing.T) {
    // モックAPIサーバーを使用
    mockServer := httptest.NewServer(...)
    defer mockServer.Close()
    
    service := NewBattleService(...)
    result, err := service.FetchBattle(ctx, testArenaInfo)
    
    assert.NoError(t, err)
    assert.NotEmpty(t, result.Teams)
}
```

**2. テストカバレッジ目標:**
- クリティカルパス: 90%以上
- ビジネスロジック: 80%以上
- インフラ層: 70%以上

**期待される効果:**
- 安心してリファクタリング可能
- バグの早期発見
- コードの自己文書化

---

## 問題7: 並行処理の複雑性

### 現状

`BattleFetcher.Invoke` で多数のgoroutineとチャネルを使用：

```go
accountInfoResult := make(chan data.Result[data.WGAccountInfo])
shipStatsResult := make(chan data.Result[data.AllPlayerShipsStats])
clanResult := make(chan data.Result[data.Clans])
shipsBadgesResult := make(chan data.Result[data.AllPlayerShipsBadges])

go b.fetchAccountInfo(accountIDs, accountInfoResult)
go b.fetchAllPlayerShipsStats(accountIDs, shipStatsResult)
go b.fetchAllPlayerShipsBadges(accountIDs, shipsBadgesResult)
go b.fetchClan(accountIDs, clanResult)

// チャネルから順次受信
accountInfo := <-accountInfoResult
shipStats := <-shipStatsResult
// ...
```

### 問題点

- **エラー処理の複雑性**: 一部が失敗しても全チャネル待機
- **リソースリーク**: チャネルのクローズ管理が不明確
- **デッドロックリスク**: チャネル受信順序に依存
- **可読性**: 並行処理の意図が不明確

### 改善提案

**errgroup を使用した改善:**

```go
import "golang.org/x/sync/errgroup"

func (s *BattleService) FetchBattle(ctx context.Context, info data.TempArenaInfo) (data.Battle, error) {
    var (
        accountInfo  data.WGAccountInfo
        shipStats    data.AllPlayerShipsStats
        clan         data.Clans
        shipsBadges  data.AllPlayerShipsBadges
    )
    
    g, ctx := errgroup.WithContext(ctx)
    
    // 並行フェッチ（エラー時は自動キャンセル）
    g.Go(func() error {
        var err error
        accountInfo, err = s.fetcher.FetchAccountInfo(ctx, accountIDs)
        return err
    })
    
    g.Go(func() error {
        var err error
        shipStats, err = s.fetcher.FetchShipStats(ctx, accountIDs)
        return err
    })
    
    g.Go(func() error {
        var err error
        clan, err = s.fetcher.FetchClan(ctx, accountIDs)
        return err
    })
    
    g.Go(func() error {
        var err error
        shipsBadges, err = s.fetcher.FetchShipsBadges(ctx, accountIDs)
        return err
    })
    
    // 全て完了またはエラーを待つ
    if err := g.Wait(); err != nil {
        return data.Battle{}, fmt.Errorf("failed to fetch battle data: %w", err)
    }
    
    // データ組み立て
    return s.composer.Compose(info, accountInfo, shipStats, clan, shipsBadges, s.cache.GetAll()), nil
}
```

**期待される効果:**
- エラー時の自動キャンセル
- デッドロック防止
- 可読性の向上
- リソースリークの防止

---

## 優先順位付き改善ロードマップ

### フェーズ1: 基盤整備（短期 - 1-2週間）
1. ✅ アーキテクチャ問題の文書化（このドキュメント）
2. テストインフラの整備
   - テストヘルパー関数の作成
   - モックデータの整備
3. エラーハンドリングの統一
   - `apperr` パッケージの改善
   - ガイドライン作成

### フェーズ2: サービス層のリファクタリング（中期 - 2-4週間）
1. `BattleFetcher` の分割
   - `BattleDataFetcher` 抽出
   - `BattleComposer` 抽出
   - `BattleDataCache` 抽出
2. 並行処理の改善（`errgroup` 導入）
3. ステートレス化

### フェーズ3: DI改善（中期 - 1-2週間）
1. Builder パターン導入または Wire 検討
2. `DependencyContainer` のリファクタリング

### フェーズ4: レイヤー分離（長期 - 3-4週間）
1. Application Service 層の追加
2. UI Adapter の分離
3. ドメインロジックの明確化

### フェーズ5: 品質向上（継続的）
1. テストカバレッジの向上（目標80%以上）
2. ドキュメントの整備
3. パフォーマンスプロファイリング

---

## 測定可能な改善目標

### コード品質メトリクス
- [ ] 関数あたりの平均行数: 50行以下
- [ ] 最大関数行数: 200行以下
- [ ] 循環的複雑度: 15以下
- [ ] テストカバレッジ: 80%以上

### アーキテクチャメトリクス
- [ ] レイヤー間の依存関係違反: 0件（arch-go で検証）
- [ ] God Object (300行以上のクラス): 0個
- [ ] 公開インターフェースあたりの実装数: 1-2個

### 保守性メトリクス
- [ ] 新機能追加時の変更ファイル数: 平均5ファイル以下
- [ ] ビルド時間: 30秒以内
- [ ] テスト実行時間: 10秒以内

---

## 結論

本プロジェクトは基本的なレイヤー分離は守られているものの、以下の改善が必要です：

1. **即座に対処すべき問題:**
   - God Object (`BattleFetcher`) の分割
   - ステートフル設計の改善
   - テストカバレッジの向上

2. **段階的に対処すべき問題:**
   - DI の改善
   - レイヤー分離の明確化
   - エラーハンドリングの統一

3. **継続的改善項目:**
   - ドキュメントの整備
   - パフォーマンス最適化
   - 技術的負債の削減
   - コードの typo 修正（例: `isNotifyExpectedStatsUnavaillalble` → `isNotifyExpectedStatsUnavailable`）

これらの改善により、コードの保守性、テスタビリティ、拡張性が大幅に向上し、長期的な開発速度の維持が可能になります。
