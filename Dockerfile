# TODO: should be a tagged release
FROM dockcross/linux-armv7

ENV PATH="/usr/local/go/bin:$PATH"
WORKDIR /src

RUN dpkg --add-architecture armhf

# TODO: this ought to work on a standard Debian or Ubuntu machine or container
# instead of using dockcross. Ubuntu 20.04 conflicted on the armhf packages,
# however, so I've elected to do this in a container for now.
RUN apt-get update -y && apt-get install -y --no-install-recommends \
    g++-arm-linux-gnueabihf \
    gcc-arm-linux-gnueabihf \
    binutils-arm-linux-gnueabihf \
    software-properties-common \
    libgl-dev:armhf \
    libfreetype-dev:armhf \
    libxi-dev:armhf \
    libxcursor-dev:armhf \
    libxrandr-dev:armhf \
    libxxf86vm-dev:armhf \
    libxinerama-dev:armhf \
    && rm -rf /var/lib/apt/lists/*

# install Go 1.18
RUN curl -o /tmp/go.tar.gz -L https://go.dev/dl/go1.18.2.linux-amd64.tar.gz && tar -C /usr/local -xzf /tmp/go.tar.gz && rm -f /tmp/go.tar.gz
