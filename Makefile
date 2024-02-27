all: test build

test:
	cd backend && $(MAKE) test

build:
	cd backend && $(MAKE) build

# copy:
# 	cp -f backend/bin/* workflow/bin/


