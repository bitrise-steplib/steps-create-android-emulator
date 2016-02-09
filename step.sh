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

if [ -z "${target}" ] ; then
	log_fail "Missing required input: target"
fi

if [ -z "${abi}" ] ; then
	log_fail "Missing required input: abi"
fi

#
# Print options
log_info 'Configs:'
log_details "name: ${name}"
log_details "target: ${target}"
log_details "abi: ${abi}"

log_info 'Check if target installed'
target_installed=true

#
# Check if target exist
system_image_path="$ANDROID_HOME/system-images/${target}/${abi}"
if [ ! -d "${system_image_path}" ] ; then
	system_image_path="$ANDROID_HOME/system-images/${target}/default/${abi}"
	if [ ! -d "${system_image_path}" ] ; then
		target_installed=false
	fi
fi

#
# Install target
if [[ $target_installed == true ]] ; then
	log_done "Target ${target} and abi ${abi} installed"
else
	log_details "Target ${target} and abi ${abi} not installed"

	log_info "Installing ${target} and abi ${abi}"
	sys_image="sys-img-${abi}-${target}"
	echo y | android update sdk --no-ui --all --filter ${sys_image}
	log_done "Target ${target} and abi ${abi} installed"
fi

log_info "Creating AVD image ${name}"
echo no | android --silent create avd --force --name $name --target ${target} --abi ${abi}

envman add --key BITRISE_EMULATOR_NAME --value $name
log_done "AVD image $name ready to use 🚀"
