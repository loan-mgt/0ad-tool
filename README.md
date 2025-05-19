# 0ad-tool


clean up data xml
find . -type f -name '*.cached.xmb' -delete


clean up icon 
find . -name "*cached.dds*" -type f -exec bash -c 'mv "$1" "${1//.cached.dds/}"' _ {} \;

