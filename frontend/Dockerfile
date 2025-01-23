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

# Serve the build folder (if needed)
RUN npm install -g serve
CMD ["serve", "-s", "build", "-l", "5000"]

# Expose the port that serve will run on
EXPOSE 5000
