run: stop up

mod:
	go mod tidy

test:
	docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit
	docker-compose -f docker-compose.test.yml down --volumes

test-db-up:
	docker-compose -f docker-compose.test.yml up --build db

test-db-down:
	docker-compose -f docker-compose.test.yml down --volumes db

up:
	docker-compose -f docker-compose.yml up -d --build

stop:
	docker-compose -f docker-compose.yml stop
