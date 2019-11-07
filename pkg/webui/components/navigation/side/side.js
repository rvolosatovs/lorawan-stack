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

import React, { Component } from 'react'
import bind from 'autobind-decorator'
import classnames from 'classnames'
import { withRouter } from 'react-router-dom'
import { defineMessages } from 'react-intl'
import PropTypes from '../../../lib/prop-types'

import Button from '../../button'
import Icon from '../../icon'
import Message from '../../../lib/components/message'
import SideNavigationList from './list'

import style from './side.styl'

const m = defineMessages({
  hideSidebar: 'Hide Sidebar',
})
@withRouter
@bind
class SideNavigation extends Component {
  state = {
    /** A flag specifying whether the side navigation is minimized or not */
    isMinimized: false,
    /**
     * A map of expanded items, where:
     *  - The key: index of the item
     *  - The value: an object consisting of:
     *    - isOpen - boolean flag specifying whether the item is opened or not
     *    - isLink - boolean flag specifying whether a link is selected within
     *                this opened item
     */
    itemsExpanded: {},
  }

  onToggle() {
    this.setState(function(prev) {
      return { isMinimized: !prev.isMinimized }
    })
  }

  onItemExpand(index, linkSelected) {
    this.setState(function(prev) {
      const oldItemsExpanded = prev.itemsExpanded
      const oldMinimized = prev.isMinimized

      // make sure that no more links are active
      if (linkSelected) {
        const itemsExpanded = Object.keys(oldItemsExpanded)
          .map(idx => +idx)
          .reduce(function(acc, idx) {
            const { isOpen, isLink } = oldItemsExpanded[idx] || {}
            if (index === idx) {
              acc[idx] = { isOpen: true, isLink: true }
            } else if (isLink) {
              acc[idx] = { isOpen, isLink: false }
            } else {
              acc[idx] = oldItemsExpanded[idx]
            }

            return acc
          }, {})

        return { itemsExpanded }
      }

      const { isOpen = false, isLink = false } = oldItemsExpanded[index] || {}

      const shouldOpen = oldMinimized || !isOpen
      const shouldLink = isLink || linkSelected
      const shouldMinimize = oldMinimized && linkSelected

      return {
        itemsExpanded: {
          ...oldItemsExpanded,
          [index]: { isOpen: shouldOpen, isLink: shouldLink },
        },
        isMinimized: shouldMinimize,
      }
    })
  }

  componentDidMount() {
    const { location, entries } = this.props
    for (let i in entries) {
      if (entries[i].path) {
        if (location.pathname === entries[i].path) break
      } else if (entries[i].nested) {
        for (let j in entries[i].items) {
          if (location.pathname === entries[i].items[j].path) {
            const itemsExpanded = { [i]: { isOpen: true, isLink: true } }
            this.setState({ itemsExpanded })
            return
          }
        }
      }
    }
  }

  render() {
    const { className, header, entries } = this.props
    const { isMinimized, itemsExpanded } = this.state

    const navigationClassNames = classnames(className, style.navigation, {
      [style.navigationMinimized]: isMinimized,
    })
    const headerClassNames = classnames(style.header, {
      [style.headerMinimized]: isMinimized,
    })

    return (
      <nav className={navigationClassNames}>
        <div>
          <div className={headerClassNames}>
            <Icon className={style.icon} icon={header.icon} />
            <Message className={style.message} content={header.title} />
          </div>
          <SideNavigationList
            itemsExpanded={itemsExpanded}
            onItemExpand={this.onItemExpand}
            items={entries}
            isMinimized={isMinimized}
          />
        </div>
        <Button
          className={style.navigationButton}
          naked
          secondary
          icon={isMinimized ? 'keyboard_arrow_right' : 'keyboard_arrow_left'}
          message={isMinimized ? null : m.hideSidebar}
          onClick={this.onToggle}
          data-hook="side-nav-hide-button"
        />
      </nav>
    )
  }
}

SideNavigation.propTypes = {
  /**
   * A list of entry objects for the side navigation
   * See `<SideNavigationList/>`'s `items` proptypes for details
   */
  entries: SideNavigationList.propTypes.items,
  /** The header for the side navigation */
  header: PropTypes.shape({
    title: PropTypes.string,
    icon: PropTypes.string,
  }),
}

export default SideNavigation
