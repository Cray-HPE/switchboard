# Copyright 2019-2021 Hewlett Packard Enterprise Development LP
#
# Permission is hereby granted, free of charge, to any person obtaining a
# copy of this software and associated documentation files (the "Software"),
# to deal in the Software without restriction, including without limitation
# the rights to use, copy, modify, merge, publish, distribute, sublicense,
# and/or sell copies of the Software, and to permit persons to whom the
# Software is furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included
# in all copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
# FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.  IN NO EVENT SHALL
# THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR
# OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
# ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
# OTHER DEALINGS IN THE SOFTWARE.
#
# (MIT License)

# DOCKER
NAME ?= cray-switchboard
VERSION ?= $(shell cat .version)

# RPM
SPEC_FILE ?= ${NAME}.spec
SPEC_VERSION ?= $(shell cat .version)
BUILD_METADATA ?= 1~development~$(shell git rev-parse --short HEAD)
SOURCE_NAME ?= ${NAME}-${SPEC_VERSION}
BUILD_DIR ?= $(PWD)/dist/rpmbuild
SOURCE_PATH := ${BUILD_DIR}/SOURCES/${SOURCE_NAME}.tar.bz2

rpm: rpm_prepare rpm_package_source rpm_build_source rpm_build

image:
	docker build --pull ${BROKER_DOCKER_ARGS} --tag '${BROKER_NAME}:${VERSION}' .

test:
	./runUnitTest.sh

rpm_prepare:
	rm -rf $(BUILD_DIR)
	mkdir -p $(BUILD_DIR)/SPECS $(BUILD_DIR)/SOURCES
	cp $(SPEC_FILE) $(BUILD_DIR)/SPECS/

rpm_package_source:
	tar --transform 'flags=r;s,^,/$(SOURCE_NAME)/,' --exclude .git --exclude dist -cvjf $(SOURCE_PATH) .

rpm_build_source:
		BUILD_METADATA=$(BUILD_METADATA) rpmbuild -ts $(SOURCE_PATH) --nodeps --define "_topdir $(BUILD_DIR)"

rpm_build:
		BUILD_METADATA=$(BUILD_METADATA) rpmbuild -ba $(SPEC_FILE) --nodeps --define "_topdir $(BUILD_DIR)"
