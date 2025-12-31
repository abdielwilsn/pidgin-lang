# Extension Images

This directory contains images for the VSCode extension.

## Required Images

### icon.png
- **Size**: 128x128 pixels
- **Format**: PNG
- **Purpose**: Extension icon shown in VSCode marketplace and extension list
- **Recommendation**: Create a simple logo with "PDG" or Pidgin-related imagery

### screenshot.png (optional but recommended)
- **Size**: 1280x720 pixels or similar 16:9 ratio
- **Format**: PNG
- **Purpose**: Main screenshot for marketplace listing
- **Content**: Show syntax-highlighted Pidgin code in VSCode

## Creating the Icon

You can create a simple icon using:

1. **Online tools**:
   - Canva: https://www.canva.com
   - Figma: https://www.figma.com
   - Pixlr: https://pixlr.com

2. **Design suggestions**:
   - Nigerian flag colors (green and white)
   - Text: "PDG" or "Pidgin"
   - Simple, recognizable at small sizes

3. **Command-line (macOS with ImageMagick)**:
   ```bash
   # Install ImageMagick
   brew install imagemagick

   # Create simple icon with text
   convert -size 128x128 xc:#008751 -gravity center \
     -font Arial-Bold -pointsize 48 -fill white \
     -annotate +0+0 'PDG' icon.png
   ```

## Taking Screenshots

1. Open a `.pdg` file in VSCode with syntax highlighting enabled
2. Choose a good color theme (e.g., "Dark+ (default dark)")
3. Use VSCode's built-in screenshot feature or system screenshot tool
4. Crop to show just the editor window with highlighted code
5. Save as `screenshot.png`

## Example Screenshot Content

```pidgin
# Fibonacci sequence
make a be 0
make b be 1

dey do while b no reach 100 {
    yarn(b)
    make temp be a + b
    make a be b
    make b be temp
}
```

## Temporary Icon

Until you create a custom icon, the extension will use VSCode's default icon. This doesn't affect functionality, only appearance in the marketplace.

## Adding Images to Extension

Once you have the images:

1. Place them in this directory
2. Update `package.json`:
   ```json
   {
     "icon": "images/icon.png"
   }
   ```
3. Rebuild the extension: `vsce package`
