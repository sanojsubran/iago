FROM  scratch
RUN apk update && apk add bash
COPY  iago /
CMD   ["/iago"]