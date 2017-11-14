
# golang:onbuild automatically copies the package source, 
# fetches the application dependencies, builds the program, 
# and configures it to run on startup 
# For more control, you can copy and build manually
FROM golang:latest 
LABEL Name=crud Version=0.0.1 
# RUN mkdir /app/src/crud
WORKDIR /app/src/crud
ADD . ${WORKDIR}    
ENV GOPATH=/app
RUN go build
EXPOSE 8080
RUN chmod +x crud
RUN ls -n
CMD ["./crud", "My personal musicbrainz feed", "0.0.1", "https://github.com/n0vice"]
