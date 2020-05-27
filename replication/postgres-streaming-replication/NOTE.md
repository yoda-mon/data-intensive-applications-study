- [公式PostgreSQLのイメージ](https://hub.docker.com/_/postgres)は`$PGDATA`以下にデータベースがあれば`initdb`をスキップ、なければ`initdb`を行ったのち、`/docker-entrypoint-initdb.d/`以下にあるスクリプトやsqlを順次実行という挙動をとる
  - `$PGDATA`以下に中途半端なファイルだけ（`postgres.conf`だけとか）があった場合、エラーを出力して終了してしまう
- ストリーミングレプリケーションにおいては、`initdb`時に決定されるdatabase identifierがマスター、スレーブで共通している必要がある（別々に生成したDBでレプリケーションを組むのが無理そう）

なので

- `postres.conf`、`pg_hba.conf`　（`replicatiion.conf`はPostgreSQL 12よりdepricateされたらしく、`postgres.conf`に統合された）については、
  1. initContainerでConfigMapで定義したものをemptyDirにコピーし編集
  2. containerでそのemptyDirをマウント 、起動時オプションで指定してPostgresに読み込ませる
- configMapを直接マウントして編集しようとするとhttps://stackoverflow.com/questions/51884999/chown-var-lib-postgresql-data-postgresql-conf-read-only-file-system
みたいなエラーが出る。configMapがReadOnlyであるため。
- スレーブ側はinitContainerでマスター側DBからのバックアップを、スレーブ本体のPostgres起動前に取得し`$PGDATA`以下に展開する（この際configと同様にemptyDirを用いてinitContainerからcontainerへ受け渡し）。これによりレプリケーション可能な状態（database identifierが同じ）での起動ができる

- 現在masterで`FATAL:  archive command failed with exit code 127 DETAIL:  The failed archive command was: ...`と出る。WALのアーカイブに失敗してる？
