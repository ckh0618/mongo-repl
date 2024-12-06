FROM mongo:7

COPY keyfile /etc/keyfile
RUN chown mongodb:mongodb /etc/keyfile && chmod 400 /etc/keyfile
