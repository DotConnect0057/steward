FROM alpine/git:latest

   # Install necessary packages
   RUN apk update && apk add --no-cache openssh

   RUN mkdir -p /var/run/sshd

   # Set root password (for testing purposes only - NOT SECURE)
   RUN echo 'root:password' | chpasswd

   # Add a non-root user
   RUN useradd myuser -m -s /bin/bash
   RUN echo 'myuser:password' | chpasswd

   # Grant sudo privileges to the user
   RUN usermod -aG sudo myuser

   # SSH configuration (adjust as needed)
   RUN sed -i 's/PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
   RUN sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/' /etc/ssh/sshd_config

   # SSH host key generation
   RUN ssh-keygen -A

   # Expose SSH port
   EXPOSE 22

   # Start SSH service
   CMD ["/usr/sbin/sshd", "-D"]