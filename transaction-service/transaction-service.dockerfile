FROM alpine:latest

RUN mkdir app

COPY transactionApp /app

CMD ["app/transactionApp"]