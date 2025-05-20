docker run --name my-postgres \
  -e POSTGRES_USER=myuser \
  -e POSTGRES_PASSWORD=mypass \
  -e POSTGRES_DB=journal \
  -p 5432:5432 \
  -v pgdata:/var/lib/postgresql/data \
  -d ankane/pgvector
docker start my-postgres
