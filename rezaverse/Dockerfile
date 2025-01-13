# Use the official Node.js image as a build stage
FROM node:14 as build

# Set the working directory
WORKDIR /app

# Copy package.json and package-lock.json
COPY rezaverse/package*.json ./
COPY rezaverse/public ./public
COPY rezaverse/src ./src

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
COPY --from=build /app/build /usr/share/nginx/html

# Expose the port the app runs on
EXPOSE 80

# Command to run the app
CMD ["nginx", "-g", "daemon off;"]