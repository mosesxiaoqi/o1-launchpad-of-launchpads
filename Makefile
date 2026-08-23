.PHONY: test test-contracts test-server test-web secrets-check

test: test-contracts test-server test-web

test-contracts:
	@test ! -d contracts/test || (cd contracts && forge test)

test-server:
	@test ! -f server/go.mod || (cd server && go test ./... -count=1)

test-web:
	@test ! -f web/package.json || (cd web && npm test -- --run && npm run build)

secrets-check:
	@! git grep -nE '(PRIVATE_KEY=0x|o1_launch_[a-f0-9]{8}_)' -- . ':!Makefile' ':!docs/superpowers/plans'
