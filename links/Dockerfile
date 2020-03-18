FROM clojure:openjdk-11-lein AS build
COPY ./links /links
WORKDIR "/links"
RUN lein uberjar

FROM openjdk:alpine

COPY --from=build /links/target/uberjar/*-standalone.jar links-api.jar
CMD ["java", "-jar", "links-api.jar"]
