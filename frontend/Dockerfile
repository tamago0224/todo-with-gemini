# Use an official Node.js runtime as a parent image
FROM node:20-alpine

# Set the working directory in the container
WORKDIR /app

# Copy package.json and package-lock.json (or yarn.lock) to leverage Docker cache
COPY package*.json ./

# Install project dependencies
RUN npm install

# Copy the rest of the application's source code
COPY . .

# Expose the port the app runs on
EXPOSE 5173

# Command to run the app
CMD ["npm", "run", "dev", "--", "--host"]