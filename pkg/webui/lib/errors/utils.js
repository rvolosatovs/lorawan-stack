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

import * as Sentry from '@sentry/browser'
import { isPlainObject } from 'lodash'

import { error as errorLog } from '@ttn-lw/lib/log'

import errorMessages from './error-messages'
import grpcErrToHttpErr from './grpc-error-map'

/**
 * Tests whether the error is a backend error object.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` is a well known backend error object.
 */
export const isBackend = error =>
  Boolean(error) &&
  typeof error === 'object' &&
  !('id' in error) &&
  error.message &&
  error.details &&
  (error.code || error.grpc_code)

/**
 * Returns whether the error is a frontend defined error object.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` is a well known frontend error object.
 */
export const isFrontend = error => Boolean(error) && typeof error === 'object' && error.isFrontend

/**
 * Returns whether `details` is a backend error details object.
 *
 * @param {object} details - The object to be tested.
 * @returns {boolean} `true` if `details` is a well known backend error details object,
 * `false` otherwise.
 */
export const isBackendErrorDetails = details =>
  Boolean(details) &&
  Boolean(details.namespace) &&
  Boolean(details.name) &&
  Boolean(details.message_format) &&
  Boolean(details.code)

/**
 * Returns whether the error has a shape that is not well-known.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` is not of a well known shape.
 */
export const isUnknown = error => !isBackend(error) && !isFrontend(error)

/**
 * Returns a frontend error object, to be passed to error components.
 *
 * @param {object} errorTitle - The error message title (i18n message).
 * @param {object} errorMessage - The error message object (i18n message).
 * @param {string} errorCode - An optional error code to be used to identify
 * a specific error type easily. E.g. `user_status_unapproved`.
 * @param {number} statusCode - An optional status code corresponding to
 * the well known HTTP status codes. This can help categorizing the error if
 * necessary.
 * @returns {object} A frontend error object to be passed to error components.
 */
export const createFrontendError = (errorTitle, errorMessage, errorCode, statusCode) => ({
  errorTitle,
  errorMessage,
  errorCode,
  statusCode,
  isFrontend: true,
})

/**
 * Maps the error type to a HTTP Status code. Useful for quickly
 * determining the type of the error. Returns false if no status code can be
 * determined.
 *
 * @param {object} error - The error to be tested.
 * @returns {number} The (closest when grpc error) HTTP Status Code, otherwise
 * `undefined`.
 */
export const httpStatusCode = error => {
  if (!Boolean(error)) {
    return undefined
  }

  let statusCode = undefined
  if (isBackend(error)) {
    statusCode = error.http_code || grpcErrToHttpErr(error.code || error.grpc_code)
  } else if (isFrontend(error)) {
    statusCode = error.statusCode
  } else if (Boolean(error.statusCode)) {
    statusCode = error.statusCode
  } else if (Boolean(error.response) && Boolean(error.response.status)) {
    statusCode = error.response.status
  }

  return Boolean(statusCode) ? parseInt(statusCode) : undefined
}
/**
 * Returns the GRPC Status code in case of a backend error.
 *
 * @param {object} error - The error to be tested.
 * @returns {number} The GRPC error code, or `false`.
 */
export const grpcStatusCode = error => isBackend(error) && (error.code || error.grpc_code)

/**
 * Tests whether the grpc error represents the not found erorr.
 *
 * @param {object} error - The error object to be tested.
 * @returns {boolean} `true` if `error` represents the not found error,
 * `false` otherwise.
 */
export const isNotFoundError = error => grpcStatusCode(error) === 5 || httpStatusCode(error) === 404

/**
 * Returns whether the grpc error represents an internal server error.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` represents an internal server error,
 * `false` otherwise.
 */
export const isInternalError = error => grpcStatusCode(error) === 13 // NOTE: HTTP 500 can also be UnknownError.

/**
 * Returns whether the grpc error represents an invalid argument or bad request
 * error.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` represents an invalid argument or bad
 * request error, `false` otherwise.
 */
export const isInvalidArgumentError = error =>
  grpcStatusCode(error) === 3 || httpStatusCode(error) === 400

/**
 * Returns whether the grpc error represents an already exists error.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` represents an already exists error,
 * `false` otherwise.
 */
export const isAlreadyExistsError = error => grpcStatusCode(error) === 6 // NOTE: HTTP 409 can also be AbortedError.

/**
 * Returns whether the grpc error represents a permission denied error.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` represents a permission denied error,
 * `false` otherwise.
 */
export const isPermissionDeniedError = error =>
  grpcStatusCode(error) === 7 || httpStatusCode(error) === 403

/**
 * Returns whether the grpc error represents an error due to not being
 * authenticated.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` represents an `Unauthenticated` error,
 * `false` otherwise.
 */
export const isUnauthenticatedError = error =>
  grpcStatusCode(error) === 16 || httpStatusCode(error) === 401

/**
 * Returns whether the grpc error represents a conflict with the current state on the server.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` represents a `Conflict` error, `false` otherwise.
 */
export const isConflictError = error =>
  grpcStatusCode(error) === 10 || httpStatusCode(error) === 409

/**
 * Returns whether `error` has translation ids.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` has translation ids, `false` otherwise.
 */
export const isTranslated = error =>
  isBackend(error) || isFrontend(error) || (typeof error === 'object' && error.id)

/**
 * Returns whether `error` is a 'network error' as JavaScript TypeError.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` is a network error, `false` otherwise.
 */
export const isNetworkError = error =>
  error instanceof Error && error.message.toLowerCase() === 'network error'

/**
 * Returns whether `error` is a 'ECONNABORTED' error as returned from axios.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` is a timeout error, `false` otherwise.
 */
export const isTimeoutError = error => error.code === 'ECONNABORTED'

/**
 * Returns whether the error is worth being sent to Sentry.
 *
 * @param {object} error - The error to be tested.
 * @returns {boolean} `true` if `error` should be forwarded to Sentry,
 * `false` otherwise.
 */
export const isSentryWorthy = error =>
  (isUnknown(error) && httpStatusCode(error) === undefined) ||
  isInvalidArgumentError(error) ||
  isInternalError(error) ||
  httpStatusCode(error) >= 500 || // Server errors.
  httpStatusCode(error) === 400 // Bad request.

/**
 * Returns an appropriate error title that can be used for Sentry.
 *
 * @param {object} error - The error object.
 * @returns {string} The Sentry error title.
 */
export const getSentryErrorTitle = error => {
  if (typeof error !== 'object') {
    return `invalid error type: ${error}`
  }

  if (isBackend(error)) {
    return error.message
  } else if (isFrontend(error)) {
    return error.errorTitle.defaultMessage
  } else if ('message' in error) {
    return error.message
  } else if ('code' in error) {
    return error.code
  } else if ('statusCode' in error) {
    return `status code: ${error.statusCode}`
  }

  return 'untitled or empty error'
}

/**
 * Returns the id of the error, used as message id.
 *
 * @param {object} error - The backend error object.
 * @returns {string} The ID.
 */
export const getBackendErrorId = error => error.message.split(' ')[0]

/**
 * Returns the id of the error details, used as message id.
 *
 * @param {object} details - The backend error details object.
 * @returns {string} The ID.
 */
export const getBackendErrorDetailsId = details => `error:${details.namespace}:${details.name}`

/**
 * Returns error details.
 *
 * @param {object} error - The backend error object.
 * @returns {object} - The details of `error`.
 */
export const getBackendErrorDetails = error => error.details[0]

/**
 * Returns the name of the error extracted from the details array.
 *
 * @param {object} error - The backend error object.
 * @returns {string} - The error name.
 */
export const getBackendErrorName = error =>
  error && error.details instanceof Array && error.details[0] && error.details[0].name
    ? error.details[0].name
    : undefined
/**
 * Returns the default message of the error, used as fallback translation.
 *
 * @param {object} error - The backend error object.
 * @returns {string} The default message.
 */
export const getBackendErrorDefaultMessage = error =>
  error.details[0].message_format || error.message.replace(/^.*\s/, '')

/**
 * Returns whether the error has one or more cause properties.
 *
 * @param {object} error - The backend error object.
 * @returns {boolean} - Whether the error has one or more cuase properties.
 */
export const hasCauses = error => isBackendErrorDetails(error) && 'cause' in error

/**
 * Returns the root cause of the backend error.
 *
 * @param {object} error - The backend error object.
 * @returns {object} - The root cause of `error`.
 */
export const getBackendErrorRootCause = error => {
  let rootCause
  if (hasCauses(error)) {
    rootCause = error.cause
  } else {
    rootCause = getBackendErrorDetails(error)
  }

  while ('cause' in rootCause) {
    rootCause = rootCause.cause
  }

  return rootCause
}

/**
 * Returns the attributes of the backend error message, if any.
 *
 * @param {object} error - The backend error object.
 * @returns {string} The attributes or undefined.
 */
export const getBackendErrorMessageAttributes = error => error.details[0].attributes

/**
 * Adapts the error object to props of message object, if possible.
 *
 * @param {object} error - The backend error object.
 * @returns {object} Message props of the error object, or generic error object.
 */
export const toMessageProps = error => {
  let props
  // Check if it is a error message and transform it to a intl message.
  if (isBackend(error)) {
    props = {
      content: {
        id: getBackendErrorId(error),
        defaultMessage: getBackendErrorDefaultMessage(error),
      },
      values: getBackendErrorMessageAttributes(error),
    }
  } else if (isBackendErrorDetails(error)) {
    props = {
      content: {
        id: getBackendErrorDetailsId(error),
        defaultMessage: error.message_format,
      },
      values: error.attributes,
    }
  } else if (isFrontend(error)) {
    props = {
      content: error.errorMessage,
      title: error.errorTitle,
    }
  } else if (isTranslated(error)) {
    // Fall back to normal message.
    props = { content: error }
  } else {
    // Fall back to generic error message.
    props = { content: errorMessages.genericError }
  }

  return props
}

/**
 * `ingestError` provides a unified error ingestion handler, which manages
 * forwarding to Sentry and other logic that should be applied when errors
 * occur. The error object is not modified.
 *
 * @param {object} error - The error object.
 * @param {object} extras - Sentry extras to be sent.
 * @param {object} tags - Sentry tags to be sent.
 */
export const ingestError = (error, extras = {}, tags = {}) => {
  // Log the error when in development mode
  errorLog(error)

  // Send to Sentry if necessary.
  if (isSentryWorthy(error)) {
    Sentry.withScope(scope => {
      scope.setTags({ ...tags, frontendOrigin: true })
      scope.setFingerprint(isBackend(error) ? getBackendErrorId(error) : error)
      if (isPlainObject(error)) {
        scope.setExtras({ ...error, ...extras })
      } else {
        scope.setExtras({ error, ...extras })
      }
      Sentry.captureException(
        error instanceof Error ? error : new Error(getSentryErrorTitle(error)),
      )
    })
  }
}
