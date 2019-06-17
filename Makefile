PROJECTNAME=restful-banking-database
BUILDDIR=./bin
SRCFILES=$(shell go list -f '{{.GoFiles}}' | tr -d '[]')
FLAGS="-d"

build: $(BUILDDIR)/$(PROJECTNAME)
.PHONY : build

$(BUILDDIR)/$(PROJECTNAME): $(SRCFILES)
	go build -o "$(BUILDDIR)/$(PROJECTNAME)" $(SRCFILES)

gorun:
	go run $(SRCFILES) $(FLAGS)

run: $(BUILDDIR)/$(PROJECTNAME)
	$(BUILDDIR)/$(PROJECTNAME) $(FLAGS)

clean:
	rm $(BUILDDIR)/$(PROJECTNAME)

