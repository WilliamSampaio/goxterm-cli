#!/usr/bin/env bash

# Install the GoXterm systemd service
# Requires sudo privileges
SERVICE_FILE="goxterm.service"
DESTINATION="/etc/systemd/system/$SERVICE_FILE"

if [ ! -f "$SERVICE_FILE" ]; then
    echo "Service file $SERVICE_FILE not found!"
    exit 1
fi

echo "Copying $SERVICE_FILE to $DESTINATION..."
sudo cp "$SERVICE_FILE" "$DESTINATION"

echo "Reloading systemd daemon..."
sudo systemctl daemon-reload

echo "Enabling GoXterm service to start on boot..."
sudo systemctl enable goxterm.service

echo "Starting GoXterm service..."
sudo systemctl start goxterm.service

echo "GoXterm service installed and started successfully."
echo "You can check the service status with: sudo systemctl status goxterm.service"
echo "To view logs, use: sudo journalctl -u goxterm.service -f"
echo "To stop the service, use: sudo systemctl stop goxterm.service"
echo "To disable the service from starting on boot, use: sudo systemctl disable goxterm.service"
echo "Installation complete."
