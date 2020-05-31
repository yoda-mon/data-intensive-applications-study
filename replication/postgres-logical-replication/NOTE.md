CREATE PUBLICATION mypub FOR TABLE pgbench_accounts, pgbench_branches, pgbench_history, pgbench_tellers;

postgres=# CREATE SUBSCRIPTION mysub CONNECTION 'dbname=test host=postgres-lr-master-lb.default.svc.cluster.local  user=repl password=repluserpassword' PUBLICATION mypub;
CREATE SUBSCRIPTION mysub CONNECTION 'dbname=test host=postgres-lr-master-lb.default.svc.cluster.local  user=repl' PUBLICATION mypub;

pg_basebackupやらずに同名のテーブルだけ作ってレプリケーション開始すると下記のようなエラーが発生
ERROR:  duplicate key value violates unique constraint "pgbench_accounts_pkey"
DETAIL:  Key (aid)=(1953) already exists.

逆にpg_backupするとinitdbが働かなくなる（srでも実はここ実行されてなかった模様？）

http://www.intellilink.co.jp/article/column/oss-postgres03.html
列名さえ包含していればレプリケーションできる模様、なのでテーブル定義を引っ張り出す

pg_dump -U postgres -s --schema=public test

CREATE TABLE public.pgbench_accounts (
    aid integer NOT NULL,
    bid integer,
    abalance integer,
    filler character(84)
)
WITH (fillfactor='100');

CREATE TABLE public.pgbench_branches (
    bid integer NOT NULL,
    bbalance integer,
    filler character(88)
)
WITH (fillfactor='100');

CREATE TABLE public.pgbench_history (
    tid integer,
    bid integer,
    aid integer,
    delta integer,
    mtime timestamp without time zone,
    filler character(22)
);

CREATE TABLE public.pgbench_tellers (
    tid integer NOT NULL,
    bid integer,
    tbalance integer,
    filler character(84)
)
WITH (fillfactor='100');

できた