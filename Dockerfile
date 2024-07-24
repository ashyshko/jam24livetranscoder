FROM debian:latest AS builder

RUN apt update && apt install -y build-essential git curl zip unzip tar

RUN mkdir /home/tools && cd /home/tools && \
    curl -o cmake.tar.gz -L https://github.com/Kitware/CMake/releases/download/v3.30.1/cmake-3.30.1-linux-x86_64.tar.gz && \
    tar xf cmake.tar.gz && \
    mv cmake-3.30.1-linux-x86_64 cmake



ENV PATH="/home/tools/cmake/bin:${PATH}" \
    CC="/usr/bin/clang" \
    CXX="/usr/bin/clang++" \
    VCPKG_ROOT="/opt/vcpkg"

RUN cd /opt && \
    git clone https://github.com/microsoft/vcpkg.git && \
    ./vcpkg/bootstrap-vcpkg.sh

WORKDIR /home/src/wdp_transcoder_lite

# RUN update-alternatives --set python3 /usr/bin/python3.8

COPY wdp_transcoder_lite/src/vcpkg.json wdp_transcoder_lite/src/vcpkg-configuration.json ./

RUN apt install -y clang pkg-config nasm

COPY wdp_transcoder_lite/src/.vcpkg .vcpkg


RUN /opt/vcpkg/vcpkg install

COPY wdp_transcoder_lite/src .

ARG VERSION=local
RUN echo $VERSION >/home/VERSION

RUN cmake -S . -B /home/build/wdp_transcoder_lite \
    -DCMAKE_BUILD_TYPE=Release \
    "-DAPP_VERSION:STRING=$VERSION" \
    "-DCMAKE_TOOLCHAIN_FILE=$VCPKG_ROOT/scripts/buildsystems/vcpkg.cmake" \
    -DCMAKE_INSTALL_PREFIX=/home/install/wdp_transcoder_lite

RUN cmake --build /home/build/wdp_transcoder_lite

RUN cmake --install /home/build/wdp_transcoder_lite

FROM golang:bookworm AS servicebuilder

WORKDIR /work

COPY go.mod go.sum ./

RUN go mod download

COPY protocol ./protocol

WORKDIR /work/transcoder
COPY transcoder/*.go . 

RUN go build

FROM debian:latest AS final

WORKDIR /home

RUN apt-get update && apt-get install libatomic1

COPY --from=builder /home/install/wdp_transcoder_lite/libwdp_transcoder_lite.so .
COPY --from=servicebuilder /work/transcoder/transcoder . 

EXPOSE 8898

ENTRYPOINT [ "./transcoder" ]