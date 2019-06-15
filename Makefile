PROJECTNAME=restful-banking-database
BUILDDIR=./bin
SRCFILES=$(wildcard *.go)
build: $(SRCFILES)
	go build -o "$(BUILDDIR)/$(PROJECTNAME)" $(SRCFILES)
run:
	go run $(SRCFILES)
