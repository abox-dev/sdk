FROM docker.io/retailcrm/agentbox-base:v0.0.1@sha256:493b753044eb82f90b99afb3974033591ccbd7037b39e7963dfe49f38790376e

RUN printf '%s\n' 'agentbox-sdk-runtime-smoke' >/opt/agentbox-sdk-runtime-smoke
