FROM node:22-alpine3.19

# Install build dependencies (needed for ffmpeg compilation)

RUN apk add make gcc g++ python3

WORKDIR /back-end

COPY package*.json ./
RUN npm install 

# Download and install ffmp

COPY . ./

EXPOSE 8000
CMD [ "npm", "start" ]