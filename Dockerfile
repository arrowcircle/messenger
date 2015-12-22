FROM busybox
ADD messenger /
ADD migrations /migration
EXPOSE 8080
ENTRYPOINT ["/messenger"]
