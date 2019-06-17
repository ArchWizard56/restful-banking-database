PROJECTNAME=restful-banking-database
BUILDDIR=./bin
SRCFILES=$(shell go list -f '{{.GoFiles}}' | tr -d '[]')

build: $(BUILDDIR)/$(PROJECTNAME)
.PHONY : build

$(BUILDDIR)/$(PROJECTNAME): $(SRCFILES)
	go build -o "$(BUILDDIR)/$(PROJECTNAME)" $(SRCFILES)

gorun:
	go run $(SRCFILES)

run: $(BUILDDIR)/$(PROJECTNAME)
	$(BUILDDIR)/$(PROJECTNAME)

clean:
	rm $(BUILDDIR)/$(PROJECTNAME)

