:: starts the database
docker network create restaurant

docker compose -f sut/database-compose.yml up --build --remove-orphans --exit-code-from mongodb

docker compose -f sut/database-compose.yml down

docker image rm restaurant-file-store-db
pause