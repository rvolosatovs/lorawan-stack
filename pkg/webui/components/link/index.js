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

import React from 'react'
import { Link as RouterLink } from 'react-router-dom'
import classnames from 'classnames'
import { defineMessages, useIntl } from 'react-intl'

import Icon from '@ttn-lw/components/icon'

import { withEnv } from '@ttn-lw/lib/components/env'
import Message from '@ttn-lw/lib/components/message'

import PropTypes from '@ttn-lw/lib/prop-types'
import { url as urlPattern } from '@ttn-lw/lib/regexp'

import style from './link.styl'

const m = defineMessages({
  glossaryTitle: 'See "{term}" in the glossary',
})

const formatTitle = function(content, values, formatter) {
  if (typeof content === 'object' && content.id && content.defaultMessage) {
    return formatter(content, values)
  }

  return content
}

const Link = function(props) {
  const {
    className,
    disabled,
    title,
    titleValues,
    id,
    children,
    to,
    replace,
    target,
    showVisited,
    onClick,
    secondary,
    primary,
  } = props

  const { formatMessage } = useIntl()
  const formattedTitle = formatTitle(title, titleValues, formatMessage)
  const classNames = classnames(style.link, className, {
    [style.linkVisited]: showVisited,
    [style.primary]: primary,
    [style.secondary]: secondary,
  })

  if (disabled) {
    return <span className={classnames(classNames, style.disabled)}>{children}</span>
  }

  return (
    <RouterLink
      className={classNames}
      id={id}
      title={formattedTitle}
      replace={replace}
      to={to}
      target={target}
      onClick={onClick}
    >
      {children}
    </RouterLink>
  )
}

Link.propTypes = {
  children: PropTypes.node,
  className: PropTypes.string,
  disabled: PropTypes.bool,
  id: PropTypes.string,
  onClick: PropTypes.func,
  primary: PropTypes.bool,
  replace: PropTypes.bool,
  secondary: PropTypes.bool,
  showVisited: PropTypes.bool,
  target: PropTypes.string,
  title: PropTypes.message,
  titleValues: PropTypes.shape({}),
  to: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.shape({
      pathname: PropTypes.string,
      search: PropTypes.string,
      hash: PropTypes.string,
      state: PropTypes.object,
    }),
  ]).isRequired,
}

Link.defaultProps = {
  children: undefined,
  className: undefined,
  disabled: false,
  id: undefined,
  onClick: () => null,
  primary: false,
  showVisited: false,
  replace: false,
  secondary: false,
  target: undefined,
  title: undefined,
  titleValues: undefined,
}

const DocLink = function(props) {
  const {
    className,
    name,
    title,
    titleValues,
    id,
    path,
    children,
    showVisited,
    secondary,
    primary,
    disabled,
    to,
    raw,
    onClick,
    env: {
      config: { documentationBaseUrl },
    },
  } = props

  const { formatMessage } = useIntl()
  const classNames = classnames(style.link, className, {
    [style.linkVisited]: showVisited,
    [style.primary]: primary,
    [style.secondary]: secondary,
  })
  if (disabled) {
    return <span className={classnames(classNames, style.disabled)}>{children}</span>
  }
  const link = documentationBaseUrl.concat(path)
  const formattedTitle = formatTitle(title, titleValues, formatMessage)

  return (
    <a
      className={classNames}
      to={to}
      title={formattedTitle}
      id={id}
      href={link}
      target="blank"
      name={name}
      onClick={onClick}
    >
      {!raw && <Icon className={style.docIcon} icon="book" />}
      {children}
      {!raw && <Icon className={style.icon} icon="launch" />}
    </a>
  )
}

DocLink.propTypes = {
  ...Link.propTypes,
  env: PropTypes.env.isRequired,
  name: PropTypes.string,
  path: PropTypes.string.isRequired,
  raw: PropTypes.bool,
  showVisited: PropTypes.bool,
  to: PropTypes.oneOfType([
    PropTypes.string,
    PropTypes.shape({
      pathname: PropTypes.string,
      search: PropTypes.string,
      hash: PropTypes.string,
      state: PropTypes.object,
    }),
  ]),
}

DocLink.defaultProps = {
  ...Link.defaultProps,
  disabled: false,
  name: undefined,
  showVisited: false,
  to: undefined,
  raw: false,
}

Link.DocLink = withEnv(DocLink)

const GlossaryLink = ({ glossaryId, term, hideTerm, primary, secondary, className }) => {
  const { formatMessage } = useIntl()
  return (
    <Link.DocLink
      primary={primary}
      secondary={secondary}
      className={className}
      path={`/reference/glossary#${glossaryId}`}
      title={m.glossaryTitle}
      titleValues={{ term: formatTitle(term, undefined, formatMessage) }}
      raw
    >
      {!hideTerm && <Message content={term} />}{' '}
      <Icon className={style.glossaryIcon} icon="help_outline" />
    </Link.DocLink>
  )
}

GlossaryLink.propTypes = {
  className: PropTypes.string,
  glossaryId: PropTypes.string.isRequired,
  hideTerm: PropTypes.bool,
  primary: PropTypes.bool,
  secondary: PropTypes.bool,
  term: PropTypes.message.isRequired,
}

GlossaryLink.defaultProps = {
  className: '',
  hideTerm: false,
  primary: false,
  secondary: false,
}

Link.GlossaryLink = GlossaryLink

const AnchorLink = function(props) {
  const {
    className,
    name,
    title,
    titleValues,
    id,
    href,
    target,
    children,
    showVisited,
    secondary,
    primary,
    disabled,
    external,
  } = props

  const { formatMessage } = useIntl()
  const formattedTitle = formatTitle(title, titleValues, formatMessage)
  const classNames = classnames(style.link, className, {
    [style.linkVisited]: showVisited,
    [style.primary]: primary,
    [style.secondary]: secondary,
  })

  if (disabled) {
    return <span className={classnames(classNames, style.disabled)}>{children}</span>
  }

  return (
    <a
      className={classNames}
      title={formattedTitle}
      id={id}
      href={href}
      target={external ? 'blank' : target}
      name={name}
    >
      {children}
      {external ? <Icon className={style.icon} icon="launch" /> : null}
    </a>
  )
}

AnchorLink.propTypes = {
  href: PropTypes.string.isRequired,
  id: PropTypes.string,
  name: PropTypes.string,
  showVisited: PropTypes.bool,
  target: PropTypes.string,
  title: PropTypes.string,
}

AnchorLink.defaultProps = {
  showVisited: false,
}

Link.Anchor = AnchorLink

const BaseAnchorLink = function({ env, href, ...rest }) {
  const { appRoot } = env

  // Prevent prefixing proper URLs.
  const path = href.match(urlPattern) ? href : appRoot + href

  return <Link.Anchor href={path} {...rest} />
}

BaseAnchorLink.propTypes = AnchorLink.propTypes
BaseAnchorLink.defaultProps = AnchorLink.defaultProps

Link.BaseAnchor = withEnv(BaseAnchorLink)

export default Link
