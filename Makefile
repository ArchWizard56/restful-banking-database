PROJECTNAME=restful-banking-database
BUILDDIR=./bin
SRCFILES=$(shell go list -f '{{.GoFiles}}' | tr -d '[]')
RUNFLAGS=-d
BUILDFLAGS=-v
DEPENDFILE="dependencies.txt"

build: $(BUILDDIR)/$(PROJECTNAME)
.PHONY : build

$(BUILDDIR)/$(PROJECTNAME): $(SRCFILES)
	go build $(BUILDFLAGS) -o "$(BUILDDIR)/$(PROJECTNAME)" $(SRCFILES)

gorun:
	go run $(SRCFILES) $(RUNFLAGS)

run: $(BUILDDIR)/$(PROJECTNAME)
	$(BUILDDIR)/$(PROJECTNAME) $(RUNFLAGS)

clean:
	rm $(BUILDDIR)/$(PROJECTNAME)
	rm $(BUILDDIR)/accounts.db

deps:
	bash -c 'while read line; do go get -u $$line; done < $(DEPENDFILE)'

