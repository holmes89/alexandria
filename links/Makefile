DATE := $(shell /bin/date -d "today" +"%Y%m%d%H%M")
create-migration:
	touch resources/migrations/$(DATE).up.sql
	touch resources/migrations/$(DATE).down.sql
