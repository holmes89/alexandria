FROM clojure:openjdk-11-lein AS build
COPY ./links /links
WORKDIR "/links"
ENV PORT 3000
EXPOSE 3000
RUN lein uberjar

FROM openjdk:11.0.6-jre-slim AS prod
ENV PORT 3000
EXPOSE 3000
COPY --from=build /links/target/uberjar/*-standalone.jar links-api.jar
CMD ["java", "-jar", "links-api.jar"]
