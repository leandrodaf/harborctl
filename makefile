# Makefile mínimo
# uso: make release 1.0.1

VERSION := $(word 2,$(MAKECMDGOALS))
PREFIX  ?= v
REMOTE  ?= origin

.PHONY: release
release:
	@test -n "$(VERSION)" || (echo "uso: make release <versao>"; exit 1)
	git tag -a "$(PREFIX)$(VERSION)" -m "Release $(PREFIX)$(VERSION)"
	git push "$(REMOTE)" "$(PREFIX)$(VERSION)"

# ignora o argumento extra (ex.: 1.0.1) como alvo
%:
	@:
