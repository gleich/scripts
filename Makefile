.PHONY = install

install:
	@for dir in $$(find . -type f -name '*.go' ! -path "./vendor/*" -exec grep -l '^package main' {} \; | xargs -n1 dirname | sort -u); do \
		echo "Installing package in $$dir"; \
		go install $$dir; \
	done

