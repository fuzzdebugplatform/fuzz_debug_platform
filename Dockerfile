FROM ubuntu:18.04

COPY ./web/build /fdplt/web
COPY ./bin/platform /bin/platform

ENTRYPOINT ["/bin/platform", "-W", "/fdplt/web"]