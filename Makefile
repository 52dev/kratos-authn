APP_VERSION=v1.0.2

PACKAGE_LIST = engine/presharedkey/ engine/oidc/ engine/jwt/ middleware/
.PHONY: tag
tag:
	git tag -f $(APP_VERSION) && $(foreach item, $(PACKAGE_LIST), git tag -f $(item)$(APP_VERSION) && ) git push --tags --force
