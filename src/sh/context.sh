
# User-facing library that allows for introspection and management of the
# INIT_SH calling context.

# This function resets any INIT_SH specific context variables, making it
# appear to the environment that we're no longer in an INIT_SH context.
# (Though any previous imports will still be accessible from the PATH.)
# After reseting, any subsequent initSh calls are guaranteed to start from
# a fresh context.
#
# The primary reason to do this is if the parent context is no longer reliable
# for some reason. Particularly in the case of daemonization, where the parent
# caller may clean up its initSh work directory while the daemon's still
# running. 
function dropInitShCall() {
    unset INIT_SH_PROJECT_CALL
    unset INIT_SH_INSTANCE_DIR
    unset INIT_SH_IMPORT_DIR
    unset INIT_SH_CALL_STACK_LEVEL
}

export -f dropInitShCall
