FROM postgres:12

#ENV POSTGRES_DB student

COPY *.sh /docker-entrypoint-initdb.d/
#COPY *.sql /docker-entrypoint-initdb.d/
#COPY *.csv /docker-entrypoint-initdb.d/

ENV PGDATA="/data"