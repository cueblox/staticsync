FROM ubuntu:20.04

COPY staticsync /

ENTRYPOINT ["/staticsync"]
