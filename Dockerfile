FROM golang:1.17.9

WORKDIR $GOPATH/src/hayvn
COPY . .


RUN make install

RUN echo $(ls -al)

RUN make build

RUN make ci

EXPOSE 8081

CMD ["hayvn"]
