docker run -d \
    --name some-postgres \
    -e POSTGRES_PASSWORD=mysecretpassword \
    -e PGDATA=/var/lib/postgresql/data/pgdata \
    -v /Users/chrishayles/Documents/containermounts/psql/data:/var/lib/postgresql/data \
    -p 5432:5432 \
    postgres:9.6.18