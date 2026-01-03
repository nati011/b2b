#!/bin/bash
set -e

echo "Stopping Docker services..."
sudo systemctl stop docker docker.socket containerd 2>/dev/null || true
sudo pkill -9 dockerd containerd docker-proxy 2>/dev/null || true

echo "Removing Docker packages..."
sudo apt-get remove -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras 2>/dev/null || true

echo "Purging Docker packages and cleaning up..."
sudo apt-get purge -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin docker-ce-rootless-extras 2>/dev/null || true
sudo apt-get autoremove -y 2>/dev/null || true

echo "Removing Docker data and configuration..."
sudo rm -rf /var/lib/docker
sudo rm -rf /var/lib/containerd
sudo rm -rf /etc/docker
sudo rm -rf ~/.docker

echo "Updating package list..."
sudo apt-get update

echo "Installing prerequisites..."
sudo apt-get install -y ca-certificates curl gnupg lsb-release

echo "Adding Docker's official GPG key..."
sudo install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /etc/apt/keyrings/docker.gpg
sudo chmod a+r /etc/apt/keyrings/docker.gpg

echo "Setting up Docker repository..."
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.gpg] https://download.docker.com/linux/ubuntu \
  $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

echo "Installing Docker..."
sudo apt-get update
sudo apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

echo "Starting Docker service..."
sudo systemctl start docker
sudo systemctl enable docker

echo "Adding user to docker group (you may need to log out and back in)..."
sudo usermod -aG docker $USER

echo "Verifying Docker installation..."
sudo docker --version
sudo docker ps -a

echo ""
echo "Docker has been reinstalled successfully!"
echo "Note: You may need to log out and log back in for the docker group changes to take effect."
echo "After that, you should be able to run docker commands without sudo."

