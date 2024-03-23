:: starts the system under test
docker network create restaurant

call file-store-svc-build.cmd

docker compose -f sut/sut-compose.yml up --build --remove-orphans --exit-code-from restaurant-file-store-svc

docker compose -f sut/sut-compose.yml down --volumes

docker image rm restaurant-file-store-svc
pause