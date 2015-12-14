FROM busybox
ADD messenger /
ADD migrations /
EXPOSE 8080
ENTRYPOINT ["/messenger"]
