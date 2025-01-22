# frontend/Dockerfile

# Use an official Node.js runtime as a parent image
FROM node:16

# Set the working directory in the container
WORKDIR /app

# Copy the package.json and package-lock.json (if they exist)
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy the entire frontend directory into the container
COPY . .

# Build the React app
RUN npm run build

# No CMD or EXPOSE here since GitLab Pages will handle hosting the static files.
