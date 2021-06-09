#!/usr/bin/env bash

# Config
readonly currentTime=$(date -u '+%Y%m%d_%H%M%S')
readonly host=${STAGING_NODE_UPDATE_HOST}
readonly port=${STAGING_NODE_UPDATE_PORT}
readonly user=${STAGING_NODE_UPDATE_USER}
readonly unit=${STAGING_NODE_UPDATE_UNIT}
readonly binName=financelime-auth
readonly remoteServiceHomePath=${STAGING_NODE_UPDATE_REMOTE_SERVICE_HOME_PATH}

readonly localSystemdHomePath=${STAGING_NODE_UPDATE_LOCAL_SYSTEMD_HOME_PATH}
readonly remoteSystemdHomePath=${STAGING_NODE_UPDATE_REMOTE_SYSTEMD_HOME_PATH}
readonly remoteSystemdFileName=${STAGING_NODE_UPDATE_REMOTE_SYSTEMD_FILENAME}

# Confirm run
read -n 1 -p "Run staging-node update (y/[a])? " userInput
if [ "${userInput}" != "y" ] ; then echo ""; echo -e "\e[1;31mRun canceled\e[0m"; exit; fi
echo ""

echo -ne "- Stop service on staging-node:\t\t\t\t\t"
ssh -p "${port}" "${user}"@"${host}" "systemctl stop ${unit}"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to stop service on staging-node [RwKmoB3Y]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Compress the new build file:\t\t\t\t\t"
cd bin || exit; gzip --keep --force ${binName}; gzip --test ${binName}.gz; cd ..
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to compress the new build file [rOeKqt1e]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Backup the prev build file on staging-node:\t\t\t"
ssh -p "${port}" "${user}"@${host} \
"cd ${remoteServiceHomePath}/bin; gzip ${binName}; gzip --test ${binName}.gz; mv ${binName}.gz ${binName}.${currentTime}.gz"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to backup the prev build file on staging-node [X0QuTyLY]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo "- Copy the archive with new build file to staging-node:"
scp -P ${port}	bin/${binName}.gz	${user}@${host}:${remoteServiceHomePath}/bin
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to copy the archive with new build file to staging-node [2y6E2Cat]\e[0m"; exit; fi

echo -ne "- Extract the new build file on staging-node:\t\t\t"
ssh -p ${port} ${user}@${host} "cd ${remoteServiceHomePath}/bin; gzip --decompress ${binName}.gz;"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to extract the new build file on staging-node [vQpM6nld]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

########################################################################################################################

echo -ne "- Compress the new assets files:\t\t\t\t"
tar --create --gzip --file=assets.tar.gz assets
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to compress the assets files [kHB4jlqD]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Backup the prev assets files on staging-node:\t\t\t"
ssh -p ${port} ${user}@${host} \
"cd ${remoteServiceHomePath}; tar --create --gzip --file=assets.${currentTime}.tar.gz assets; rm --force --dir --recursive assets"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to backup the prev assets files on staging-node [s65iZcf9]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo "- Copy the archive with assets files to staging-node:"
scp -P ${port}	assets.tar.gz	${user}@${host}:${remoteServiceHomePath}
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to Copy the archive with assets files to staging-node [dihL5qGM]\e[0m"; exit; fi

echo -ne "- Extract the new assets files on staging-node:\t\t\t"
ssh -p ${port} ${user}@${host} \
"cd ${remoteServiceHomePath}; tar --extract --gzip --file=assets.tar.gz; rm assets.tar.gz"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to extract the new assets files on staging-node [s65iZcf9]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Remove the archive with with assets files to staging-node:\t"
rm assets.tar.gz
if [ $? -ne 0 ] ; then echo -e "\e[1;31mRemove the archive with with assets files to staging-node [E622yCat]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

########################################################################################################################

echo -ne "- Backup the prev systemd file on staging-node:\t\t\t"
ssh -p ${port} ${user}@${host} \
"cd ${remoteSystemdHomePath}; tar --create --gzip --file=${remoteSystemdFileName}.${currentTime}.tar.gz ${remoteSystemdFileName}"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to backup the prev systemd files on staging-node [Zs65icf9]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo "- Copy the new systemd file to staging-node:"
scp -P ${port}	${localSystemdHomePath}/${remoteSystemdFileName}	${user}@${host}:${remoteSystemdHomePath}/${remoteSystemdFileName}
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to copy the new systemd file to staging-node [22y6ECat]\e[0m"; exit; fi

echo -ne "- Reload Systemd daemon:\t\t\t\t\t"
ssh -p ${port} ${user}@${host} "systemctl daemon-reload"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to reload the Systemd daemon [RmB3owKY]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo -ne "- Start service on staging-node:\t\t\t\t"
ssh -p ${port} ${user}@${host} "systemctl start ${unit}"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to start service on staging-node [oRwKmB3Y]\e[0m"; exit; fi
echo -e "\e[32mOK\e[0m"

echo "- Status service on staging-node:"
ssh -p ${port} ${user}@${host} "systemctl status ${unit} --full --lines=30"
if [ $? -ne 0 ] ; then echo -e "\e[1;31mFailed to check status of the service on staging-node [oRwKmB3Y]\e[0m"; exit; fi

echo -e "\e[32mSuccessful completion\e[0m"
echo "--"
