# Use the official Node.js image as a build stage
FROM node:14 AS build

# Set the working directory
WORKDIR /frontend

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY . .

# ARG for backend API URL
ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=${REACT_APP_API_URL}

# ARG for domain name
ARG REACT_APP_DOMAIN_NAME
ENV REACT_APP_DOMAIN_NAME=${REACT_APP_DOMAIN_NAME}

# Build the application
RUN npm run build

# Use a lightweight web server to serve the built app
FROM nginx:alpine

# Copy the built app from the previous stage
COPY --from=build /frontend/build /usr/share/nginx/html

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose the port the app runs on
EXPOSE 80/tcp 

# Command to run the app
CMD ["/usr/sbin/nginx", "-g", "daemon off;"]  