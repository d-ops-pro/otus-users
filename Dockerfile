FROM alpine

ARG PACKAGE
ENV LISTEN=:8000
ENV DB_URI=""

COPY $PACKAGE /app

CMD ["sh", "-c", "/app --http-addr=$LISTEN --db-uri=$DB_URI"]
