
SHELL=/bin/bash


.PHONY: install
install:
	go install go.osspkg.com/goppy/v3/cmd/goppy@latest
	goppy setup-lib

.PHONY: lint
lint:
	goppy lint

.PHONY: npm-lint
npm-lint:
	pnpm lint

.PHONY: npm-format
npm-format:
	pnpm format

.PHONY: npm-format-check
npm-format-check:
	pnpm format:check

.PHONY: license
license:
	goppy license

.PHONY: build
build:
	goppy build --arch=amd64

.PHONY: tests
tests:
	goppy test

.PHONY: pre-commit
pre-commit: install license lint tests build

.PHONY: ci
ci: pre-commit

.PHONY: examples-dev
examples-dev:
	pnpm --filter @osspkg/ui-examples dev

.PHONY: npm-install
npm-install:
	pnpm install --frozen-lockfile

.PHONY: npm-build
npm-build: npm-install
	pnpm --filter @osspkg/ui-core build
	pnpm --filter @osspkg/ui-transport build
	pnpm --filter @osspkg/ui-react build

.PHONY: npm-publish
npm-publish: npm-build
	pnpm --filter @osspkg/ui-core publish --access public
	pnpm --filter @osspkg/ui-transport publish --access public
	pnpm --filter @osspkg/ui-react publish --access public
