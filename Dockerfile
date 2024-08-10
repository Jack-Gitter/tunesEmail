# Copy all file and build
FROM golang as golang
WORKDIR /tunesEmail
COPY . . 
RUN go build

# Copy only binary and .env to other container
FROM golang
WORKDIR /tunesEmail
COPY --from=golang ./tunesEmail/tunesEmailService ./tunes
COPY --from=golang ./tunesEmail/.env ./.env

CMD ["./tunesEmailService"]

