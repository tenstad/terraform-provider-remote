FROM alpine:latest

RUN apk add --no-cache \
        bash \
        openssh \
        sudo \
    && ssh-keygen -A \
    && sed -i "s/#\?PermitRootLogin.*/PermitRootLogin yes/" /etc/ssh/sshd_config \
    && echo "root:password" | chpasswd

COPY key.pub /root/.ssh/authorized_keys

EXPOSE 22
CMD ["/usr/sbin/sshd", "-D"]
