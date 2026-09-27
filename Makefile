
SHELL=/bin/bash


.PHONY: install
install:
	pnpm install --frozen-lockfile
	go install go.osspkg.com/goppy/v3/cmd/goppy@latest
	goppy setup-lib

.PHONY: lint
lint:
	pnpm lint
	pnpm format
	pnpm format:check
	goppy lint

.PHONY: license
license:
	goppy license

.PHONY: build
build:
	pnpm --filter @osspkg/ui-core build
	pnpm --filter @osspkg/ui-transport build
	pnpm --filter @osspkg/ui-react build
	pnpm typecheck
	goppy build --arch=amd64

.PHONY: tests
tests:
	pnpm test
	goppy test

.PHONY: pre-commit
pre-commit: install license lint tests build

.PHONY: examples-dev
examples-dev:
	pnpm --filter @osspkg/ui-examples dev

.PHONY: npm-publish
npm-publish: npm-build
	pnpm --filter @osspkg/ui-core publish --access public
	pnpm --filter @osspkg/ui-transport publish --access public
	pnpm --filter @osspkg/ui-react publish --access public
