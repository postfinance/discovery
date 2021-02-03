 if [ $1 -eq 0 ]; then
    systemctl stop discovery.service
    systemctl disable discovery.service
fi
