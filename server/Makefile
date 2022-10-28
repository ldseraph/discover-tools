.PHONY: build
build:
	go build -o discover-service *.go 

.PHONY: dev
dev:
	go run *.go

.PHONY: install
install:
	install discover-service /usr/local/bin/discover-service
	install discover.service /etc/systemd/system/discover.service
	systemctl daemon-reload
	systemctl enable discover.service
	systemctl restart discover.service

.PHONY: deb
deb:
	mkdir -pv build/discover_service_0.1_arm64/DEBIAN
	mkdir -pv build/discover_service_0.1_arm64/usr/local/bin
	mkdir -pv build/discover_service_0.1_arm64/etc/systemd/system/
	cp -r DEBIAN build/discover_service_0.1_arm64
	cp discover-service build/discover_service_0.1_arm64/usr/local/bin/discover-service
	cp discover.service build/discover_service_0.1_arm64/etc/systemd/system/discover.service
	cd build && dpkg-deb --build --root-owner-group discover_service_0.1_arm64
	mv build/discover_service_0.1_arm64.deb ./