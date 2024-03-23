:: attach to the system under test, only works if you define a CMD instead of an ENTRYPOINT
docker network create restaurant

call file-store-svc-build.cmd

docker compose -f sut/sut-compose.yml build

docker run -t -i --name restaurant-file-store-svc restaurant-file-store-svc bash

docker rm restaurant-file-store-svc

docker image rm restaurant-file-store-svc

pause