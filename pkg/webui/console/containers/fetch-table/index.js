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
import { connect } from 'react-redux'
import { push } from 'connected-react-router'
import bind from 'autobind-decorator'
import classnames from 'classnames'

import debounce from '../../../lib/debounce'

import sharedMessages from '../../../lib/shared-messages'
import Tabular from '../../../components/table'
import Input from '../../../components/input'
import Button from '../../../components/button'
import Tabs from '../../../components/tabs'

import style from './fetch-table.styl'

const DEFAULT_PAGE = 1
const DEFAULT_TAB = 'all'
const ALLOWED_TABS = ['all']
const ALLOWED_ORDERS = ['asc', 'desc', undefined]

const filterValidator = function(filters) {
  if (!ALLOWED_TABS.includes(filters.tab)) {
    filters.tab = DEFAULT_TAB
  }

  if (!ALLOWED_ORDERS.includes(filters.order)) {
    filters.order = undefined
    filters.orderBy = undefined
  }

  if (
    (Boolean(filters.order) && !Boolean(filters.orderBy)) ||
    (!Boolean(filters.order) && Boolean(filters.orderBy))
  ) {
    filters.order = undefined
    filters.orderBy = undefined
  }

  if (!Boolean(filters.page) || filters.page < 0) {
    filters.page = DEFAULT_PAGE
  }

  return filters
}

@connect(function(state, props) {
  const base = props.baseDataSelector(state, props)

  return {
    items: base[props.entity] || [],
    totalCount: base.totalCount || 0,
    fetching: base.fetching,
    fetchingSearch: base.fetchingSearch,
    pathname: state.router.location.pathname,
    mayAdd: 'mayAdd' in base ? base.mayAdd : true,
  }
})
@bind
class FetchTable extends Component {
  constructor(props) {
    super(props)

    this.state = {
      query: '',
      page: 1,
      tab: 'all',
      order: undefined,
      orderBy: undefined,
    }

    const { debouncedFunction, cancel } = debounce(this.requestSearch, 350)

    this.debouncedRequestSearch = debouncedFunction
    this.debounceCancel = cancel
  }

  componentDidMount() {
    this.fetchItems()
  }

  componentWillUnmount() {
    this.debounceCancel()
  }

  fetchItems() {
    const { dispatch, pageSize, searchItemsAction, getItemsAction } = this.props

    const filters = { ...this.state, limit: pageSize }

    if (filters.query) {
      dispatch(searchItemsAction(filters))
    } else {
      dispatch(getItemsAction(filters))
    }
  }

  async onPageChange(page) {
    await this.setState(
      this.props.filterValidator({
        ...this.state,
        page,
      }),
    )

    this.fetchItems()
  }

  async requestSearch() {
    await this.setState(
      this.props.filterValidator({
        ...this.state,
        page: 1,
      }),
    )

    this.fetchItems()
  }

  async onQueryChange(query) {
    await this.setState(
      this.props.filterValidator({
        ...this.state,
        query,
      }),
    )

    this.debouncedRequestSearch()
  }

  async onOrderChange(order, orderBy) {
    await this.setState(
      this.props.filterValidator({
        ...this.state,
        order,
        orderBy,
      }),
    )

    this.fetchItems()
  }

  async onTabChange(tab) {
    await this.setState(
      this.props.filterValidator({
        ...this.state,
        tab,
      }),
    )
    this.fetchItems()
  }

  onItemClick(index) {
    const {
      dispatch,
      pathname,
      items,
      entity,
      itemPathPrefix,
      getItemPathPrefix,
      handlesPagination,
      pageSize,
    } = this.props
    const { page } = this.state

    let itemIndex = index
    if (handlesPagination) {
      const pageNr = page - 1 // switch to 0-based pagination
      itemIndex += pageSize * pageNr
    }

    const entitySingle = entity.substr(0, entity.length - 1)
    let entityPath
    if (Boolean(getItemPathPrefix)) {
      entityPath = getItemPathPrefix(items[itemIndex])
    } else {
      const item_id = items[itemIndex].id || items[itemIndex].ids[`${entitySingle}_id`]
      entityPath = `${itemPathPrefix}/${item_id}`
    }

    dispatch(push(`${pathname}${entityPath}`))
  }

  render() {
    const {
      items,
      totalCount,
      fetching,
      fetchingSearch,
      mayAdd,
      pageSize,
      addMessage,
      tableTitle,
      headers,
      tabs,
      searchable,
      handlesPagination,
      itemPathPrefix,
      pathname,
      actionItems,
    } = this.props
    const { page, query, tab } = this.state

    const buttonClassNames = classnames(style.filters, {
      [style.topRule]: Boolean(tabs || tableTitle),
    })

    return (
      <div>
        <div className={buttonClassNames}>
          <div className={style.filtersLeft}>
            {tabs && (
              <Tabs
                active={tab}
                className={style.tabs}
                tabs={tabs}
                onTabChange={this.onTabChange}
              />
            )}
            {tableTitle && (
              <div className={style.tableTitle}>
                {tableTitle} ({totalCount})
              </div>
            )}
          </div>
          <div className={style.filtersRight}>
            {searchable && (
              <Input
                value={query}
                icon="search"
                loading={fetchingSearch}
                onChange={this.onQueryChange}
              />
            )}
            {actionItems}
            {mayAdd && (
              <Button.Link
                className={style.addButton}
                message={addMessage}
                icon="add"
                to={`${pathname}${itemPathPrefix}/add`}
              />
            )}
          </div>
        </div>
        <Tabular
          paginated
          page={page}
          totalCount={totalCount}
          pageSize={pageSize}
          onRowClick={this.onItemClick}
          onPageChange={this.onPageChange}
          loading={fetching}
          headers={headers}
          data={items}
          emptyMessage={sharedMessages.noMatch}
          handlesPagination={handlesPagination}
        />
      </div>
    )
  }
}

FetchTable.defaultProps = {
  pageSize: 20,
  filterValidator,
  itemPathPrefix: '',
}

export default FetchTable
