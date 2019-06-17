PROJECTNAME=restful-banking-database
BUILDDIR=./bin
SRCFILES := $(go list -f '{{.GoFiles}}' | tr -d "[]")

build: $(BUILDDIR)/$(PROJECTNAME)
.PHONY : build

$(BUILDDIR)/$(PROJECTNAME): $(SRCFILES)
	echo $(SRCFILES)
	go build -o "$(BUILDDIR)/$(PROJECTNAME)" $(SRCFILES)

gorun:
	go run $(SRCFILES)

run: $(BUILDDIR)/$(PROJECTNAME)
	$(BUILDDIR)/$(PROJECTNAME)

clean:
	rm $(BUILDDIR)/$(PROJECTNAME)

