FROM gesquive/go-builder:latest AS builder

ENV APP=gop

COPY dist/ /dist/
RUN copy-release

# =============================================================================
FROM gesquive/docker-base:latest
LABEL maintainer="Gus Esquivel <gesquive@gmail.com>"

COPY --from=builder /app/${APP} /app/

ENTRYPOINT ["/app/gop"]
