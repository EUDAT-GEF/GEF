#!/bin/bash

# Shell script for use within the ENES CDO example for the EUDAT Generic Execution Framework (GEF)
# by Asela Rajapakse (asela.rajapakse@mpimet.mpg.de).
# Adaptation of a shell script prepared for Sakradzija, M. and C. Hohenegger, 2017:
# What determines the distribution of shallow convective mass flux through cloud base?
# J. Atmos. Sci., https://doi.org/10.1175/JAS-D-16-0326.1.
# (http://journals.ametsoc.org/doi/10.1175/JAS-D-16-0326.1)
# by the authors.

# Uses the CDO select operator to select a given climate variable (lwp in this case)
# from a set of NetCDF files into intermediate NetCDF files
# and then uses the CDO gather operator to merge all the files into one.

function gathervar {
echo $@
fstem=$1
nx=$2
ny=$3
var=$4
nx=1
  if [ -s /output/gathered_$fstem.out.xy.$var.15ts.nc ]; then
    echo " File gathered_$fstem.out.xy.$var.15ts.nc exists already, skipping this variable"
  else
    echo " Do cdo gather for var=$var in file=$fstem"

    for n in $(seq 0 $ny); do
      nstring=$(printf %04d $n)
      if [ -s /output/gather.xy.$nstring.$var.15ts.nc ]; then
        echo " File gather.xy.$nstring.$var.15ts.nc exists already, using this file"
      else
        echo " Do cdo gather for n=$nstring for all var=$var"
        for m in $(seq 0 $nx); do
          mstring=$(printf %04d $m)
          cdo selname,$var /output/$fstem.out.xy.$var.${nstring}.${mstring}.15ts.nc /output/selname.xy.$var.$nstring.$mstring.15ts.nc
        done
        cdo gather /output/selname.xy.$var.$nstring.????.15ts.nc /output/gather.xy.$var.$nstring.15ts.nc
        rm /output/selname.xy.$var.$nstring.????.15ts.nc
      fi
    done
    echo "cdo gather for gathered slices of var=$var"
    cdo gather /output/gather.xy.$var.????.15ts.nc /output/gathered_$fstem.out.xy.$var.15ts.nc

    rm /output/gather.xy.$var.????.15ts.nc
  fi

# rico_gcss.out.xy.lwp.0052.0014.15ts.nc

}

export -f gathervar

# The cdo_gather_lwp_sh script has the following configuration variables:
# 1) The number of processes to be used for computation
# 2) The name prefix of the data files to be used (rico_gcss in the ENES CDO example)

nproc=2
fstem="rico_gcss"

# The dataset for this example is downloaded by the GEF and put into the /input directory and the archive unpacked into /output.

cd /input

tar -xvf rico_gcss_out_xy_lwp_15ts.tar -C /output

cd /output

# Determines the number of NetCDF files and calls the gathervar function for the climate variables specified (only lwp in this case).

nn=$( ls -l /output/${fstem}.out.xy.lwp.????.????.15ts.nc | wc -l )
nx=$( ls -l /output/${fstem}.out.xy.lwp.0000.????.15ts.nc | wc -l )
ny=$( ls -l /output/${fstem}.out.xy.lwp.????.0000.15ts.nc | wc -l )

echo " cdo_gather_cross for $fstem dataset"
echo " Looking for files in dir=$dir"
echo " Found nn=$nn files from nx=$nx ny=$ny proc"

nx=$(( $nx - 1 ))
ny=$(( $ny - 1 ))

varnames="lwp"

echo " Calling gathervar with var=$varnames"

echo -e "$varnames" | xargs -n1 -P$nproc -I{} -d\  bash -c "gathervar $fstem $nx $ny {}"

rm /output/rico_gcss*
