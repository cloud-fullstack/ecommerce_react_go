# Use the official Node.js image as a build stage
FROM node:14 AS builder

# Set the working directory to /app
WORKDIR /app

# Copy package.json and package-lock.json
COPY frontend/package*.json ./

# Install dependencies
RUN npm install

# Copy the rest of the application code
COPY frontend/ .

# Build the application
RUN npm run build

# Use the official Nginx image for the production stage
FROM nginx:stable

# Correct this line by referencing the named stage 'builder'
COPY --from=builder /app/build /usr/share/nginx/html

# Copy custom Nginx configuration
COPY nginx.conf /etc/nginx/nginx.conf

# Expose the port the app runs on
EXPOSE 80

# Command to run the app
CMD ["nginx", "-g", "daemon off;"]
