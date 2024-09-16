#FROM alpine
#COPY thunder /bin
#CMD ["thunder"]

# Do not work without sslmode=disable
FROM scratch
COPY thunder /
CMD ["/thunder"]