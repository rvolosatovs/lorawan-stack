# Copyright © 2018 The Things Network Foundation, The Things Industries B.V.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Include this makefile to enable webpack-related rules

# The location of the config files
CONFIG_DIR ?= config

# The place where we keep intermediate build files
CACHE_DIR ?= .cache

# Webpack
WEBPACK ?= ./node_modules/.bin/webpack
WEBPACK_FLAGS ?= --colors $(if $(CI),,--progress)

# The config file to use for client
WEBPACK_CONFIG ?= $(CONFIG_DIR)/webpack.config.js

# Pre-build config files for quicker builds
$(CACHE_DIR)/config/%.js: $(CONFIG_DIR)/%.js
	@$(log) pre-building config files [babel $<]
	@mkdir -p $(CACHE_DIR)/config
	@$(BABEL) $< >| $@

# The location of the cached config file
WEBPACK_CONFIG_BUILT = $(subst $(CONFIG_DIR)/,$(CACHE_DIR)/config/,$(WEBPACK_CONFIG))

js.build: $(PUBLIC_DIR)/console.html

$(PUBLIC_DIR)/console.html: $(WEBPACK_CONFIG_BUILT) $(shell $(JS_SRC_FILES)) $(JS_SRC_DIR)/index.html package.json yarn.lock
	@$(log) "building client [webpack -c $(WEBPACK_CONFIG_BUILT) $(WEBPACK_FLAGS)]"
	@$(JS_ENV) $(WEBPACK) --config $(WEBPACK_CONFIG_BUILT) $(WEBPACK_FLAGS)

# build in dev mode
js.build-dev: NODE_ENV =
js.build-dev: js.build

# watch files
.PHONY: js.watch
js.watch: NODE_ENV = development
js.watch: WEBPACK_FLAGS += -w
js.watch: js.build

$(CACHE_DIR)/make/%.js: .make/js/%.js
	@$(log) "pre-building translation scrips [babel $<]"
	@mkdir -p $(CACHE_DIR)/make
	@$(BABEL) $< >| $@

SUPPORT_LOCALES ?= en,ja
DEFAULT_LOCALE ?= en
OUTPUT_MESSAGES ?= messages.yml,messages.xlsx
UPDATES_FILES ?= messages.xlsx

# update translations
js.translations: $(CACHE_DIR)/make/translations.js $(CACHE_DIR)/make/xls.js $(CACHE_DIR)/make/xx.js
	@$(log) "gathering translations [translations --output messages.{yml,xlsx]"
	@$(NODE) $(CACHE_DIR)/make/translations.js --output $(OUTPUT_MESSAGES) --support $(SUPPORT_LOCALES) --updates $(UPDATES_FILES)

# vim: ft=make
