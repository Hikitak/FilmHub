# simple node image for vite dev/build
FROM node:20-alpine
WORKDIR /app
COPY package.json package-lock.json* ./
RUN npm install --legacy-peer-deps --silent
COPY . .
EXPOSE 3000
CMD ["npm","run","dev","--","--host","0.0.0.0"] 