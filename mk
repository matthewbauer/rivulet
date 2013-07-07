#!/bin/sh

./bump

echo 'making'

coffee --compile --output build coffee/*
#closure build/script.js > build/script.js

jade -o build jade
cp build/{landing,articles,feeds,user}.html templates

lessc less/style.less build/style.css

yuicompressor build/script.js -o build/script.js
yuicompressor build/bookmark.js -o build/bookmark.js
yuicompressor build/style.css -o build/style.css

sizes='16 24 32 48 57 64 72 96 114 128 144 195 256 512'
for size in $sizes
do
  inkscape -C -e build/icon$size.png -w $size -h $size static/icon.svg
done

#./bump

#cd crx
#zip -r ../rivulet.crx .
#cd ..

#cd wgt
#cp ../config.xml .
#cp ../icon/icon.png .
#cp ../build/offline.html index.html
#zip -r ../rivulet.wgt .
#zip -r ../rivulet.oex .
#cd ..

