#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username student --dbname student <<-EOSQL
  CREATE TABLE amzn_revs.multailingual_jp (
    marketplace       varchar(2),
    customer_id       integer,
    review_id         varchar(14),
    product_id        varchar(10),
    product_parent    integer,
    product_title     varchar(400),
    product_category  varchar(40),
    star_rating       integer,
    helpful_votes     integer,
    total_votes       integer,
    vine              varchar(1),
    verified_purchase varchar(1),
    review_headline   varchar(300),
    review_body       text,
    review_date       date
  );
  
  -- \COPY amzn_revs.multailingual_jp FROM '/docker-entrypoint-initdb.d/amazon_reviews_multilingual_JP_v1_00.csv' (FORMAT csv, HEADER true);
EOSQL