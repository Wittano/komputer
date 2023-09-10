FROM gradle AS BUILDER
ENV APP_HOME=/usr/app/
WORKDIR $APP_HOME
COPY build.gradle settings.gradle $APP_HOME

COPY gradle $APP_HOME/gradle
COPY --chown=gradle:gradle . /home/gradle/src
USER root
RUN chown -R gradle /home/gradle/src

RUN gradle :cli:shadowJar || return 0
COPY . .

FROM openjdk17:alpine-jre
ENV ARTIFACT_NAME=komputer-cli-1.0.jar
ENV APP_HOME=/usr/app/

WORKDIR $APP_HOME
COPY --from=BUILDER $APP_HOME/build/libs/$ARTIFACT_NAME .

EXPOSE 8080
COPY .env $APP_HOME
RUN -java jar ${ARTIFACT_NAME} init
ENTRYPOINT exec java -jar ${ARTIFACT_NAME}