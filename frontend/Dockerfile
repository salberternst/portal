FROM node:18

WORKDIR /app
COPY package-lock.json /app
COPY package.json /app
RUN npm install
COPY docs /app/docs
COPY src /app/src
COPY static /app/static
COPY babel.config.js /app/babel.config.js
COPY docusaurus.config.js /app/docusaurus.config.js
COPY sidebars.js /app/sidebars.js

CMD npm run start