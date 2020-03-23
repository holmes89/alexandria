FROM rust:latest as build
RUN rustup default nightly

WORKDIR /journal
COPY ./journal/migrations ./migrations
COPY ./journal/src ./src
COPY ./journal/Cargo.* ./

ENV PORT 8000
EXPOSE 8000

RUN cargo build --release --bin journal
CMD ["./target/release/journal"]

FROM debian:stretch-slim as prod
RUN apt update && apt install libpq-dev -y
ENV PORT 8000
EXPOSE 8000

COPY --from=build /journal/target/release /release
CMD ["/release/journal"]
