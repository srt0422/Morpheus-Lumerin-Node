import React, { useState } from 'react';
import { List as RVList, AutoSizer, InfiniteLoader } from 'react-virtualized';
import styled from 'styled-components';

import withTxListState from '../../../store/hocs/withTxListState';
import ScanningTxPlaceholder from './ScanningTxPlaceholder';
import NoTxPlaceholder from './NoTxPlaceholder';
import { ItemFilter, Flex } from '../../common';
import Header from './Header';
import TxRow from './row/Row';
import Spinner from '../../common/Spinner';

const Container = styled.div`
  margin-top: 2.4rem;
  height: 100%;

  @media (min-width: 960px) {
  }
`;

const LoadingRov = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  color: ${p => p.theme.colors.primary};
`;

const Transactions = styled.div`
  margin: 1.6rem 0 1.6rem;
  border-radius: 15px;
  background-color: transparent;
`;

const ListContainer = styled.div`
  height: calc(100vh - 490px);
  border-radius: 0.375rem;

  background: rgba(255,255,255, 0.04);
  border-width: 1px;
  border: 1px solid rgba(255, 255, 255, 0.04);
  color: white;
`;

const TxRowContainer = styled.div`
  &:hover {
    background-color: rgba(0,0,0,0.5);
  }
`;

const Title = styled.div`
  font-size: 2.4rem;
  line-height: 3rem;
  white-space: nowrap;
  margin: 0;
  font-weight: 500;
  color: ${p => p.theme.colors.morMain};
  margin-bottom: 4.8px;
  margin-right: 2.4rem;
  cursor: default;

  @media (min-width: 1140px) {
    margin-right: 0.8rem;
  }

  @media (min-width: 1200px) {
    margin-right: 1.6rem;
  }
`;

export const TxList = ({
  transactions,
  loadNextTransactions,
  hasTransactions,
  syncStatus,
  hasNextPage,
  ...props
}) => {

  const [nextPageLoading, setNextPageLoading] = useState(false);

  const load = async () => {
    setNextPageLoading(true);
    await loadNextTransactions();
    setNextPageLoading(false);
  }
  return (
    <Container data-testid="tx-list">
      <Flex.Row grow="1">
        <Title>Last Transactions</Title>
      </Flex.Row>
      <Transactions>
        <ItemFilter
          extractValue={({ txType }) => txType}
          items={transactions.filter(({ txType }) => txType)}
        >
          {({ filteredItems, onFilterChange, activeFilter }) => {
            // const rowCount = hasNextPage
            //   ? filteredItems.length + 1
            //   : filteredItems.length;
            const rowCount = filteredItems.length;

            // Only load 1 page of items at a time.
            // Pass an empty callback to InfiniteLoader in case it asks us to load more than once.
            const loadMoreRows = nextPageLoading
              ? () => {}
              : () => {
                // load();
                };

            // Every row is loaded except for our loading indicator row.
            const isRowLoaded = ({ index }) => {
              return !transactions.length || index < rowCount - 1;
            };

            const rowRenderer = ({ key, style, index }) =>
              isRowLoaded({ index }) ? (
                <TxRowContainer
                  style={style}
                  key={`${key}-${filteredItems[index].hash}`}
                >
                  <TxRow
                    data-testid="tx-row"
                    data-hash={filteredItems[index].hash}
                    tx={filteredItems[index]}
                  />
                </TxRowContainer>
              ) : (
                null
                // <LoadingRov key={key} style={style}>
                //   Loading... <Spinner></Spinner>
                // </LoadingRov>
              );

            return (
              <React.Fragment>
                <Header
                  onWalletRefresh={() => {}}
                  hasTransactions={hasTransactions}
                  onFilterChange={onFilterChange}
                  activeFilter={activeFilter}
                  syncStatus={syncStatus}
                />

                <ListContainer>
                  {!transactions.length &&
                    (syncStatus === 'syncing' ? (
                      <ScanningTxPlaceholder />
                    ) : (
                      <NoTxPlaceholder />
                    ))}
                  {+transactions.length > 0 && (
                    <AutoSizer>
                          {({ width, height }) => (
                            <RVList
                              rowRenderer={rowRenderer}
                              rowHeight={66}
                              rowCount={rowCount}
                              height={height || 500} // defaults for tests
                              width={width || 500} // defaults for tests
                            />
                          )}
                        </AutoSizer>
                    // <InfiniteLoader
                    //   isRowLoaded={isRowLoaded}
                    //   loadMoreRows={loadMoreRows}
                    //   rowCount={rowCount}
                    //   threshold={activeFilter ? 1 : 10}
                    // >
                    //   {({ onRowsRendered, registerChild }) => (
                    //     <AutoSizer>
                    //       {({ width, height }) => (
                    //         <RVList
                    //           ref={registerChild}
                    //           onRowsRendered={onRowsRendered}
                    //           rowRenderer={rowRenderer}
                    //           rowHeight={66}
                    //           rowCount={rowCount}
                    //           height={height || 500} // defaults for tests
                    //           width={width || 500} // defaults for tests
                    //         />
                    //       )}
                    //     </AutoSizer>
                    //   )}
                    // </InfiniteLoader>
                  )}
                </ListContainer>
              </React.Fragment>
            );
          }}
        </ItemFilter>
      </Transactions>
    </Container>
  );
};

export default withTxListState(TxList);
