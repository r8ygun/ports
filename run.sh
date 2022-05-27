#!/usr/bin/env bash

sudo docker compose down
sudo docker compose build
sudo docker compose up --force-recreate --remove-orphans --abort-on-container-exit
