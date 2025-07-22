# Poker Penalty Paradise Pickem Selector - 実装方針

## 概要
このプロジェクトは、4枚のカードが配られた時に最適なポーカーゲームバリアントを選択するシステムです。各ゲームでの勝率をシミュレーションし、最も高い期待値を持つゲームを推奨します。

## アーキテクチャ

### パッケージ構成
```
pkg/poker/
├── card.go           # カードの基本表現
├── evaluator.go      # 各種ハンド評価関数
├── game.go           # ゲームインターフェースと実装
├── simulator.go      # モンテカルロシミュレーション
├── hidugi_simulator.go # HiDuGi専用シミュレーター
└── parser.go         # 入力パース処理
```

### 主要コンポーネント

#### 1. カード表現 (card.go)
- カードは0-51の整数で表現
- `Rank()`: カードのランク (0=2, 12=A)
- `Suit()`: カードのスート (0=c, 1=d, 2=h, 3=s)
- 文字列との相互変換機能

#### 2. ハンド評価 (evaluator.go)
- `Evaluate5CardHigh()`: 5枚ポーカーのハンド評価
- `Evaluate4CardHigh()`: 4枚ポーカーのハンド評価
- `EvaluateBadugi()`: バドゥーギのハンド評価
- `EvaluateHiDuGi()`: HiDuGiの複合評価

#### 3. ゲーム実装 (game.go)
各ゲームは`Game`インターフェースを実装：
- `DrawmahaHi`: ドローマハハイ
- `BadugiGame`: バドゥーギ
- `HiDuGiGame`: ハイドゥーギ（スプリットポット）
- `StubGame`: 未実装ゲームのプレースホルダー

#### 4. シミュレーション (simulator.go)
- モンテカルロ法による勝率計算
- デフォルト100,000回の試行で高精度を実現
- 並列実行による高速化

## 特殊な実装

### HiDuGiの評価ロジック
HiDuGiはハイハンドとバドゥーギハンドでポットを分け合うスプリットポットゲームです。

#### 評価戦略
1. **基本評価**: ハイスコアとバドゥーギスコアを組み合わせ
2. **8-バドゥーギボーナス**: 8以下のバドゥーギには10倍のボーナス
3. **強いハイハンドボーナス**: トリップス以上には3倍のボーナス

#### スプリットポット戦略
スプリットポットゲームでは「片方のポットを確実に取る」ことが重要です。

```go
// シミュレーション時のボーナス適用
if hasVeryStrongHigh && equity >= 0.5 {
    if highCategory >= 7 { // フォーカード
        equity = equity * 1.10  // 10%ボーナス
    } else {
        equity = equity * 1.05  // 5%ボーナス
    }
}
```

このボーナスにより、4枚のエース（AAAA）のような手は：
- Drawmaha-Hiではストレートフラッシュに負ける可能性がある
- HiDuGiではハイポットを確実に獲得できる
→ HiDuGiが選択される

## パフォーマンス最適化

### 1. 高速ハンド評価
- ビット演算を活用した効率的な評価
- カテゴリー別の早期終了ロジック

### 2. シミュレーション精度
- 100,000回の試行で統計的信頼性を確保
- 実行時間は約300-400ms（実用的な範囲）

### 3. メモリ効率
- カードは整数で表現（メモリ使用量最小化）
- 配列の事前確保でアロケーション削減

## テスト戦略

### ユニットテスト
- 各評価関数の正確性を検証
- エッジケース（ストレート、フラッシュ等）のテスト
- ゲーム選択ロジックのテスト

### 統合テスト
- 実際のハンド例での動作確認
- シミュレーション精度の検証

## 今後の拡張

### 実装予定のゲーム
- Drawmaha-2-7
- Prime
- Omaha DoubleBoard

### 機能拡張
- 複数対戦相手への対応
- より詳細な統計情報の提供
- Webインターフェースの追加

## 新しいゲームの実装方針

### 1. ゲームインターフェースの実装
新しいゲームを追加する際は、`Game`インターフェースを実装します：

```go
type MyNewGame struct{}

func (g MyNewGame) Name() string {
    return "MyNewGame"
}

func (g MyNewGame) Evaluate(hand []Card, board []Card) int64 {
    // ゲーム固有の評価ロジック
}
```

### 2. 評価関数の実装
ゲームのルールに応じて適切な評価関数を実装：

#### ハイゲームの場合
```go
func (g MyHighGame) Evaluate(hand []Card, board []Card) int64 {
    best5 := findBest5Cards(hand, board)
    return Evaluate5CardHigh(best5)
}
```

#### ローゲームの場合
```go
func (g MyLowGame) Evaluate(hand []Card, board []Card) int64 {
    best5 := findBest5Cards(hand, board)
    // ローゲームでは小さい値が強い
    return -EvaluateLow(best5)
}
```

#### スプリットポットゲームの場合
```go
func (g MySplitGame) Evaluate(hand []Card, board []Card) int64 {
    // HiDuGiGame の実装を参考に
    // 複数の評価を組み合わせる
}
```

### 3. シミュレーション戦略

#### 標準シミュレーション
ほとんどのゲームでは`SimulateEquity`をそのまま使用できます。

#### 特殊シミュレーション
スプリットポットゲームなど特殊なルールがある場合：

1. 専用のシミュレーター関数を作成
2. `SimulateEquity`内で特殊処理を追加

```go
// simulator.go
if _, ok := g.(MySplitGame); ok {
    return SimulateMySplitGameEquity(my4, iters)
}
```

### 4. テストの追加

#### 評価関数のテスト
```go
func TestMyNewGameEvaluation(t *testing.T) {
    tests := []struct {
        name     string
        hand     []Card
        expected string  // 期待されるハンドの強さ
    }{
        // テストケースを追加
    }
}
```

#### シミュレーションのテスト
```go
func TestMyNewGameSimulation(t *testing.T) {
    // 特定のハンドで期待される勝率範囲をテスト
}
```

### 5. 実装チェックリスト

新しいゲームを追加する際のチェックリスト：

- [ ] `Game`インターフェースの実装
- [ ] 評価関数の作成または既存関数の活用
- [ ] 必要に応じて専用シミュレーターの実装
- [ ] `PickBestGame`関数への追加
- [ ] ユニットテストの作成
- [ ] 特徴的なハンドでの動作確認
- [ ] ドキュメントの更新

### 6. 特殊ルールへの対応

#### ドローゲーム
- カード交換のシミュレーションが必要
- 最適な交換戦略の実装

#### ボードゲーム
- 共有カードの考慮
- `board`パラメータの活用

#### ワイルドカード
- 評価関数でワイルドカードの処理
- 最適なワイルドカードの使い方を決定

### 7. パフォーマンス考慮事項

- 評価関数は高速である必要がある（数百万回呼ばれる）
- 可能な限り事前計算やキャッシュを活用
- 不必要なメモリアロケーションを避ける