#!/bin/bash
go build
sudo setcap 'cap_net_bind_service=+ep' "$(realpath program)"
source prod.env
./program