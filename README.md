# rapid-go

## stack

- golang 1.11
- mysql (cloud sql)
- Chi (as Router)
- squirrel (as query builder)
- gin (for hot reload ,not framwork)
- docker 

## local

- package
  - 実行するのは docker container の中なので関係ないが、vscode でローカルにパッケージがないとエラーでまくってうざいので、dep でローカルにライブラリをインストールすることをお勧めします。

```bash
> cd api/src
> dep ensure
```

- build

```bash
> docker-compose build
```

- start

```bash
> docker-compose up -d
```

- test

```bash
> http://localhost:3001/ping
```

- request
  - ローカル開発の際にはリクエストに毎度 auth のトークン入れるのはめんどくさいと思うので、ローカル環境のみ/noauth/v1 という感じで、noauth を挟むことで dummy の authID を入れてリクエストを通せる。

## develop flow

- model を書く
  - DBschema と一致するであろう entitiy
- repository を書く
  - DB や外部 API に対してクエリーやリクエストを投げる処理を書く
  - データの CRUD のみを扱う
- service を書く
  - ロジックをメインに書く
  - 権限コントロールなども行う
  - 複数の repository を利用することも多々ある
  - ここで API を叩いたりクエリーを書いてはいけない
- handler を書く

  - リクエストとレスポンスの定義を行う
  - URL パラメータや auth の uid 等もここで取得する。
  - えた値を全て service の関数の引数に入れて、service を呼ぶ。
  - 必ずサービスは一つ

- dependency や routing は適宜追加してください。
- 新しいライブラリーを使いたくなった場合は都度相談して利用すること(ほぼないはず)
- むやみに goroutine を使わない。(特に DB 周りは気をつけないとすぐに connection が枯渇する)
