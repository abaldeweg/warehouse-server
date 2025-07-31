release:
	git tag -a v$(TAG) -m "v$(TAG)"
	git push origin v$(TAG)

	git tag -a framework/v$(TAG) -m "v$(TAG)"
	git push origin framework/v$(TAG)

	git tag -a gateway/v$(TAG) -m "v$(TAG)"
	git push origin gateway/v$(TAG)

	git tag -a static/v$(TAG) -m "v$(TAG)"
	git push origin static/v$(TAG)

	git tag -a logs_import/v$(TAG) -m "v$(TAG)"
	git push origin logs_import/v$(TAG)

	git tag -a logs_web/v$(TAG) -m "v$(TAG)"
	git push origin logs_web/v$(TAG)
