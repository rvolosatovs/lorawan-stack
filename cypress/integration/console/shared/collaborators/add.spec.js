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

describe('Collaborators', () => {
  const userNotExist = 'not-exist-user'
  const orgNotExist = 'not-exist-org'
  const userId = 'main-collab-user'
  const user = {
    ids: { user_id: userId },
    primary_email_address: 'main-collab-user@example.com',
    password: 'ABCDefg123!',
    password_confirm: 'ABCDefg123!',
  }
  const collaboratorId = 'collab-test-user'
  const collaboratorUser = {
    ids: { user_id: collaboratorId },
    primary_email_address: 'collab-test-user@example.com',
    password: 'ABCDefg123!',
    password_confirm: 'ABCDefg123!',
  }
  const orgUserId = 'org-test-user'
  const orgUser = {
    ids: { user_id: orgUserId },
    primary_email_address: 'org-test-user@example.com',
    password: 'ABCDefg123!',
    password_confirm: 'ABCDefg123!',
  }
  const organizationId = 'test-collab-org'
  const organization = {
    ids: { organization_id: organizationId },
  }

  before(() => {
    cy.dropAndSeedDatabase()

    cy.createUser(user)
    cy.createUser(collaboratorUser)
    cy.createUser(orgUser)

    cy.loginConsole({
      user_id: orgUserId,
      password: orgUser.password,
    })
    cy.createOrganization(organization, orgUserId)
    cy.clearLocalStorage()
    cy.clearCookies()
  })

  describe('Application', () => {
    const applicationId = 'collab-test-app'
    const application = { ids: { application_id: applicationId } }

    before(() => {
      cy.loginConsole({ user_id: userId, password: user.password })
      cy.createApplication(application, userId)
      cy.clearLocalStorage()
      cy.clearCookies()
    })

    beforeEach(() => {
      cy.loginConsole({ user_id: user.ids.user_id, password: user.password })
      cy.visit(
        `${Cypress.config('consoleRootPath')}/applications/${applicationId}/collaborators/add`,
      )
    })

    it('succeeds adding user as a collaborator', () => {
      cy.findByLabelText('Collaborator ID').type(collaboratorId)
      cy.findByLabelText('User').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification').should('not.exist')
      cy.location('pathname').should(
        'eq',
        `${Cypress.config('consoleRootPath')}/applications/${applicationId}/collaborators`,
      )
    })

    it('fails adding non-existent user', () => {
      cy.findByLabelText('Collaborator ID').type(userNotExist)
      cy.findByLabelText('User').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification')
        .should('be.visible')
        .findByText(`user \`${userNotExist}\` not found`)
        .should('be.visible')
      cy.visit(
        `${Cypress.config('consoleRootPath')}/applications/${applicationId}/collaborators/add`,
      )
    })

    it('succeeds adding organization as a collaborator', () => {
      cy.findByLabelText('Collaborator ID').type(organizationId)
      cy.findByLabelText('Organization').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification').should('not.exist')
      cy.location('pathname').should(
        'eq',
        `${Cypress.config('consoleRootPath')}/applications/${applicationId}/collaborators`,
      )
    })

    it('fails adding non-existent organization', () => {
      cy.findByLabelText('Collaborator ID').type(orgNotExist)
      cy.findByLabelText('Organization').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification')
        .should('be.visible')
        .findByText(`organization \`${orgNotExist}\` not found`)
        .should('be.visible')
      cy.visit(
        `${Cypress.config('consoleRootPath')}/applications/${applicationId}/collaborators/add`,
      )
    })
  })

  describe('Gateway', () => {
    const gatewayId = 'collab-test-gtw'
    const gateway = { ids: { gateway_id: gatewayId } }

    before(() => {
      cy.loginConsole({ user_id: userId, password: user.password })
      cy.createGateway(gateway, userId)
      cy.clearLocalStorage()
      cy.clearCookies()
    })

    beforeEach(() => {
      cy.loginConsole({ user_id: user.ids.user_id, password: user.password })
      cy.visit(`${Cypress.config('consoleRootPath')}/gateways/${gatewayId}/collaborators/add`)
    })

    it('succeeds adding user as a collaborator', () => {
      cy.findByLabelText('Collaborator ID').type(collaboratorId)
      cy.findByLabelText('User').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification').should('not.exist')
      cy.location('pathname').should(
        'eq',
        `${Cypress.config('consoleRootPath')}/gateways/${gatewayId}/collaborators`,
      )
    })

    it('fails adding non-existent user', () => {
      cy.findByLabelText('Collaborator ID').type(userNotExist)
      cy.findByLabelText('User').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification')
        .should('be.visible')
        .findByText(`user \`${userNotExist}\` not found`)
        .should('be.visible')
      cy.visit(`${Cypress.config('consoleRootPath')}/gateways/${gatewayId}/collaborators/add`)
    })

    it('succeeds adding organization as a collaborator', () => {
      cy.findByLabelText('Collaborator ID').type(organizationId)
      cy.findByLabelText('Organization').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification').should('not.exist')
      cy.location('pathname').should(
        'eq',
        `${Cypress.config('consoleRootPath')}/gateways/${gatewayId}/collaborators`,
      )
    })

    it('fails adding non-existent organization', () => {
      cy.findByLabelText('Collaborator ID').type(orgNotExist)
      cy.findByLabelText('Organization').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification')
        .should('be.visible')
        .findByText(`organization \`${orgNotExist}\` not found`)
        .should('be.visible')
      cy.visit(`${Cypress.config('consoleRootPath')}/gateways/${gatewayId}/collaborators/add`)
    })
  })

  describe('Organization', () => {
    const testOrgId = 'collab-test-org-2'
    const testOrg = {
      ids: { organization_id: testOrgId },
    }

    before(() => {
      cy.loginConsole({ user_id: userId, password: user.password })
      cy.createOrganization(testOrg, userId)
      cy.clearLocalStorage()
      cy.clearCookies()
    })

    beforeEach(() => {
      cy.loginConsole({ user_id: user.ids.user_id, password: user.password })
      cy.visit(`${Cypress.config('consoleRootPath')}/organizations/${testOrgId}/collaborators/add`)
    })

    it('succeeds adding user as a collaborator', () => {
      cy.findByLabelText('Collaborator ID').type(collaboratorId)
      cy.findByLabelText('User').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification').should('not.exist')
      cy.location('pathname').should(
        'eq',
        `${Cypress.config('consoleRootPath')}/organizations/${testOrgId}/collaborators`,
      )
    })

    it('fails adding non-existent user', () => {
      cy.findByLabelText('Collaborator ID').type(userNotExist)
      cy.findByLabelText('User').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification')
        .should('be.visible')
        .findByText(`user \`${userNotExist}\` not found`)
        .should('be.visible')
      cy.visit(`${Cypress.config('consoleRootPath')}/organizations/${testOrgId}/collaborators/add`)
    })

    it('fails adding organization as a collaborator', () => {
      cy.findByLabelText('Collaborator ID').type(organizationId)
      cy.findByLabelText('Organization').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification')
        .should('be.visible')
        .findByText('account of type `organization` can not collaborate on `organization`')
        .should('be.visible')
      cy.visit(`${Cypress.config('consoleRootPath')}/organizations/${testOrgId}/collaborators/add`)
    })

    it('fails adding non-existent organization', () => {
      cy.findByLabelText('Collaborator ID').type(orgNotExist)
      cy.findByLabelText('Organization').check()
      cy.findByLabelText('Grant all current and future rights').check()
      cy.findByRole('button', { name: 'Add collaborator' }).click()

      cy.findByTestId('error-notification')
        .should('be.visible')
        .findByText(`organization \`${orgNotExist}\` not found`)
        .should('be.visible')
      cy.visit(`${Cypress.config('consoleRootPath')}/organizations/${testOrgId}/collaborators/add`)
    })
  })
})
