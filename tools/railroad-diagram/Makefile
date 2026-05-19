REPO_URL := https://github.com/tabatkins/railroad-diagrams
BRANCH   := gh-pages
DEST     := ./tabatkins

clone:
	rm -rf $(DEST)
	git clone --depth=1 --branch=$(BRANCH) --single-branch $(REPO_URL) $(DEST)
	rm -rf $(DEST)/.git

run:
	@if [ ! -d $(DEST) ]; then $(MAKE) clone; fi
	go run ./main.go
