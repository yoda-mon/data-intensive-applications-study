amazonのレビューをPostgresに放り込む

まだ書き捨て状態なので

- ファイル名決め打ち
- 並列書き込み数決め打ち
- ロードできない行あり（カスタムのTSVパーサー作った）

に注意（そのうち修正しようかしまいか...）

### Postgresへのロード方法
di-postgres-tmpが起動している状態にしておいて、このディレクトリから

```sh
go run main.go data.go
```

もしくは
```sh
go build
./amzn-revs-loader
```

今のところ  
262,432 行 あるデータの内、SELECT COUNTすると 
262,414 行
なのでいくつかとりこぼしてる状態

ロードした後はdocker commitして保存しておくこと。


### メモ
golangのCSVReaderはtsvも対応しているものの、標準化されている？CSVの形でないと読み込んでくれない。今回のアマゾンのデータセットは

```
あいうえお\tかき"くけこ"\t\t"たち:つてと
```

のように各要素中にダブルクォーテーションが現れるのでそのままCSVReaderを使うとErrorが発生（`"`は要素の両端にしか現れてはいけない仕様）  
https://golang.org/src/encoding/csv/reader.go から元コードを取ってきて、`"`が出てきたら`\`を差しこみEscapeするようにした（`csv/`に配置）

また`'`をPosgtresはカラムとして認識してしまうので書き込み時に`pq.QuoteLiteral`で無理やりエスケープ

他にも謎エラー出てるけど取り敢えずは放置