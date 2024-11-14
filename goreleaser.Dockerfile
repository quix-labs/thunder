# Allow support for SSL + /tmp
FROM gcr.io/distroless/static
COPY thunder /
CMD ["/thunder"]