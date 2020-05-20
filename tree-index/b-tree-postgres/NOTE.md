ロード完了までが前提  
JOINはデータ用意するのが面倒（`product_id`がどうやらASINみたいなので、スクレイプして値段取ってくるとかすれば面白そう）  
なのでproductだけ分離して試してみる  


```sql
-- productテーブル
CREATE TABLE amzn_revs.multailingual_jp_products AS
SELECT
  DISTINCT(product_id),
  product_parent,
  product_title,
  product_cate前提  
FROM
  amzn_revs.multailingual_jp
;

-- reviewテーブル
CREATE TABLE amzn_revs.multailingual_jp_reviews AS
SELECT
  DISTINCT(review_id),
  marketplace,
  customer_id,
  product_id,
  star_rating,
  helpful_votes,
  total_votes,
  vine,
  verified_purchase,
  review_headline,
  review_body,
  review_date
FROM
  amzn_revs.multailingual_jp
;

-- シンプルにJOIN
SELECT
  p.product_id,
  p.product_title,
  r.star_rating
FROM amzn_revs.multailingual_jp_reviews AS r
LEFT JOIN amzn_revs.multailingual_jp_products AS p
ON r.product_id = p.product_id
;

-- 星が多い商品と、その星ごとの数
SELECT 
  by_product_star.product_id, 
  p.product_title,
  by_product_star.star_rating, 
  by_product_star.stars_by_rate, 
  by_product.stars_total 
FROM (
  SELECT
    product_id,
    star_rating,
    COUNT(star_rating) AS stars_by_rate
  FROM 
    amzn_revs.multailingual_jp_reviews
  GROUP BY
    product_id,
    star_rating
  ) AS by_product_star
LEFT JOIN (
  SELECT
    product_id,
    COUNT(star_rating) AS stars_total
  FROM 
    --by_product_star
    amzn_revs.multailingual_jp_reviews
  GROUP BY
    product_id
  ) AS by_product
ON by_product_star.product_id = by_product.product_id
LEFT JOIN 
  amzn_revs.multailingual_jp_products AS p
ON by_product_star.product_id = p.product_id
ORDER BY 
  stars_total DESC,
  star_rating DESC
;

```

どれも割と一瞬で終わって試験にならなそう...て事でメモリを制限してみる

```sh
docker commit di-postgres di-postgres-r1
docker stop di-postgres
docker run -d --rm \
   --name di-postgres-r1 \
   -p 5432:5432 \
   --cpuset-cpus=0 \
   --memory=128m \
   di-postgres-r1
docker stats # 確認
```

うーんまぁ取り敢えずこれで確認してみよう

```sql
-- 実行時間表示
\timing
-- 出力を見せない
\o /dev/null
```

- シンプル Time: 1103.003 ms (00:01.103)
- ちょっと複雑 Time: 1161.049 ms (00:01.161)

```SQL
CREATE INDEX ON amzn_revs.multailingual_jp_reviews (product_id);
CREATE INDEX ON amzn_revs.multailingual_jp_products (product_id);
```
- シンプル　1230.988 ms (00:01.231)
- 複雑 Time: 1168.812 ms (00:01.169)

うーんキャッシュ載っちゃってるかな
`"echo 3 > /proc/sys/vm/drop_caches"`でできるようだけどコンテナだとホスト側になっちゃうか。。。

データサイズについて
`SELECT datname, pg_size_pretty(pg_database_size(datname)) FROM pg_database;` でデータベース単位

https://qiita.com/awakia/items/99c3d114aa16099e825d の謎クエリ走らせるとテーブル単位で表示される

```sh
# Index貼る前
amzn_revs | multailingual_jp_reviews  | 208 MB     | 
amzn_revs | multailingual_jp          | 197 MB     | 
amzn_revs | multailingual_jp_products | 3672 kB    |

# Index貼る後
amzn_revs | multailingual_jp_reviews                 | 208 MB     | pg_toast_24589           |                          |       24589 | r       |    262414 |    26681
amzn_revs | multailingual_jp                         | 197 MB     | pg_toast_16387           |                          |       16387 | r       |    245331 |    25276
amzn_revs | multailingual_jp_reviews_product_id_idx  | 8120 kB    |                          |                          |       32777 | i       |    262414 |     1015
amzn_revs | multailingual_jp_products                | 3672 kB    |                          |                          |       24585 | r       |     38461 |      459
amzn_revs | multailingual_jp_products_product_id_idx | 1200 kB    |
```

インデックス貼った分増加している。

- チューニングしてPostgres自体の持つメモリを減らす
- データ自体を水増しする
- Diskが遅いやつ使う

などすれば差が見れるかな