# 0ad-tool


clean up data xml
find . -type f -name '*.cached.xmb' -delete

# Icon

path to icon:
from: [gitea](https://gitea.wildfiregames.com/0ad/0ad/src/branch/main/binaries/data/mods/public/art/textures/ui/session/portraits/units)
`binaries/data/mods/public/art/textures/ui/session/portraits/units`


# Dev

`podman build -t 0ad .`

`podman run -it --rm 0ad`