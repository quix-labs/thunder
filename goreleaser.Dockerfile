#FROM alpine
#COPY thunder /bin
#CMD ["thunder"]
FROM scratch
COPY thunder /
CMD ["/thunder"]