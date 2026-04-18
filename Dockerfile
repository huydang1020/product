FROM alpine:3.14

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root

# Copy binary đã build sẵn và assets
COPY product .

EXPOSE 8000 8001

CMD ["./product", "start"]