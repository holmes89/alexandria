FROM rust:latest as build
RUN rustup default nightly

WORKDIR /journal
COPY ./journal/src ./src
COPY ./journal/Cargo.* ./

ENV PORT 8000
EXPOSE 8000

RUN cargo build --release

FROM scratch as prod

ENV PORT 8000
EXPOSE 8000
COPY --from=build /journal/target/release/journal /
CMD ["/journal"]
