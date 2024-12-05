# Use a base image with Java
FROM openjdk:17-jdk-slim as builder

# Set the working directory
WORKDIR /app

# Copy the Maven or Gradle build file and the source code
COPY pom.xml ./
COPY src ./src

# Build the application
RUN ./mvnw clean package -DskipTests

# Use a smaller base image for the final image
FROM openjdk:17-jre-slim

# Set the working directory
WORKDIR /app

# Copy the jar file from the builder stage
COPY --from=builder /app/target/*.jar app.jar

# Expose the port the app runs on
EXPOSE 8080

# Command to run the application
ENTRYPOINT ["java", "-jar", "app.jar"]