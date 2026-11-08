docker network create shoplite-net
docker volume create pgdata
docker run -d --name db --restart unless-stopped --network shoplite-net \
  -v pgdata:/var/lib/postgresql/data \
  -e POSTGRES_USER=shop -e POSTGRES_PASSWORD=ChangeMe123 -e POSTGRES_DB=shoplite \
  -p 10.0.2.20:5432:5432 postgres:16-alpine
