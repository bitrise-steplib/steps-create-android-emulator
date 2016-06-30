#!/bin/bash

set -e

# -----------------------
# --- Functions
# -----------------------

RESTORE='\033[0m'
RED='\033[00;31m'
YELLOW='\033[00;33m'
BLUE='\033[00;34m'
GREEN='\033[00;32m'

function color_echo {
	color=$1
	msg=$2
	echo -e "${color}${msg}${RESTORE}"
}

function echo_fail {
	msg=$1
	echo
	color_echo "${RED}" "${msg}"
	exit 1
}

function echo_warn {
	msg=$1
	color_echo "${YELLOW}" "${msg}"
}

function echo_info {
	msg=$1
	echo
	color_echo "${BLUE}" "${msg}"
}

function echo_details {
	msg=$1
	echo "  ${msg}"
}

function echo_done {
	msg=$1
	color_echo "${GREEN}" "  ${msg}"
}

function validate_required_input {
	key=$1
	value=$2
	if [ -z "${value}" ] ; then
		echo_fail "[!] Missing required input: ${key}"
	fi
}

function print_and_run {
  cmd="$1"
  echo_details "${cmd}"
	echo
  eval "${cmd}"
}

# -----------------------
# --- Main
# -----------------------

#
# Validate options
validate_required_input "name", "${name}"

if [ -z "${name}" ] ; then
	echo_fail "Missing required input: name"
fi

if [ -z "${platform}" ] ; then
	echo_fail "Missing required input: platform"
fi

if [ -z "${abi}" ] ; then
	echo_fail "Missing required input: abi"
fi

#
# Print options
echo_info 'Configs:'
echo_details "name: ${name}"
echo_details "platform: ${platform}"
echo_details "abi: ${abi}"

#
# Check if platform installed
echo_info 'Check if platform installed'
platform_installed=true

platform_path="${ANDROID_HOME}/platforms/${platform}"
echo_details "checking path: ${platform_path}"
if [[ ! -d "${platform_path}" ]] ; then
  platform_installed=false
fi

#
# Install platform if needed
if [[ ${platform_installed} == true ]] ; then
  echo_done "platform ${platform} installed"
else
  echo_details "platform ${platform} not installed"
  echo_details "installing ${platform}"

	install_cmd="echo y | android update sdk --no-ui --all --filter ${platform}"
	print_and_run "${install_cmd}"
  if [ $? -ne 0 ] ; then
    echo_fail "command failed"
  fi

  echo_done "platform ${platform} installed"
fi

#
# Check if system image installed
echo_info 'Check if system image installed'
system_image_installed=true

system_image_path="${ANDROID_HOME}/system-images/${platform}/${abi}"
echo_details "checking path: ${system_image_path}"

if [ ! -d "${system_image_path}" ] ; then
	system_image_path="${ANDROID_HOME}/system-images/${platform}/default/${abi}"
	echo_details "checking path: ${system_image_path}"

	if [ ! -d "${system_image_path}" ] ; then
		system_image_installed=false
	fi
fi

#
# Install system image if needed
system_image="sys-img-${abi}-${platform}"

if [[ ${system_image_installed} == true ]] ; then
	echo_done "system image ${system_image} installed"
else
	echo_details "system image ${system_image} not installed"

	echo_info "Installing ${system_image}"

	install_cmd="echo y | android update sdk --no-ui --all --filter ${system_image}"
	print_and_run "${install_cmd}"

	# Check if install succed
	system_image_path="${ANDROID_HOME}/system-images/${platform}/${abi}"
	echo_details "checking path: ${system_image_path}"

	if [ ! -d "${system_image_path}" ] ; then
		system_image_path="${ANDROID_HOME}/system-images/${platform}/default/${abi}"
		echo_details "checking path: ${system_image_path}"

		if [ ! -d "${system_image_path}" ] ; then
			echo_fail "system image ${system_image} not installed"
		fi
	fi
	##

	echo_done "system image ${system_image} installed"
fi

#
# Create AVD image
echo_info "Creating AVD image ${name}"
create_cmd="echo no | android create avd --force --name ${name} --target ${platform} --abi ${abi}"
print_and_run "${create_cmd}"
if [ $? -ne 0 ] ; then
  echo_fail "command failed"
fi

envman add --key BITRISE_EMULATOR_NAME --value ${name}
echo_done "AVD image ${name} ready to use 🚀"
