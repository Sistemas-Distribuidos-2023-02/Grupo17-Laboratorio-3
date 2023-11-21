docker-mv-1:
	docker-compose -f docker-compose.yml up -d informante fulcrum

docker-mv-2:
	docker-compose -f docker-compose.yml up -d informante fulcrum

docker-mv-3:
	docker-compose -f docker-compose.yml up -d vanguardia fulcrum

docker-mv-4:
	docker-compose -f docker-compose.yml up -d broker-luna

docker-down:
	docker-compose -f docker-compose.yml down

docker-clean:
	docker system prune -a
