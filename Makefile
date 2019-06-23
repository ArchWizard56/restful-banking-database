PROJECTNAME=restful-banking-database
BUILDDIR=./bin
RUNFLAGS=-d
BUILDFLAGS=-v
OSFLAG='linux'
SRCFILES=$(shell env GOOS=$(OSFLAG) go list -f '{{.GoFiles}}' | tr -d '[]')
DEPENDFILE="dependencies.txt"

build: $(BUILDDIR)/$(PROJECTNAME)
.PHONY : build

$(BUILDDIR)/$(PROJECTNAME): $(SRCFILES) deps
	env CGO_ENABLED=1 GOOS=$(OSFLAG) go build $(BUILDFLAGS) -o "$(BUILDDIR)/$(PROJECTNAME)" $(SRCFILES)

gorun:
	go run $(SRCFILES) $(RUNFLAGS)

run: $(BUILDDIR)/$(PROJECTNAME)
	$(BUILDDIR)/$(PROJECTNAME) $(RUNFLAGS)

clean:
	rm $(BUILDDIR)/$(PROJECTNAME)
	rm $(BUILDDIR)/accounts.db

.PHONY: deps
deps:
	bash -c 'while read line; do go get -u $$line; done < $(DEPENDFILE)'

