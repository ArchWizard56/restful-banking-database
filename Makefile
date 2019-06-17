PROJECTNAME=restful-banking-database
BUILDDIR=./bin
SRCFILES=$(wildcard *.go)

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

