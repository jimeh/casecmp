FROM scratch
COPY ./casecmp /
ENV PORT 8080
EXPOSE 8080
WORKDIR /
ENTRYPOINT ["/casecmp"]
