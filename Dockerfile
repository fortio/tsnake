FROM scratch
COPY tsnake /usr/bin/tsnake
ENV HOME=/home/user
ENTRYPOINT ["/usr/bin/tsnake"]
