// Copyright © 2019 The Things Network Foundation, The Things Industries B.V.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

export const configSelector = () => window

export const selectApplicationRootPath = () => configSelector().APP_ROOT

export const selectAssetsRootPath = () => configSelector().ASSETS_ROOT

export const selectApplicationConfig = () => configSelector().APP_CONFIG

export const selectApplicationSiteName = () => configSelector().SITE_NAME

export const selectApplicationSiteTitle = () => configSelector().SITE_TITLE

export const selectApplicationSiteSubTitle = () => configSelector().SITE_SUB_TITLE

export const selectStackConfig = () => selectApplicationConfig().stack_config

export const selectGsConfig = () => selectStackConfig().gs

export const selectIsConfig = () => selectStackConfig().is

export const selectNsConfig = () => selectStackConfig().ns

export const selectJsConfig = () => selectStackConfig().js

export const selectAsConfig = () => selectStackConfig().as

export const selectLanguageConfig = () => selectApplicationConfig().language

export const selectSupportLinkConfig = () => selectApplicationConfig().support_link

export const selectPageData = () => configSelector().PAGE_DATA
