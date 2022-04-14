FROM ubuntu:focal
ARG UMAGNUSVERSION 2.2.1

RUN /bin/sh -c set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		ca-certificates \
		curl \
		netbase \
		wget \
	; \
	rm -rf /var/lib/apt/lists/*

RUN /bin/sh -c set -ex; \
	if ! command -v gpg > /dev/null; then \
		apt-get update; \
		apt-get install -y --no-install-recommends \
			gnupg \
			dirmngr \
		; \
		rm -rf /var/lib/apt/lists/*; \
	fi

RUN apt-get update; \
	apt-get install -y --no-install-recommends \
		git \
		mercurial \
		openssh-client \
		subversion \
		procps \
	; \
	rm -rf /var/lib/apt/lists/*

RUN /bin/sh -c set -eux; \
	apt-get update; \
	apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
	; \
	rm -rf /var/lib/apt/lists/*

ENV PATH /usr/local/golang/1.17.8/bin:$PATH

ENV GOLANG_VERSION 1.17.8

RUN mkdir -p /usr/local/golang

RUN /bin/sh -c set -eux; \
	arch="$(dpkg --print-architecture)"; arch="${arch##*-}"; \
	url=; \
	case "$arch" in \
		'amd64') \
			url='https://dl.google.com/go/go1.17.8.linux-amd64.tar.gz'; \
			sha256='980e65a863377e69fd9b67df9d8395fd8e93858e7a24c9f55803421e453f4f99'; \
			;; \
		'armel') \
			export GOARCH='arm' GOARM='5' GOOS='linux'; \
			;; \
		'armhf') \
			url='https://dl.google.com/go/go1.17.8.linux-armv6l.tar.gz'; \
			sha256='3287ca2fe6819fa87af95182d5942bf4fa565aff8f145812c6c70c0466ce25ae'; \
			;; \
		'arm64') \
			url='https://dl.google.com/go/go1.17.8.linux-arm64.tar.gz'; \
			sha256='57a9171682e297df1a5bd287be056ed0280195ad079af90af16dcad4f64710cb'; \
			;; \
		'i386') \
			url='https://dl.google.com/go/go1.17.8.linux-386.tar.gz'; \
			sha256='a826cd599828aeefc86d742e1a8ce8ab7e0251b2429568ad5633e21c8a769053'; \
			;; \
		'mips64el') \
			export GOARCH='mips64le' GOOS='linux'; \
			;; \
		'ppc64el') \
			url='https://dl.google.com/go/go1.17.8.linux-ppc64le.tar.gz'; \
			sha256='2077dd2fc57a74b0630b0c239ae4e3114607311778effd43fcfe5174133ee188'; \
			;; \
		's390x') \
			url='https://dl.google.com/go/go1.17.8.linux-s390x.tar.gz'; \
			sha256='3fac23801644a2f93a1643acecd5a94a5ea05d88e19467092fb6e64205710f61'; \
			;; \
		*) echo >&2 "error: unsupported architecture '$arch' (likely packaging update needed)"; exit 1 ;; \
	esac; \
	build=; \
	if [ -z "$url" ]; then \
		build=1; \
		url='https://dl.google.com/go/go1.17.8.src.tar.gz'; \
		sha256='2effcd898140da79a061f3784ca4f8d8b13d811fb2abe9dad2404442dabbdf7a'; \
		echo >&2; \
		echo >&2 "warning: current architecture ($arch) does not have a compatible Go binary release; will be building from source"; \
		echo >&2; \
	fi; \
		wget -O go.tgz.asc "$url.asc"; \
	wget -O go.tgz "$url" --progress=dot:giga; \
	echo "$sha256 *go.tgz" | sha256sum -c -; \
		GNUPGHOME="$(mktemp -d)"; export GNUPGHOME; \
	gpg --batch --keyserver keyserver.ubuntu.com --recv-keys 'EB4C 1BFD 4F04 2F6D DDCC  EC91 7721 F63B D38B 4796'; \
	gpg --batch --keyserver keyserver.ubuntu.com --recv-keys '2F52 8D36 D67B 69ED F998  D857 78BD 6547 3CB3 BD13'; \
	gpg --batch --verify go.tgz.asc go.tgz; \
	gpgconf --kill all; \
	rm -rf "$GNUPGHOME" go.tgz.asc; \
		tar -C /usr/local/golang -xzf go.tgz; \
	mv /usr/local/golang/go /usr/local/golang/1.17.8; \
	rm go.tgz; \
		if [ -n "$build" ]; then \
		savedAptMark="$(apt-mark showmanual)"; \
		apt-get update; \
		apt-get install -y --no-install-recommends golang-go; \
				( \
			cd /usr/local/golang/1.17.8/src; \
			export GOROOT_BOOTSTRAP="$(go env GOROOT)" GOHOSTOS="$GOOS" GOHOSTARCH="$GOARCH"; \
			./make.bash; \
		); \
				apt-mark auto '.*' > /dev/null; \
		apt-mark manual $savedAptMark > /dev/null; \
		apt-get purge -y --auto-remove -o APT::AutoRemove::RecommendsImportant=false; \
		rm -rf /var/lib/apt/lists/*; \
				rm -rf \
			/usr/local/golang/1.17.8/pkg/*/cmd \
			/usr/local/golang/1.17.8/pkg/bootstrap \
			/usr/local/golang/1.17.8/pkg/obj \
			/usr/local/golang/1.17.8/pkg/tool/*/api \
			/usr/local/golang/1.17.8/pkg/tool/*/go_bootstrap \
			/usr/local/golang/1.17.8/src/cmd/dist/dist \
		; \
	fi; \
		go version

ENV GOPATH /usr/local/umagnus/$UMAGNUSVERSION

ENV PATH /usr/local/umagnus/$UMAGNUSVERSION/bin:$PATH

RUN mkdir -p "$GOPATH/src/github.com/seeder-research" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

WORKDIR "$GOPATH/src/github.com/seeder-research"

RUN git clone https://github.com/seeder-research/uMagNUS -b 2.2.1

WORKDIR "$GOPATH/src/github.com/seeder-research/uMagNUS"

RUN make base

WORKDIR /usr/local/umagnus/2.2.1
