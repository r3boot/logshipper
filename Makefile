NAME = "logshipper"

BUILDROOT = "work/buildroot"

WHEEZY_BRANCH = "debian_wheezy"
WHEEZY_PKGS = "build/wheezy"

JESSIE_BRANCH = "debian_jessie"
JESSIE_PKGS = "build/jessie"

SRCDIR = ${PWD}

all: clean $(NAME)

$(NAME):
	go build -v

jessie-package:
	git checkout master
	git branch -D ${JESSIE_BRANCH} || true
	git checkout -b ${JESSIE_BRANCH}
	install -dv ${BUILDROOT}
	install -dv ${BUILDROOT}/etc
	install -dv ${BUILDROOT}/usr/sbin
	install -dv ${BUILDROOT}/lib/systemd/system
	install -dv ${BUILDROOT}/var/lib/logshipper
	install -v -m 0755 logshipper ${BUILDROOT}/usr/sbin/logshipper
	install -v -m 0644 logshipper.yml ${BUILDROOT}/etc/logshipper.yml
	install -v -m 0644 rcscripts/logshipper.service \
		${BUILDROOT}/lib/systemd/system/logshipper.service
	install -v -m 0644 README.md ${BUILDROOT}/README.Debian
	install -v -m 0644 debian/changelog ${BUILDROOT}/changelog
	(cd ${BUILDROOT}; equivs-build ${SRCDIR}/debian/logshipper-jessie.equivs)
	install -dv ${JESSIE_PKGS}
	cp -v ${BUILDROOT}/logshipper_*_amd64.deb ${JESSIE_PKGS}
	git checkout master
	git branch -D ${JESSIE_BRANCH}

wheezy-package:
	git checkout master
	git branch -D ${WHEEZY_BRANCH} || true
	git checkout -b ${WHEEZY_BRANCH}
	install -dv ${BUILDROOT}
	install -dv ${BUILDROOT}/etc
	install -dv ${BUILDROOT}/usr/sbin
	install -dv ${BUILDROOT}/etc/init.d
	install -dv ${BUILDROOT}/var/lib/logshipper
	install -v -m 0755 logshipper ${BUILDROOT}/usr/sbin/logshipper
	install -v -m 0644 logshipper.yml ${BUILDROOT}/etc/logshipper.yml
	install -v -m 0644 rcscripts/logshipper.rc \
		${BUILDROOT}/etc/init.d/logshipper
	install -v -m 0644 README.md ${BUILDROOT}/README.Debian
	install -v -m 0644 debian/changelog ${BUILDROOT}/changelog
	(cd ${BUILDROOT}; equivs-build ${SRCDIR}/debian/logshipper-wheezy.equivs)
	install -dv ${WHEEZY_PKGS}
	cp -v ${BUILDROOT}/logshipper_*_amd64.deb ${WHEEZY_PKGS}
	git checkout master
	git branch -D ${WHEEZY_BRANCH}

debian-packages: wheezy-package jessie-package

clean:
	rm -f logshipper
