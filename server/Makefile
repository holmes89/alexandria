.PHONY: start-db sample-call

start-db:
	docker run -d -p 5432:5432 -e POSTGRES_PASSWORD=password postgres

sample-call:
	curl -X POST localhost:8080/books/ -F 'file=@/Users/joel/Downloads/the-joy-of-clojure.pdf' -F 'name=Joy of Clojure'
