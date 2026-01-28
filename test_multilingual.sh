#!/bin/bash
# Test Hugo multilingual switching

echo "=== Testing Hugo Multilingual Configuration ==="
echo ""

# Check if Hugo is installed
if ! command -v hugo &>/dev/null; then
	echo "Error: Hugo is not installed. Please install Hugo first."
	exit 1
fi

# Change to website directory
cd "$(dirname "$0")/website"

# Check configuration
echo "1. Checking Hugo configuration..."
echo "Default language: $(grep defaultContentLanguage hugo.toml | cut -d' ' -f3)"
echo ""

# Check content structure
echo "2. Checking content structure..."
echo "Content directories:"
ls -la content/ | grep -E "en|zh|docs"
echo ""

# Run Hugo in server mode for testing (but don't wait)
echo "3. Running Hugo to generate site..."
hugo --minify --buildDrafts --buildFuture >/tmp/hugo_build.log 2>&1

if [ $? -eq 0 ]; then
	echo "✓ Hugo build successful!"
	echo "Site generated in public/ directory"
	echo ""
	echo "To test locally, run:"
	echo "  cd website && hugo server -D"
	echo ""
	echo "Then visit:"
	echo "  - English: http://localhost:1313"
	echo "  - Chinese: http://localhost:1313/zh/"
else
	echo "✗ Hugo build failed. Check the log:"
	tail -n 20 /tmp/hugo_build.log
fi

# Clean up
rm -f /tmp/hugo_build.log
