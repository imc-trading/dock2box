all: readme

readme:
	godoc2md github.com/mickep76/hwinfo | grep -v Generated >README.md ;\
	for dir in $$(find * -maxdepth 1 -type d); do \
		godoc2md github.com/mickep76/hwinfo/$${dir} | grep -v Generated >>README.md ;\
	done
