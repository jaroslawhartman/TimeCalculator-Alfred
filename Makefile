all: build copy

build:
	cd backend && $(MAKE) build

copy:
	cp -f backend/bin/* workflow/bin/


