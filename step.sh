#!/bin/bash

set -e

# -----------------------
# --- Functions
# -----------------------

log_fail() { # red
  echo -e "\033[31m$1\033[0m"
  exit 1
}

log_warn() { # yellow
  echo -e "\033[33m$1\033[0m"
}

log_info() { # blue
echo
  echo -e "\033[34m$1\033[0m"
}

log_details() { # white
  echo -e "  \033[97m$1\033[0m"
}

log_done() { # green
  echo -e "  \033[32m$1\033[0m"
}

# -----------------------
# --- Main
# -----------------------

#
# Validate options
if [ -z "${name}" ] ; then
	log_fail "Missing required input: name"
fi

if [ -z "${platform}" ] ; then
	log_fail "Missing required input: platform"
fi

if [ -z "${abi}" ] ; then
	log_fail "Missing required input: abi"
fi

#
# Print options
log_info 'Configs:'
log_details "name: ${name}"
log_details "platform: ${platform}"
log_details "abi: ${abi}"

#
# Check if platform installed
log_info 'Check if platform installed'
platform_installed=true

platform_path="${ANDROID_HOME}/platforms/${platform}"
if [[ ! -d "${platform_path}" ]] ; then
  platform_installed=false
fi

#
# Install platform if needed
if [[ ${platform_installed} == true ]] ; then
  log_done "Platform ${platform} installed"
else
  log_details "Platform ${platform} not installed"

  log_info "Installing ${platform}"
  out=$(echo y | android update sdk --no-ui --all --filter ${platform})
  if [ $? -ne 0 ]; then
    echo "out: $out"
  fi
  log_done "Platform ${platform} installed"
fi

#
# Check if system image installed
log_info 'Check if system image installed'
system_image_installed=true

system_image_path="${ANDROID_HOME}/system-images/${platform}/${abi}"
if [ ! -d "${system_image_path}" ] ; then
	system_image_path="${ANDROID_HOME}/system-images/${platform}/default/${abi}"
	if [ ! -d "${system_image_path}" ] ; then
		system_image_installed=false
	fi
fi

#
# Install system image if needed
system_image="sys-img-${abi}-${platform}"

if [[ ${system_image_installed} == true ]] ; then
	log_done "System image ${system_image} installed"
else
	log_details "System image ${system_image} not installed"

	log_info "Installing ${system_image}"
	out=$(echo y | android update sdk --no-ui --all --filter ${system_image})
  if [ $? -ne 0 ]; then
    echo "out: $out"
  fi
	log_done "System image ${system_image} installed"
fi

#
# Create AVD image
log_info "Creating AVD image ${name}"
out=$(echo no | android create avd --force --name ${name} --target ${platform} --abi ${abi})
if [ $? -ne 0 ]; then
  echo "out: $out"
fi

envman add --key BITRISE_EMULATOR_NAME --value ${name}
log_done "AVD image ${name} ready to use 🚀"
