#!/usr/bin/env bash
# Copyright © 2020. Financelime, https://financelime.com. All rights reserved.
# Author: DmAlix. Contacts: <dmalix@financelime.com>, <dmalix@yahoo.com>
# License: GNU General Public License v3.0, https://www.gnu.org/licenses/gpl-3.0.html

# Config
readonly currentTime=$(date -u '+%Y%m%d_%H%M%S')
readonly host=${TEST_BENCH_UPDATE_HOST}
readonly port=${TEST_BENCH_UPDATE_PORT}
readonly user=${TEST_BENCH_UPDATE_USER}
readonly unit=${TEST_BENCH_UPDATE_UNIT}
readonly binName=financelime-rest-api
readonly remoteServiceHomePath=${TEST_BENCH_UPDATE_REMOTE_SERVICE_HOME_PATH}

readonly localSystemdHomePath=${TEST_BENCH_UPDATE_LOCAL_SYSTEMD_HOME_PATH}
readonly localSystemdFileName=${TEST_BENCH_UPDATE_LOCAL_SYSTEMD_FILENAME}
readonly remoteSystemdHomePath=${TEST_BENCH_UPDATE_REMOTE_SYSTEMD_HOME_PATH}
readonly remoteSystemdFileName=${TEST_BENCH_UPDATE_REMOTE_SYSTEMD_FILENAME}

# Confirm run
read -n 1 -p "Run test-bench update (y/[a])? " userInput
if [ "${userInput}" != "y" ] ; then echo ""; echo -e "\e[1;31mRun canceled\e[0m"; exit; fi
echo ""

echo -ne "- Stop service on test-bench:\t\t\t\t\t"
ssh -p ${port} ${user}@${host} "systemctl stop ${unit}"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to stop service on test-bench [RwKmoB3Y]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Compress the new build file:\t\t\t\t\t"
cd bin; gzip --keep --force ${binName}; gzip --test ${binName}.gz; cd ..
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to compress the new build file [rOeKqt1e]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Backup the prev build file on test-bench:\t\t\t"
ssh -p ${port} ${user}@${host} \
"cd ${remoteServiceHomePath}/bin; gzip ${binName}; gzip --test ${binName}.gz; mv ${binName}.gz ${binName}.${currentTime}.gz"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to backup the prev build file on test-bench [X0QuTyLY]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo "- Copy the archive with new build file to test-bench:"
scp -P ${port}	bin/${binName}.gz	${user}@${host}:${remoteServiceHomePath}/bin
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to copy the archive with new build file to test-bench [2y6E2Cat]\e[0m"; exit; fi

echo -ne "- Extract the new build file on test-bench:\t\t\t"
ssh -p ${port} ${user}@${host} "cd ${remoteServiceHomePath}/bin; gzip --decompress ${binName}.gz;"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to extract the new build file on test-bench [vQpM6nld]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Compress the new migrate files:\t\t\t\t"
tar --create --gzip --file=migrate.tar.gz migrate
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to compress the migrate files [kHB4jlqD]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Backup the prev migrate files on test-bench:\t\t\t"
ssh -p ${port} ${user}@${host} \
"cd ${remoteServiceHomePath}; tar --create --gzip --file=migrate.${currentTime}.tar.gz migrate; rm --force --dir --recursive migrate"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to backup the prev migrate files on test-bench [s65iZcf9]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo "- Copy the archive with migrate files to test-bench:"
scp -P ${port}	migrate.tar.gz	${user}@${host}:${remoteServiceHomePath}
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to Copy the archive with migrate files to test-bench [dihL5qGM]\e[0m"; exit; fi

echo -ne "- Extract the new migrate files on test-bench:\t\t\t"
ssh -p ${port} ${user}@${host} \
"cd ${remoteServiceHomePath}; tar --extract --gzip --file=migrate.tar.gz; rm migrate.tar.gz"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to extract the new migrate files on test-bench [s65iZcf9]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Remove the archive with with migrate files to test-bench:\t"
rm migrate.tar.gz
if [ $? -ne 0 ] ; then echo -e "\e[1;31mRemove the archive with with migrate files to test-bench [E622yCat]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Backup the prev systemd file on test-bench:\t\t\t"
ssh -p ${port} ${user}@${host} \
"cd ${remoteSystemdHomePath}; tar --create --gzip --file=${remoteSystemdFileName}.${currentTime}.tar.gz ${remoteSystemdFileName}"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to backup the prev systemd files on test-bench [Zs65icf9]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo "- Copy the new systemd file to test-bench:"
scp -P ${port}	${localSystemdHomePath}/${remoteSystemdFileName}	${user}@${host}:${remoteSystemdHomePath}/${remoteSystemdFileName}
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to copy the new systemd file to test-bench [22y6ECat]\e[0m"; exit; fi

echo -ne "- Reload Systemd daemon:\t\t\t\t\t"
ssh -p ${port} ${user}@${host} "systemctl daemon-reload"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to reload the Systemd daemon [RmB3owKY]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Start service on test-bench:\t\t\t\t\t"
ssh -p ${port} ${user}@${host} "systemctl start ${unit}"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to start service on test-bench [oRwKmB3Y]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo "- Status service on test-bench:"
ssh -p ${port} ${user}@${host} "systemctl status ${unit} --full --lines=30"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to check status of the service on test-bench [oRwKmB3Y]\e[0m"; exit; fi

echo -e "\e[32mSuccessful completion\e[0m"
echo "--"
