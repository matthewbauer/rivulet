#!/bin/sh

echo 'making'

mkdir -p build templates

coffee --compile --output build coffee/script.coffee

jade -o build jade/{api,about,offline,extras,help}.jade
jade -o templates jade/{landing,articles,feeds,user}.jade

lessc less/style.less build/style.css

yuicompressor build/script.js -o build/script.min.js
yuicompressor build/style.css -o build/style.min.css

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

