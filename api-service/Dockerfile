FROM node:22-alpine3.19

# Install build dependencies (needed for ffmpeg compilation)

WORKDIR /api-service

COPY package*.json ./
RUN npm install 



COPY . .

EXPOSE 8000

CMD [ "npm", "start" ]