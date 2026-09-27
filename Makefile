
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
pre-commit: install license lint build tests 

.PHONY: examples-dev
examples-dev:
	pnpm --filter @osspkg/ui-examples dev

.PHONY: npm-version-patch
npm-version-patch:
	cd react-sdk/ui-core && npm version patch --no-git-tag-version
	cd react-sdk/ui-transport && npm version patch --no-git-tag-version
	cd react-sdk/ui-react && npm version patch --no-git-tag-version

.PHONY: npm-publish
npm-publish: build npm-version-patch
	pnpm --filter @osspkg/ui-core publish --access public
	pnpm --filter @osspkg/ui-transport publish --access public
	pnpm --filter @osspkg/ui-react publish --access public
