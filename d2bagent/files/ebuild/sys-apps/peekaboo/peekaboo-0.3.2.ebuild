# Copyright 1999-2015 Gentoo Foundation
# Distributed under the terms of the GNU General Public License v2
# $Id$
# By Jean-Michel Smith, first created 9/4/15

EAPI=5

inherit user git-r3 systemd

DESCRIPTION="Expose hardware info using JSON/REST and provide a system HTML Front-End"
HOMEPAGE="https://github.com/imc-trading/peekaboo.git"
SRC_URI=""

LICENSE="Apache-2.0"
SLOT="0"
KEYWORDS="amd64"
IUSE=""

DEPEND="dev-lang/go"
RDEPEND="${DEPEND}
	sys-apps/ethtool
	sys-apps/hwdata-redhat"

EGIT_REPO_URI="https://github.com/imc-trading/peekaboo.git"
EGIT_COMMIT="${PV}"

PEEKABOODIR="${WORKDIR}/peekaboo-${PV}"
PEEKABOOSRC="${PEEKABOODIR}/src/github.com/imc-trading/peekaboo"

src_compile() {
	ebegin "Building peekaboo ${PV}"
	export GOPATH=${PEEKABOODIR}
	export PATH=${GOPATH}/bin:${PATH}
	cd ${PEEKABOODIR}
	./build
	cd
	eend ${?}
}

src_install() {
	ebegin "installing peekaboo ${PV}"
	dobin ${PEEKABOODIR}/bin/peekaboo
	newinitd ${FILESDIR}/peekaboo.init peekaboo
	systemd_dounit "${FILESDIR}"/peekaboo.service
	newconfd ${FILESDIR}/peekaboo.conf peekaboo
	# templates and static pages:
	dodir /var/lib/peekaboo
	cp -a "${PEEKABOOSRC}"/templates "${D}"/var/lib/peekaboo
	cp -a "${PEEKABOOSRC}"/static "${D}"/var/lib/peekaboo
}
