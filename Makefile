docker-mv-1:
	docker-compose -f docker-compose.yml up -d informante1 fulcrum1

docker-mv-2:
	docker-compose -f docker-compose.yml up -d informante2 fulcrum2

docker-mv-3:
	docker-compose -f docker-compose.yml up -d vanguardia fulcrum3

docker-mv-4:
	docker-compose -f docker-compose.yml up -d broker-luna

docker-down:
	docker-compose -f docker-compose.yml down

docker-clean:
	docker system prune -a
