#!/bin/bash

HOOK_DIR=.git/hooks
HOOK_FILES=.github/hooks/*

for hook in $HOOK_FILES; do
    # Get just the filename
    hook_name=$(basename "$hook")
    # Create the symlink, overwriting any existing ones
    ln -sf "../../.github/hooks/$hook_name" "$HOOK_DIR/$hook_name"
done

echo "Git hooks installed!"