FROM scratch
COPY citrix-netscaler-exporter /app
EXPOSE 9280
ENTRYPOINT ["/app"]
