FROM scratch
COPY thunder /
ENTRYPOINT ["/thunder"]
