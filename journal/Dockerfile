FROM rust:latest as build
RUN rustup default nightly

WORKDIR /journal
COPY ./src ./Cargo.* ./

RUN cargo build --release

FROM scratch

COPY --from=build /journal/target/journal /
CMD ["/journal"]
