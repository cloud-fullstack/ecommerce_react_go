# Use an official Node.js runtime as a parent image
FROM node:16

# Set the working directory in the container
WORKDIR /app

# Copy the package.json and package-lock.json (if they exist)
COPY frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy the entire frontend directory into the container
COPY frontend/ .

# Set the PUBLIC_URL build argument
ARG PUBLIC_URL
ENV PUBLIC_URL=$PUBLIC_URL

# Build the React app
RUN npm run build

# Install Nginx
RUN apt-get update && apt-get install -y nginx

# Copy the build files to the Nginx directory
RUN mkdir -p /var/www/html && cp -r build/* /var/www/html/

# Copy the Nginx configuration file
COPY ../nginx.conf /etc/nginx/sites-available/default

# Expose port 80
EXPOSE 80

# Start Nginx
CMD ["nginx", "-g", "daemon off;"]