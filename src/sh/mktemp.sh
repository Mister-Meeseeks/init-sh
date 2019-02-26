
# Creates facility for user code to easily create temporary files
# and directories and have them cleaned up at the end of initSh
# run.

# Call with any arguments valid for system's mktemp utility. Besides
# tmpDir, which is always set to the INIT_SH runtime tmp directory.
# Only works on system's where mktemp supports --tmpdir flag.
function mkTempInitSh() {
    local args=$@
    local tmpDir=$INIT_SH_INSTANCE_DIR/tmpUser/
    mkdir -p $tmpDir
    mktemp --tmpdir=$tmpDir $args
}

export -f mkTempInitSh
