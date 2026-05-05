# Justfile for pngdeity-github-io

HUGO_VERSION := "0.154.4"
DART_SASS_VERSION := "1.97.2"

# Default recipe
default:
    @just --list

# The entrypoint for CI
ci: build

# Main build recipe
build: setup
    #!/usr/bin/env bash
    set -euo pipefail
    echo "--- Executing Build ---"
    
    # Add local bin to path
    export PATH="$PWD/ci/bin:$PATH"
    
    # Clean previous build artifacts
    rm -rf ci/out
    mkdir -p ci/out/{blog,app}
    
    # 1. Build Hugo Blog
    echo "Building Hugo Blog..."
    # Ensure theme dependencies are installed if package.json exists in theme
    if [ -f "hugo-src/themes/ananke/package.json" ]; then
        echo "Installing theme dependencies..."
        cd hugo-src/themes/ananke && npm ci
        cd ../../..
    fi
    hugo --gc --minify --source ./hugo-src
    
    # 2. Build Blazor App
    echo "Building Blazor App..."
    dotnet publish my-blazor-app/BlazorApp.csproj -c Release -o release
    
    echo "Updating Blazor base href..."
    # Portable replacement for sed -i using perl
    perl -pi -e 's|<base href="/" />|<base href="/app/" />|g' release/wwwroot/index.html
    
    # 3. Standardize Output
    echo "Packaging Deployment Artifacts..."
    cp -r hugo-src/public/* ./ci/out/blog/
    cp -r release/wwwroot/* ./ci/out/app/
    
    # Prevent GitHub Pages from ignoring folders with underscores
    touch ./ci/out/.nojekyll
    
    echo "Build execution finished successfully."

# Setup toolchain binaries
setup:
    #!/usr/bin/env bash
    set -euo pipefail
    
    mkdir -p ci/bin
    
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    # Normalize Architecture names
    if [ "$ARCH" = "x86_64" ]; then
        HUGO_ARCH="amd64"
        SASS_ARCH="x64"
    elif [ "$ARCH" = "arm64" ] || [ "$ARCH" = "aarch64" ]; then
        HUGO_ARCH="arm64"
        SASS_ARCH="arm64"
    else
        echo "Unsupported architecture: $ARCH"
        exit 1
    fi
    
    # Detect OS-specific naming
    if [ "$OS" = "darwin" ]; then
        HUGO_OS="darwin"
        SASS_OS="macos"
        # Hugo uses 'universal' for macOS usually, or specific if available
        HUGO_PLATFORM="darwin-universal"
    elif [ "$OS" = "linux" ]; then
        HUGO_OS="linux"
        SASS_OS="linux"
        HUGO_PLATFORM="linux-${HUGO_ARCH}"
    else
        echo "Unsupported OS: $OS"
        exit 1
    fi

    # 1. Install Hugo
    if [ ! -f "ci/bin/hugo" ] || [ "$(ci/bin/hugo version | awk '{print $2}' | cut -d 'v' -f 2)" != "{{HUGO_VERSION}}" ]; then
        echo "Installing Hugo Extended v{{HUGO_VERSION}} for ${HUGO_PLATFORM}..."
        URL="https://github.com/gohugoio/hugo/releases/download/v{{HUGO_VERSION}}/hugo_extended_{{HUGO_VERSION}}_${HUGO_PLATFORM}.tar.gz"
        curl -sL "$URL" | tar -xz -C ci/bin hugo
    else
        echo "Hugo v{{HUGO_VERSION}} already installed."
    fi

    # 2. Install Dart Sass
    if [ ! -f "ci/bin/sass" ]; then
        echo "Installing Dart Sass v{{DART_SASS_VERSION}} for ${SASS_OS}-${SASS_ARCH}..."
        URL="https://github.com/sass/dart-sass/releases/download/{{DART_SASS_VERSION}}/dart-sass-{{DART_SASS_VERSION}}-${SASS_OS}-${SASS_ARCH}.tar.gz"
        # Extract and move binary. Dart Sass tarball has a nested folder.
        curl -sL "$URL" | tar -xz -C ci/bin --strip-components=1
    else
        echo "Dart Sass v{{DART_SASS_VERSION}} already installed."
    fi

    echo "Environment setup complete."
