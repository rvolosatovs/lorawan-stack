// Copyright © 2021 The Things Network Foundation, The Things Industries B.V.
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

import { defineSmokeTest } from '../utils'

const validatePasswordLinkRegExp = `http:\\/\\/localhost:\\d{4}\\/[a-zA-Z0-9-_]+\\/validate\\?.+&token=[A-Z0-9]+`

const contactInfoValidation = defineSmokeTest('succeeds validating contact info', () => {
  const user = {
    ids: { user_id: 'test-user-id1' },
    primary_email_address: 'test-user1@example.com',
    password: 'ABCDefg123!',
    password_confirm: 'ABCDefg123!',
  }

  cy.createUser(user)

  cy.task('findInStackLog', validatePasswordLinkRegExp).then(validationUri => {
    cy.log(validationUri)
    cy.visit(validationUri)
    cy.findByTestId('notification')
      .should('be.visible')
      .should('contain', 'Validation successful')
  })

  cy.reload()
  cy.findByTestId('error-notification')
    .should('be.visible')
    .should('contain', 'token already used')
})

export default [contactInfoValidation]
