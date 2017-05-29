FROM jamesandersen/alpine-golang-opencv3:edge

WORKDIR /go/src/github.com/jamesandersen/gosudoku
COPY . .

RUN go-wrapper download   # "go get -d -v ./..."

RUN go-wrapper install    # "go install -v ./..."

# Set the PORT environment variable inside the container
ENV PORT 8080
# Expose port 8080 to the host so we can access our application
EXPOSE 8080
# Now tell Docker what command to run when the container starts
CMD ["go-wrapper", "run", "-filename=samples/800wi.png"]