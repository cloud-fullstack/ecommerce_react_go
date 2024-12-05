# Use the official Node.js image as a build stage
FROM node:16 AS build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# Define build arguments
ARG VITE_API_URL
ARG API_URL
ARG DOMAIN_NAME

# Set environment variables for the build
ENV VITE_API_URL=$VITE_API_URL
ENV API_URL=$API_URL
ENV DOMAIN_NAME=$DOMAIN_NAME

# Build the application
RUN npm run build

# Use a lightweight web server to serve the built app
FROM nginx:alpine

# Copy the built app from the previous stage
COPY --from=build /app/build /usr/share/nginx/html

# Expose the port the app runs on
EXPOSE 80

# Command to run the app
CMD ["nginx", "-g", "daemon off;"]