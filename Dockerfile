FROM centurylink/ca-certs

ADD ./dummy-portal /

ENTRYPOINT [ "/dummy-portal" ]
