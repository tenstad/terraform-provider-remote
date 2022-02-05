FROM alpine:latest

COPY key.pub /root/.ssh/authorized_keys

RUN apk add --no-cache \
        bash \
        openssh \
        sudo \
    && ssh-keygen -A \
    && sed -i "s/#\?PermitRootLogin.*/PermitRootLogin yes/" /etc/ssh/sshd_config \
    && adduser -D bob \
    && echo "root:password" | chpasswd \
    && echo "bob:pwd" | chpasswd \
    && chmod 600 /root/.ssh/authorized_keys

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
