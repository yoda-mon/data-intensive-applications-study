
### ビルド
```sh
docker build -t di-postgres-tmp .
```


### 初期化
runするとinitdbが連番振られたファイルを順次実行してくれるらしい

```sh
docker run -d --rm \
  --name di-postgres-tmp \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=m9YfasiM \
  -p 5432:5432 \
  di-postgres-tmp
```

### 保存
ローダーでデータ保存したら
```sh
docker commit di-postgres-tmp di-postgres
```

保存したやつ再走行
```
docker run -d --rm \
  --name di-postgres \
  -p 5432:5432 \
  di-postgres-tmp
```


### PSQL
# ルートでアクセス
PGPASSWORD=m9YfasiM psql -h localhost -p 5432 -U postgres
# 一般でアクセス
PGPASSWORD=1234 psql -h localhost -p 5432 -U student
```

```sql
\l   -- データベース一覧

\dn  -- スキーマ一覧

\d   -- テーブル一覧
\d amzn_revs.multailingual_jp -- データセット用テーブル 
```
