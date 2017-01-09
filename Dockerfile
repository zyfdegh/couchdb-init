FROM alpine
MAINTAINER zyfdegh <zyfdegg@gmail.com>

ENV COUCHDB_URL="http://127.0.0.1:5984/"
ENV COUCHDB_USER=""
ENV COUCHDB_PASS=""

# fix library dependencies
# otherwise golang binary may encounter 'not found' error
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

COPY couchdb-init /couchdb-init
CMD /couchdb-init
