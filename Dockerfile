FROM busybox
ADD messenger /
ADD migrations /migrations
EXPOSE 8080
ENTRYPOINT ["/messenger"]
