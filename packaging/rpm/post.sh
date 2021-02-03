# systemd
systemctl daemon-reload
if [ $1 -ge 1 ]; then
  systemctl try-restart discovery.service
fi
