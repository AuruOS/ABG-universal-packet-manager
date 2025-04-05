PREFIX := /usr
DESTDIR := /
BINARY_NAME := abg

GO := go

all: clean build

build: ${BINARY_NAME}

${BINARY_NAME}:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ${GO} build -a -tags netgo -ldflags '-w -extldflags "-static"' -o $@

install: build
	install -Dm755 ${BINARY_NAME} ${DESTDIR}${PREFIX}/bin/${BINARY_NAME}
	mkdir -p ${DESTDIR}/etc/abg
	sed -i 's|/usr/share/abg/distrobox|${PREFIX}/share/abg/distrobox|g' config/abg.json
	install -Dm644 config/abg.json ${DESTDIR}/etc/abg/apx.json
	mkdir -p ${DESTDIR}${PREFIX}/share/abg/distrobox
	sh distrobox/install --prefix ${DESTDIR}${PREFIX}/share/abg/distrobox
	mv ${DESTDIR}${PREFIX}/share/abg/distrobox/bin/distrobox* ${DESTDIR}${PREFIX}/share/abg/distrobox/.

install-manpages:
	mkdir -p ${DESTDIR}${PREFIX}/share/man/man1
	cp -r man/* ${DESTDIR}${PREFIX}/share/man/.
	chmod 644 ${DESTDIR}${PREFIX}/share/man/man1/abg*

uninstall: uninstall-manpages
	rm ${DESTDIR}${PREFIX}/bin/${BINARY_NAME}
	rm -rf ${DESTDIR}/etc/abg
	rm -rf ${DESTDIR}${PREFIX}/share/abg

uninstall-manpages:
	rm -rf ${DESTDIR}${PREFIX}/share/man/man1/abg*

clean:
	rm -f ${BINARY_NAME}
	${GO} clean
