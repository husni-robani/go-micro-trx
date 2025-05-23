FROM alpine:latest

# make directory
RUN mkdir app

# copy binary file to /app
COPY brokerApp /app

# run binaryfile in container
CMD ["app/brokerApp"]