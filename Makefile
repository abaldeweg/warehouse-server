release:
	git tag -a v$(TAG) -m "v$(TAG)"
	git push origin v$(TAG)

	git tag -a blog/v$(TAG) -m "v$(TAG)"
	git push origin blog/v$(TAG)

	git tag -a admincli/v$(TAG) -m "v$(TAG)"
	git push origin admincli/v$(TAG)

	git tag -a framework/v$(TAG) -m "v$(TAG)"
	git push origin framework/v$(TAG)

	git tag -a gateway/v$(TAG) -m "v$(TAG)"
	git push origin gateway/v$(TAG)
