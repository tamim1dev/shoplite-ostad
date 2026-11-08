docker network create shoplite-net
docker run -d --name backend --restart unless-stopped --network shoplite-net \
  -e DATABASE_URL='postgres://shop:ChangeMe123@10.0.2.20:5432/shoplite?sslmode=disable' \
  -p 8080:8080 YOURUSER/shoplite-backend:v1.0.0
