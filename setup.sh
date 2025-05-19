#!/bin/bash

# --- Configuration: List your Go tools here ---
# Add each tool's import path to this array
declare -a GO_TOOLS=(
    "github.com/magefile/mage"
    "github.com/air-verse/air"
    "github.com/golangci/golangci-lint/v2"
)
# --- End Configuration ---

echo "Checking and installing Go tools..."

# Determine GOBIN or default GOPATH/bin
# GOBIN takes precedence if set
if [ -n "$GOBIN" ]; then
    TOOL_INSTALL_PATH="$GOBIN"
elif [ -n "$GOPATH" ]; then
    TOOL_INSTALL_PATH="$GOPATH/bin"
else
    # Default Go path if neither GOBIN nor GOPATH are set
    TOOL_INSTALL_PATH="$HOME/go/bin"
fi

# Ensure the tool installation path is in the PATH for `command -v` to work reliably for newly installed tools in the same session
# This is more for the user's convenience after the script runs.
# The `command -v` check itself will search the current PATH.
export PATH=$TOOL_INSTALL_PATH:$PATH

for tool_path in "${GO_TOOLS[@]}"; do

    # Extract the binary name from the tool path.
    # This handles common Go tool path patterns:
    # 1. Standard path: "golang.org/x/tools/cmd/godoc" -> "godoc"
    # 2. Simple name: "mvdan.cc/gofumpt" -> "gofumpt"
    # 3. Path with version suffix: "github.com/golangci/golangci-lint/v2" -> "golangci-lint"
    # 4. Name with @version suffix: "tool@v1.2.3" -> "tool"
    
    potential_last_segment=$(basename "$tool_path")

    # Check if the last path segment is a version (e.g., v2, v3)
    # and the tool_path contains directory separators.
    # Corrected regex: ^v[0-9]+$
    if [[ "$potential_last_segment" =~ ^v[0-9]+$ && "$tool_path" == *"/"* ]]; then
        # If so, the binary name is the directory name before the version segment
        # e.g., "path/to/tool/v2" -> "tool"
        binary_name=$(basename "$(dirname "$tool_path")")
    else
        # Otherwise, the last segment is the starting point for the binary name
        # e.g., "path/to/tool" -> "tool", or "tool@version" -> "tool@version"
        binary_name="$potential_last_segment"
    fi

    # Remove any @version suffix from the determined binary name
    # e.g., "tool@v1.2.3" -> "tool", or "gopls@latest" -> "gopls"
    binary_name=${binary_name%@*}

    echo "----------------------------------------"
    echo "Processing: $tool_path (binary: $binary_name)"

    # Check if the binary is already in PATH and executable
    if command -v "$binary_name" &> /dev/null; then
        echo "$binary_name is already installed and in PATH."
        echo "Location: $(command -v "$binary_name")"
    else
        echo "$binary_name not found in PATH. Attempting to install..."
        # Using @latest to get the newest version.
        if go install "${tool_path}@latest"; then
            echo "Successfully installed $tool_path (binary: $binary_name)."
            # Verify installation again, in case go install puts it somewhere unexpected
            # or if the binary name differs from the package name component.
            if command -v "$binary_name" &> /dev/null; then
                 echo "Verified: $binary_name is now in PATH at $(command -v "$binary_name")"
            else
                 echo "Warning: $binary_name still not found in PATH after installation. Check your GOBIN/GOPATH and PATH settings."
                 echo "Expected install location: $TOOL_INSTALL_PATH/$binary_name"
            fi
        else
            echo "Error installing $tool_path."
        fi
    fi
done

echo "----------------------------------------"
echo "All specified tools processed."