#!/bin/sh

CONFIG_DIR=/usr/local/share/etc/waveform
CONFIG_FILE_NAME=config.yaml
CONFIG_FILE_PATH=${CONFIG_DIR}/${CONFIG_FILE_NAME}

GITHUB_RAW_URL="https://raw.githubusercontent.com/woshikedayaa/waveform-backend/main/"

# 安装配置文件
if [ -f ${CONFIG_FILE_PATH} ];then
  mkdir -p ${CONFIG_DIR}
  echo Download:${GITHUB_RAW_URL}config/config_full.yaml ">" ${CONFIG_FILE_PATH}
  curl -SLo $CONFIG_FILE_PATH
fi

# todo 自动安装最新的 release