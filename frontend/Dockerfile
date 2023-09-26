FROM node:18

WORKDIR /app
COPY package-lock.json /app
COPY package.json /app
RUN npm install
COPY src /app/src
COPY public /app/public
COPY tsconfig.json /app/tsconfig.json
COPY vite.config.ts /app/vite.config.ts

CMD npm run dev