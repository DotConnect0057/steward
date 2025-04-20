FROM golang:latest

   # Install necessary packages
   RUN git clone https://github.com/DotConnect0057/steward.git

   RUN cd steward && \
       go mod download && \
       go build -o st && \
       chmod +x st

   CMD ["sleep", "infinity"]
   # Expose the port that the application will run on