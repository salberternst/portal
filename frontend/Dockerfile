FROM node:18

WORKDIR /app
COPY package-lock.json /app
COPY package.json /app
RUN npm ci
# COPY src /app/src
COPY index.html /app/index.html
COPY public /app/public
COPY tsconfig.json /app/tsconfig.json
COPY vite.config.ts /app/vite.config.ts

CMD npm run dev