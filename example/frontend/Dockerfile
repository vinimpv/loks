# Stage 1: Build the React application
FROM node:16 AS build

WORKDIR /app
COPY ./package.json ./package-lock.json ./
RUN npm install
COPY ./ ./
RUN npm run build

# Stage 2: Serve the React application
FROM nginx:alpine

COPY --from=build /app/build /usr/share/nginx/html

# Copy the default nginx.conf provided by the image (optional)
# We can customize this later if needed.
COPY ./nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]
