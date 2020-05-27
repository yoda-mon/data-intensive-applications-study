
## PostgreSQL Streaming Replication on K8s

ストリーミングレプリケーションの挙動確認用マニフェストファイル類  
- [Let's Postgres](https://lets.postgresql.jp/documents/technical/replication/1)の構築ページを参考に、pgbenchによるテーブル作成を初回起動時に実行  
- docker for mac でKubernetesオプションONにした状態で動作確認  
- **例によって手元で動かすためのものなのでパブリックな環境にDeploy等はしないこと**  

```sh
# 起動
kubectl apply -f master -f slave -f common
# 停止
kubectl apply -f master -f slave -f common
```

```sh
# ホスト端末（mac）からMasterに接続
PGPASSWORD=mysecretpassword psql -h localhost -p 5432 -U postgres
# ホスト端末（mac）からSlaveに接続
PGPASSWORD=mysecretpassword psql -h localhost -p 5433 -U postgres

# MasterのPodからSlaveに対し接続
kubectl exec -it {Master Pod名} bash
PGPASSWORD=mysecretpassword psql -h postgres-sr-slave-lb.default.svc.cluster.local -p 5433 -U postgres
```

```sh
# pgbenchによる更新
kubectl exec -i {Master Pod名} -- pgbench -p 5432 -T 180 -U postgres
# master側のレプリケーションの情報のチェック
PGPASSWORD=mysecretpassword psql -h localhost -p 5432 -U postgres -c "SELECT * FROM pg_stat_replication"
# slave側のレプリケーション最終更新時刻の取得
PGPASSWORD=mysecretpassword psql -h localhost -p 5433 -U postgres -c "SELECT pg_last_xact_replay_timestamp()"
```


### Note
雑多ノート: [NOTE.md](./NOTE.md)

### 参考リンク
- https://wiki.postgresql.org/wiki/Streaming_Replication 
- https://stacksoft.io/blog/postgres-statefulset/ Statefulsetの例。Postgres 10なところ以外は参考になる