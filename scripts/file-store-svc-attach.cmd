:: attach to the system under test, only works if you define a CMD instead of an ENTRYPOINT
docker network create restaurant

call file-store-svc-build.cmd

docker-compose -f sut/sut-compose.yml up -d --build --remove-orphans

docker-compose -f sut/sut-compose.yml exec restaurant-file-store-svc bash

docker compose -f sut/sut-compose.yml down

docker image rm restaurant-file-store-svc

pause