FROM maven:3.8.3-openjdk-11

ENV WORK_DIR /flex-templates
RUN mkdir ${WORK_DIR}
WORKDIR ${WORK_DIR}

COPY v2/ ./

ENTRYPOINT [ "mvn", "spotless:check" ]