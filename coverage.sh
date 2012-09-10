#!/bin/sh

OUT=`pwd`/coverage

cd wbxml

rm -rf $OUT > /dev/null 2>&1
mkdir $OUT

DATA=$OUT/coverage.data

gocov test > $DATA 2>/dev/null


counter=0
for line in `gocov report $DATA | grep -v 100.00%`
do
  if [ $counter -eq 0 ]; then
    FILE=`echo $line | cut -f2 -d/`
  fi

  if [ $counter -eq 1 ]; then
    gocov annotate $DATA wbxml.$line >> $OUT/$FILE.cov
  fi

  counter=`expr $counter + 1`
  if [ $counter  -eq 4 ]; then
    counter=0
  fi
done

cd ..
