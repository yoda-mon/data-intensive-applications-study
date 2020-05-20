#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
  CREATE DATABASE student;
  CREATE ROLE student WITH LOGIN PASSWORD '1234'; 

  \connect student;
  
  CREATE SCHEMA amzn_revs AUTHORIZATION student;
  
  GRANT ALL ON ALL TABLES IN SCHEMA amzn_revs TO student;
  GRANT ALL ON ALL SEQUENCES IN SCHEMA amzn_revs TO student;
  GRANT ALL ON ALL FUNCTIONS IN SCHEMA amzn_revs TO student;

  ALTER ROLE student SET search_path TO public, amzn_revs;
EOSQL

