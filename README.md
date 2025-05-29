# 0ad-tool


clean up data xml
find . -type f -name '*.cached.xmb' -delete


clean up icon 
find . -name "*cached.dds*" -type f -exec bash -c 'mv "$1" "${1//.cached.dds/}"' _ {} \;

find . -name "*.cached.dds" -exec bash -c 'ffmpeg -i "$1" "${1%.cached.dds}.png"' _ {} \;

# Dev

`podman build -t 0ad .`

`podman run -it --rm 0ad`