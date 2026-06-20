#!/bin/bash
# Repository Setup and Deployment Script
# Run this after setup-ec2.sh and logout/login

set -e

echo "📁 Setting up Krisho repository..."

# Clone repository (replace with your actual repo URL)
# git clone https://github.com/your-username/krishaksh.git
# cd krishaksh/server

# Or if you have the code locally, upload it via scp
echo "💡 If you haven't cloned yet, run:"
echo "git clone https://github.com/your-username/krishaksh.git"
echo "cd krishaksh/server"

# Create directories
echo "📁 Creating directories..."
mkdir -p deploy/nginx/certs
mkdir -p deploy/secrets

echo "📋 Next steps:"
echo "1. Place SSL certificates in deploy/nginx/certs/"
echo "2. Place Firebase credentials in deploy/secrets/"
echo "3. Create .env.prod file"
echo "4. Run: ./scripts/deploy-prod.sh .env.prod"

echo "✅ Repository setup complete!"