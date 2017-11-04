FROM ubuntu:16.04
RUN apt-get update
RUN mkdir /app
COPY graphql-golang-hw /app
CMD ./app/graphql-golang-hw