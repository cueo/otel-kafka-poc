#!/usr/bin/env bash

# run below in background
nohup below record --retain-for-s 604800 --compress &

# run collector
/collector
