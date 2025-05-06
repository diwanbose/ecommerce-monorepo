#!/bin/bash

# Create a temporary directory
TEMP_DIR=$(mktemp -d)
PROJECT_NAME="ecommerce-monorepo"

# Copy all necessary files
cp -r frontend backend k8s .github scripts Makefile docker-compose.yml README.md "$TEMP_DIR/$PROJECT_NAME"

# Remove unnecessary files
find "$TEMP_DIR/$PROJECT_NAME" -type d -name "node_modules" -exec rm -rf {} +
find "$TEMP_DIR/$PROJECT_NAME" -type d -name "build" -exec rm -rf {} +
find "$TEMP_DIR/$PROJECT_NAME" -type d -name "bin" -exec rm -rf {} +

# Create zip file
cd "$TEMP_DIR"
zip -r "$PROJECT_NAME.zip" "$PROJECT_NAME"

# Move zip file to project root
mv "$PROJECT_NAME.zip" "$OLDPWD"

# Clean up
rm -rf "$TEMP_DIR"

echo "Project has been packaged into $PROJECT_NAME.zip" 