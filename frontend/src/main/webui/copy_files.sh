#!/bin/bash
echo Copying files...

CSS_DIR=../resources/assets/lib/css
FONTS_DIR=../resources/assets/lib/fonts
IMG_DIR=../resources/assets/lib/img

mkdir -p $CSS_DIR $FONTS_DIR $IMG_DIR

cp node_modules/bootstrap/dist/css/*   $CSS_DIR/
cp node_modules/bootstrap/dist/fonts/* $FONTS_DIR/

cp node_modules/font-awesome/css/*   $CSS_DIR/
cp node_modules/font-awesome/fonts/* $FONTS_DIR/

# cp node_modules/react-toggle/style.css $CSS_DIR/toggle-style.css

# cp -r node_modules/react-widgets/dist/css/*   $CSS_DIR/
# cp -r node_modules/react-widgets/dist/fonts/* $FONTS_DIR/
# cp -r node_modules/react-widgets/dist/img/*   $IMG_DIR/
