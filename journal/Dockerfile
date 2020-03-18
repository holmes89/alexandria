FROM rust:latest as build
RUN rustup default nightly

WORKDIR /journal
COPY ./journal/src ./src
COPY ./journal/Cargo.* ./

RUN cargo build --release

FROM scratch as prod

COPY --from=build /journal/target/release/journal /
CMD ["/journal"]
