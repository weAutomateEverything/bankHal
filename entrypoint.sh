#!/bin/sh
# Override user ID lookup to cope with being randomly assigned IDs using
# the -u option to 'docker run'.
USER_ID=$(id -u)
if [ x"$USER_ID" != x"0" -a x"$USER_ID" != x"1001" ]; then
    NSS_WRAPPER_PASSWD=/tmp/passwd.nss_wrapper
    NSS_WRAPPER_GROUP=/etc/group
    cat /etc/passwd | sed -e 's/^ipython:/builder:/' > $NSS_WRAPPER_PASSWD
    echo "ipython:x:$USER_ID:0:IPython,,,:/home/ipython:/bin/bash" >> $NSS_WRAPPER_PASSWD
    export NSS_WRAPPER_PASSWD
    export NSS_WRAPPER_GROUP
    LD_PRELOAD=/usr/local/lib64/libnss_wrapper.so
    export LD_PRELOAD
fi
/opt/dynatrace/oneagent/dynatrace-agent64.sh