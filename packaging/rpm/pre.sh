 if [ $1 -eq 0 ]; then
    systemctl stop discovery.service
    systemctl stop exporter.service
    systemctl disable discovery.service
    systemctl disable exporter.service
fi
