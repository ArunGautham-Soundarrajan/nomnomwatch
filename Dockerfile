FROM public.ecr.aws/docker/library/golang:1.27.3-trixie AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
# timetzdata embeds the timezone database, so TZ works in the distroless image.
# templ and sqlc output is committed, so no generate step is needed.
RUN CGO_ENABLED=0 go build -tags timetzdata -o /out/nomnomwatch ./cmd/nomnomwatch

FROM gcr.io/distroless/static-debian12

COPY --from=build /out/nomnomwatch /usr/local/bin/nomnomwatch

USER 10001:10001
EXPOSE 8080

ENTRYPOINT ["/usr/local/bin/nomnomwatch"]
