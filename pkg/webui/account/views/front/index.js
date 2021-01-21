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
import { Switch, Route, Redirect } from 'react-router'

import Footer from '@ttn-lw/components/footer'

import Header from '@account/containers/header'

import Login from '@account/views/login'
import CreateAccount from '@account/views/create-account'
import ForgotPassword from '@account/views/forgot-password'
import UpdatePassword from '@account/views/update-password'
import FrontNotFound from '@account/views/front-not-found'
import Validate from '@account/views/validate'
import Code from '@account/views/code'

import style from './front.styl'

const FrontView = () => {
  return (
    <div className={style.container}>
      <section className={style.content}>
        <Header className={style.header} />
        <div className={style.main}>
          <Switch>
            <Route path="/login" component={Login} />
            <Route path="/register" component={CreateAccount} />
            <Route path="/forgot-password" component={ForgotPassword} />
            <Route path="/update-password" component={UpdatePassword} />
            <Route path="/validate" component={Validate} />
            <Route path="/code" component={Code} />
            <Redirect exact from="/" to="/login" />
            <Route component={FrontNotFound} />
          </Switch>
        </div>
      </section>
      <section className={style.visual} />
      <Footer className={style.footer} />
    </div>
  )
}

export default FrontView
