#!/bin/bash

# Script to generate mock files using mockery

# Check if 'mockery' is installed, and install it if not
$(dirname "$0")/ensure_mockery.sh

mockery --log-level=ERROR