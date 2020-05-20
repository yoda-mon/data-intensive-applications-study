インデックスについて触ってみるため用

```sh
tree-index/
├── README.md
├── amzn-revs-loader  # ローダー
├── b-tree-postgres   # PostgreSQLのDocker用
└── data              # 利用データ格納
```

##  利用データ
Amazonのレビューのオープンデータ 
- [ドキュメント](https://s3.amazonaws.com/amazon-reviews-pds/readme.html)
  - [各カラムの説明など](https://s3.amazonaws.com/amazon-reviews-pds/tsv/index.txt)


カラム名           | 説明
:--               | :--
marketplace       | レビューの書かれたマーケットプレイスの2文字の国別コード。
customer_id       | あるひとりの人に書かれたレビューを集約することのできる、ランダムなID。
review_id         | レビューのユニークなID。
product_id        | レビューが関連する商品のユニークなID。多言語データセットでは、異なる国における同じ商品は同じproduct_idでGROUP BYできる。
product_parent    | 同じ商品に対するレビューを集約することのできる、ランダムなID。
product_title     | 商品名
product_category  | レビューを分類することのできる、商品の大きなカテゴリ（データセットを同質のものでまとめるのにも利用可能）。
star_rating       | 1-5の評価のスター。
helpful_votes     | 「参考になったレビュー」の投票数。
total_votes       | レビューに対する投票の総数。
vine              | 先取りプログラムで寄せられたレビュー。
verified_purchase | 検証された購入でのレビュー。
review_headline   | レビューのタイトル。
review_body       | レビューの本文。
review_date       | レビューの書かれた日付。

```sh
aws s3 ls s3://amazon-reviews-pds/tsv/  --human-readable

# 一番大きなEbookのレビュー取得
aws s3 cp s3://amazon-reviews-pds/tsv/amazon_reviews_us_Digital_Ebook_Purchase_v1_00.tsv.gz data/
# 小さめの日本語のレビュー取得
aws s3 cp s3://amazon-reviews-pds/tsv/amazon_reviews_multilingual_JP_v1_00.tsv.gz data/

# 一括解凍
find ./ -type f -name "*.gz" -exec gunzip {} \;

# TSV -> CSV
cat amazon_reviews_multilingual_JP_v1_00.tsv |tr "\\t" "," > amazon_reviews_multilingual_JP_v1_00.csv
```

色々やってみたけどCOPY文だとフォーマットエラーになるのでローダー作製
