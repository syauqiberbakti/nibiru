###############################################################################
###                           Simplenet Localnet                             ###
###############################################################################

# Triggers a force rebuild of the simplenet images
.PHONY: simplenet-build
simplenet-build:
	docker compose -f ./contrib/docker-compose/docker-compose-simplenet.yml build --no-cache --pull

# Run a simplenet testnet locally
.PHONY: simplenet
simplenet: simplenet-down
	docker compose -f ./contrib/docker-compose/docker-compose-simplenet.yml up --detach --build

# Stop simplenet
.PHONY: simplenet-down
simplenet-down:
	docker compose -f ./contrib/docker-compose/docker-compose-simplenet.yml down --timeout 1 --volumes

###############################################################################
###                               simplenet Logs                             ###
###############################################################################

.PHONY: simplenet-logs
simplenet-logs:
	docker compose -f ./contrib/docker-compose/docker-compose-simplenet.yml logs

.PHONY: simplenet-logs-pricefeeder
simplenet-logs-pricefeeder:
	docker compose -f ./contrib/docker-compose/docker-compose-simplenet.yml logs pricefeeder --follow

.PHONY: simplenet-logs-go-heartmonitor
simplenet-logs-go-heartmonitor:
	docker compose -f ./contrib/docker-compose/docker-compose-simplenet.yml logs go-heartmonitor

###############################################################################
###                              simplenet SSH                               ###
###############################################################################

.PHONY: simplenet-ssh-nibiru
simplenet-ssh-nibiru:
	docker compose -f ./contrib/docker-compose/docker-compose-simplenet.yml exec -it nibiru /bin/sh

.PHONY: simplenet-ssh-go-heartmonitor
simplenet-ssh-go-heartmonitor:
	docker compose -f ./contrib/docker-compose/docker-compose-simplenet.yml exec -it go-heartmonitor /bin/sh

