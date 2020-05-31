
## PostgreSQL Logical Replication on K8s

ロジカル`レプリケーションの挙動確認用マニフェストファイル類  
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
PGPASSWORD=mysecretpassword psql -h postgres-lr-slave-lb.default.svc.cluster.local -p 5433 -U postgres
```

```sh
# pgbenchによる更新
kubectl exec -i {Master Pod名} -- pgbench -p 5432 -T 180 -U postgres
```


### Note
雑多ノート: [NOTE.md](./NOTE.md)

### 参考リンク
- https://www.postgresql.org/docs/10/logical-replication-quick-setup.html
- http://www.intellilink.co.jp/article/column/oss-postgres03.html