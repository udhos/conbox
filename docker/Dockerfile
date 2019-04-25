# STEP 1 build executable binary

#FROM golang:alpine as builder
FROM golang:1.12.4-alpine as builder
#RUN apk add --no-cache git
RUN apk add git

# Create appuser on builder image
RUN adduser -D -g '' conboxuser

COPY ./ /tmp/conbox/
WORKDIR /tmp/conbox

#build the binary
#RUN CGO_ENABLED=0 go build -o /bin/conbox
RUN CGO_ENABLED=0 go install -v ./conbox
RUN ls /go/bin/conbox

# STEP 2 build a small image

# start from scratch
FROM scratch

# copy appuser from builder image
COPY --from=builder /etc/passwd /etc/passwd

# Copy our static executable
COPY --from=builder /go/bin/conbox /bin/conbox

USER conboxuser

ENTRYPOINT ["/bin/conbox"]
