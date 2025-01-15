# Use the official Node.js image as a build stage
FROM node:18 AS build

# Set the working directory to /frontend
WORKDIR /frontend

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies and build the application
RUN npm install && npm run build

# ARG for backend API URL
ARG REACT_APP_API_URL
ENV REACT_APP_API_URL=${REACT_APP_API_URL}

# ARG for domain name
ARG REACT_APP_DOMAIN_NAME
ENV REACT_APP_DOMAIN_NAME=${REACT_APP_DOMAIN_NAME}

# Use a lightweight web server to serve the built app
FROM nginx:alpine

# Copy the built app from the previous stage
COPY --from=build /frontend/build /usr/share/nginx/html

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Set permissions for Nginx
RUN chown -R nginx:nginx /usr/share/nginx/html

# Switch to non-root user
USER nginx

# Expose the port the app runs on
EXPOSE 80

# Command to run the app
CMD ["nginx", "-g", "daemon off;"]