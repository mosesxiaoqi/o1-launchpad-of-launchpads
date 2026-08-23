.PHONY: test test-contracts test-server test-web secrets-check dev dev-rpc dev-api dev-indexer dev-web smoke migrate

test: secrets-check test-contracts test-server test-web

test-contracts:
	@test ! -d contracts/test || (cd contracts && forge test)

test-server:
	@test ! -f server/go.mod || (cd server && go test ./... -count=1)

test-web:
	@test ! -f web/package.json || (cd web && npm test -- --run && npm run lint && npm run build && npm run test:e2e)

dev:
	@$(MAKE) -j4 dev-rpc dev-api dev-indexer dev-web

dev-rpc:
	@cd server/service/launchpad/rpc && go run . -f etc/launchpad-rpc.yaml

dev-api:
	@cd server/service/launchpad/api && go run . -f etc/launchpad-api.yaml

dev-indexer:
	@cd server/service/indexer && go run . -f etc/indexer.yaml

dev-web:
	@cd web && npm run dev

migrate:
	@test -n "$(DATABASE_URL)" || (echo 'DATABASE_URL is required' >&2; exit 1)
	@psql "$(DATABASE_URL)" -v ON_ERROR_STOP=1 -f server/migrations/001_init.sql

smoke:
	@./scripts/demo-smoke.sh

secrets-check:
	@! git ls-files -co --exclude-standard -z | xargs -0 grep -nE '(PRIVATE_KEY=0x|o1_launch_[a-f0-9]{8}_)' 2>/dev/null | grep -vE '^(Makefile|docs/superpowers/plans/)'
