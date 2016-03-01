#!/sbin/runscript
# Copyright 1999-2004 Gentoo Foundation
# Distributed under the terms of the GNU General Public License v2
# $Header: /var/cvsroot/gentoo-x86/app-admin/pwgen/files/pwgen.rc,v 1.4 2006/09/05 16:58:01 wolf31o2 Exp $

depend() {
	before local
}

start() {
	ebegin "Setting default Dock2Box bootstrap image root password"
    # this is a bogus password and should be changed to whatever is appropriate for your environment before generating the bootstrap image
    # for example purposes, this INSECURE password is: "Welcome@d2b".  CHANGE THIS to something unique to your organization!
    usermod -p '{{ ENCRYPTED_ROOT_PASSWORD }}' root
	eend $? "Failed to set root password."
}

stop() {
	ebegin "Stopping pwgen"
	eend $? "Failed to stop pwgen."
}
