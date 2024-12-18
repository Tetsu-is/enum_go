
1. 偶数型を実装してそれを検査するanalyzerを作成する
2. enum型を実装してそれを検査するanalyzerを作成する

### 使用するときのイメージ

本ツールの行うこと

- satisfy関数によって代入文の検査を行う

ユーザーが行うこと

- カスタム型の定義
- カスタム型の値が満たすべきsatisfy関数の定義


### 偶数型のマイルストーン
- [x] var宣言をBasicLitで初期化
- [ ] 単純な代入文で右辺にUnderlyingTypeのBasicLitを代入
- [ ] for文のiteratorでvar宣言をBasicLitで初期化