.PHONY: releasegen
releasegen:
	docker run \
		--volume "$(CURDIR):/go/src/github.com/kramphub/kiya" \
		--workdir "/go/src/github.com/kramphub/kiya" \
		golang:1.10.0 \
		bash -x etc/bin/releasegen.sh

# go get github.com/aktau/github-release
# export GITHUB_TOKEN=...
createrelease:
	github-release info -u kramphub -r kiya
	TAG=$(git tag -l --points-at HEAD)
	github-release release \
		--user kramphub \
		--repo kiya \
		--tag $(TAG) \
		--name "kiya" \
		--description "Kiya - secrets management tool"

uploadrelease:
	github-release upload \
		--user kramphub \
		--repo kiya \
		--tag $(shell git tag -l --points-at HEAD) \
		--name "kiya-Linux-x86_64" \
		--file release/kiya-Linux-x86_64

	github-release upload \
		--user kramphub \
		--repo kiya \
		--tag $(shell git tag -l --points-at HEAD) \
		--name "kiya-Windows-x86_64" \
		--file release/kiya-Windows-x86_64

	github-release upload \
		--user kramphub \
		--repo kiya \
		--tag $(shell git tag -l --points-at HEAD) \
		--name "kiya-Darwin-x86_64" \
		--file release/kiya-Darwin-x86_64