#!/bin/bash

TRXSET=$1
STATSCODE=$TRXSET


STATSFILE="$HOME/prebakedstats.txt"
NET="prebaked"
GIVDIR="$HOME/.dogye/networks/$NET"
TRX="$GIVDIR/trx"
PREFIX="Test"
SUFFIX=".json"
FAUCET="Faucet"
NODEHOST="172.77.5.10"
NODEPORT="8080"
NODENAME="Node0"
OUTDIR=$TRX/$TRXSET

STEM=$(basename $TRXSET)
STEM=${STEM##TRX_}
ACCTCNT=${STEM%%_*}
TRXCNT=${STEM##*_}


TRXFILE=""
CONFIGDIR=""
PRE=""

# Store current path
mydir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" > /dev/null && pwd )"


if [ ! -d "$OUTDIR" ] ; then
    >&2 echo "Cannot find $OUTDIR"
    exit 1
fi

# Start TRX Timestamp
res2=$(date +%s.%N)


# Launch signed transactions processing as a background process
PIDS=""
for i in $(seq 1 $ACCTCNT)
do
    ( $mydir/run-trx.sh $OUTDIR/$PREFIX$i$SUFFIX  ) & PIDS="$PIDS $!"
done

# Wait for background tasks to finish
FAIL=0
for job in $PIDS
do
    wait $job || let "FAIL+=1"
done

echo ""

# Timings
# Finish timer
res3=$(date +%s.%N)
dt=$(echo "$res2 - $res1" | bc)
dt2=$(echo "$res3 - $res2" | bc)


# Check values of accounts as expected
echo node $mydir/index.js --account=$FAUCET --nodename=$NODENAME --nodehost=$NODEHOST \
 --nodeport=$NODEPORT --TRXfile=$TRXFILE --configdir=$CONFIGDIR  --pre=$PRE
exitcode=$?

echo "Preparing $TRXCNT transactions took $dt seconds"
echo "$TRXCNT transactions applying took $dt2 seconds"
rate=$(echo "scale=4;$TRXCNT / $dt2" | bc)
echo "$rate transactions per second"


if [ $exitcode -ne 0 ] ; then
    echo "Balance checks failed."
    exit $exitcode
fi


if [ "$FAIL" == "0" ];
then
    echo "PASSED"

    if [ ! -z "$STATSCODE" ] ; then
      echo "$TRXCNT $ACCTCNT $dt2 $STATSCODE" >> $STATSFILE
    fi


else
    echo "FAIL! ($FAIL)"
    exit 5
fi