#!/bin/sh
export ARGS0=$0
export BASE_DIR=$(dirname "$ARGS0")
export CURRENT_DIR=$(dirname "$BASE_DIR")
export CURRENT_DIR=$(dirname "$CURRENT_DIR")
export CURRENT_DIR=$(dirname "$CURRENT_DIR")
export EXEC_BIN="$BASE_DIR/projectLauncher-darwin"
export TMP_SCRIPT="/tmp/projectLauncher-$RANDOM.sh"
# echo $CURRENT_DIR
# echo $EXEC_BIN


# cd "$CURRENT_DIR" 
echo "#!/bin/sh\n sudo $EXEC_BIN -p $CURRENT_DIR" > $TMP_SCRIPT
chmod +x $TMP_SCRIPT
open -a Terminal.app "$TMP_SCRIPT"