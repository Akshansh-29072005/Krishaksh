#!/bin/bash
# EC2 Setup Script for Krisho Backend
# Run this after connecting to your EC2 instance

set -e

echo "🚀 Setting up Krisho Backend on EC2..."

# Update system
echo "📦 Updating system packages..."
sudo apt update && sudo apt upgrade -y

# Install Docker
echo "🐳 Installing Docker..."
sudo apt install -y docker.io docker-compose-plugin
sudo systemctl enable docker
sudo systemctl start docker
sudo usermod -aG docker $USER

echo "✅ Docker installed. Please logout and login again, then run the next script."

# Note: User needs to logout/login for docker group to take effect
echo "🔄 Please run: logout, then ssh back in, then run ./setup-repo.sh"