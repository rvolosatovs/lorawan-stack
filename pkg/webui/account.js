// Copyright © 2020 The Things Network Foundation, The Things Industries B.V.
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

import React from 'react'
import DOM from 'react-dom'
import { Provider } from 'react-redux'
import { createBrowserHistory } from 'history'
import * as Sentry from '@sentry/browser'

import sentryConfig from '@ttn-lw/constants/sentry'

import WithLocale from './lib/components/with-locale'
import { EnvProvider } from './lib/components/env'
import Init from './lib/components/init'
import env from './lib/env'
import { selectApplicationRootPath } from './lib/selectors/env'
import './lib/yup'
import createStore from './account/store'

const appRoot = selectApplicationRootPath()
const history = createBrowserHistory({ basename: `${appRoot}/` })
// Initialize sentry before creating store
if (env.sentryDsn) Sentry.init(sentryConfig)
const store = createStore(history)
const rootElement = document.getElementById('app')

const render = () => {
  const App = require('./account/views/app').default

  DOM.render(
    <EnvProvider env={env}>
      <Provider store={store}>
        <WithLocale>
          <Init>
            <App history={history} env={env} />
          </Init>
        </WithLocale>
      </Provider>
    </EnvProvider>,
    rootElement,
  )
}

if (module.hot) {
  module.hot.accept('./account/views/app', () => {
    setTimeout(render)
  })
}

render()
