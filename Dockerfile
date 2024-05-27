FROM node:22.2.0-alpine as base
WORKDIR /usr/src/app
COPY package.json ./

FROM base as builder
RUN npm install
COPY index.ts ./
RUN npm run build

FROM base
RUN npm install --production
COPY --from=builder /usr/src/app/dist /usr/src/app/dist
CMD ["npm", "start"]
