FROM scratch
COPY goseal /bin/goseal

EXPOSE 18212
ENTRYPOINT ["/bin/goseal"]