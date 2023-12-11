:: starts the system under test
docker network create restaurant

call file-store-svc-build.cmd

docker compose -f sut/docker-compose.yml up --build --remove-orphans --exit-code-from restaurant-file-store-svc

docker compose -f sut/docker-compose.yml down

docker image rm restaurant-file-store-svc
pause