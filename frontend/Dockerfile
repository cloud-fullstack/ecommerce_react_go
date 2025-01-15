# Use the official Node.js image as a build stage
FROM node:14 AS build

# Set the working directory to /frontend
WORKDIR /frontend

# Copy package.json and package-lock.json
COPY package*.json ./

# Copy the root configuration files
COPY postcss.config.js ./
COPY tailwind.config.js ./
COPY tsconfig.json ./

# Copy the public directory
COPY public ./public

# Copy the src directory (including App.tsx, index.tsx, etc.)
COPY src ./src

# Install dependencies
RUN npm install

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

# Expose the port the app runs on
EXPOSE 80

# Command to run the app
CMD ["nginx", "-g", "daemon off;"]