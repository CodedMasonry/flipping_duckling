#!/bin/bash

# Exit the script if any command fails
set -e

# Install garble for obfuscation
echo "Installing garble..."
go install mvdan.cc/garble@latest

# Set up environment variables for Windows compilation
echo "Setting up environment variables for Windows compilation..."
export GOOS=windows
export GOARCH=amd64

# Compile the Go code with garble
echo "Compiling Go code for Windows..."
garble build -o flipping_duckling.exe -tiny

# Compilation completed
echo "Compilation completed. The executable is flipping_duckling.exe"
