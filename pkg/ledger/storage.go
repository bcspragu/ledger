package ledger

import (
	"context"
	"time"

	"github.com/formancehq/go-libs/api"
	"github.com/numary/ledger/pkg/core"
)

type Store interface {
	GetLastTransaction(ctx context.Context) (*core.ExpandedTransaction, error)
	CountTransactions(context.Context, TransactionsQuery) (uint64, error)
	GetTransactions(context.Context, TransactionsQuery) (api.Cursor[core.ExpandedTransaction], error)
	GetTransaction(ctx context.Context, txid uint64) (*core.ExpandedTransaction, error)
	GetAccount(ctx context.Context, accountAddress string) (*core.Account, error)
	GetAssetsVolumes(ctx context.Context, accountAddress string) (core.AssetsVolumes, error)
	GetAccountWithVolumes(ctx context.Context, account string) (*core.AccountWithVolumes, error)
	GetVolumes(ctx context.Context, accountAddress, asset string) (core.Volumes, error)
	CountAccounts(context.Context, AccountsQuery) (uint64, error)
	GetAccounts(context.Context, AccountsQuery) (api.Cursor[core.Account], error)
	GetBalances(context.Context, BalancesQuery) (api.Cursor[core.AccountsBalances], error)
	GetBalancesAggregated(context.Context, BalancesQuery) (core.AssetsBalances, error)
	GetLastLog(context.Context) (*core.Log, error)
	GetLogs(context.Context, *LogsQuery) (api.Cursor[core.Log], error)
	LoadMapping(context.Context) (*core.Mapping, error)
	GetMigrationsAvailable() ([]core.MigrationInfo, error)
	GetMigrationsDone(context.Context) ([]core.MigrationInfo, error)

	UpdateTransactionMetadata(ctx context.Context, txid uint64, metadata core.Metadata, at time.Time) error
	UpdateAccountMetadata(ctx context.Context, address string, metadata core.Metadata, at time.Time) error
	Commit(ctx context.Context, txs ...core.ExpandedTransaction) error
	SaveMapping(ctx context.Context, m core.Mapping) error
	Name() string
	Initialize(context.Context) (bool, error)
	Close(context.Context) error
}

const (
	QueryDefaultPageSize = 15
)

type TransactionsQuery struct {
	PageSize  uint
	AfterTxID uint64
	Filters   TransactionsQueryFilters
}

type TransactionsQueryFilters struct {
	Reference   string
	Destination string
	Source      string
	Account     string
	EndTime     time.Time
	StartTime   time.Time
	Metadata    map[string]string
}

func NewTransactionsQuery() *TransactionsQuery {
	return &TransactionsQuery{
		PageSize: QueryDefaultPageSize,
	}
}

func (a *TransactionsQuery) WithPageSize(pageSize uint) *TransactionsQuery {
	if pageSize != 0 {
		a.PageSize = pageSize
	}

	return a
}

func (a *TransactionsQuery) WithAfterTxID(after uint64) *TransactionsQuery {
	a.AfterTxID = after

	return a
}

func (a *TransactionsQuery) WithStartTimeFilter(start time.Time) *TransactionsQuery {
	if !start.IsZero() {
		a.Filters.StartTime = start
	}

	return a
}

func (a *TransactionsQuery) WithEndTimeFilter(end time.Time) *TransactionsQuery {
	if !end.IsZero() {
		a.Filters.EndTime = end
	}

	return a
}

func (a *TransactionsQuery) WithAccountFilter(account string) *TransactionsQuery {
	a.Filters.Account = account

	return a
}

func (a *TransactionsQuery) WithDestinationFilter(dest string) *TransactionsQuery {
	a.Filters.Destination = dest

	return a
}

func (a *TransactionsQuery) WithReferenceFilter(ref string) *TransactionsQuery {
	a.Filters.Reference = ref

	return a
}

func (a *TransactionsQuery) WithSourceFilter(source string) *TransactionsQuery {
	a.Filters.Source = source

	return a
}

func (a *TransactionsQuery) WithMetadataFilter(metadata map[string]string) *TransactionsQuery {
	a.Filters.Metadata = metadata

	return a
}

type AccountsQuery struct {
	PageSize     uint
	Offset       uint
	AfterAddress string
	Filters      AccountsQueryFilters
}

type AccountsQueryFilters struct {
	Address         string
	Balance         string
	BalanceOperator BalanceOperator
	Metadata        map[string]string
}

type BalanceOperator string

const (
	BalanceOperatorE   BalanceOperator = "e"
	BalanceOperatorGt  BalanceOperator = "gt"
	BalanceOperatorGte BalanceOperator = "gte"
	BalanceOperatorLt  BalanceOperator = "lt"
	BalanceOperatorLte BalanceOperator = "lte"
	BalanceOperatorNe  BalanceOperator = "ne"

	DefaultBalanceOperator = BalanceOperatorGte
)

func (b BalanceOperator) IsValid() bool {
	switch b {
	case BalanceOperatorE,
		BalanceOperatorGt,
		BalanceOperatorGte,
		BalanceOperatorLt,
		BalanceOperatorNe,
		BalanceOperatorLte:
		return true
	}

	return false
}

func NewBalanceOperator(s string) (BalanceOperator, bool) {
	if !BalanceOperator(s).IsValid() {
		return "", false
	}

	return BalanceOperator(s), true
}

func NewAccountsQuery() *AccountsQuery {
	return &AccountsQuery{
		PageSize: QueryDefaultPageSize,
	}
}

func (a *AccountsQuery) WithPageSize(pageSize uint) *AccountsQuery {
	if pageSize != 0 {
		a.PageSize = pageSize
	}

	return a
}

func (a *AccountsQuery) WithOffset(offset uint) *AccountsQuery {
	a.Offset = offset

	return a
}

func (a *AccountsQuery) WithAfterAddress(after string) *AccountsQuery {
	a.AfterAddress = after

	return a
}

func (a *AccountsQuery) WithAddressFilter(address string) *AccountsQuery {
	a.Filters.Address = address

	return a
}

func (a *AccountsQuery) WithBalanceFilter(balance string) *AccountsQuery {
	a.Filters.Balance = balance

	return a
}

func (a *AccountsQuery) WithBalanceOperatorFilter(balanceOperator BalanceOperator) *AccountsQuery {
	a.Filters.BalanceOperator = balanceOperator

	return a
}

func (a *AccountsQuery) WithMetadataFilter(metadata map[string]string) *AccountsQuery {
	a.Filters.Metadata = metadata

	return a
}

type BalancesQuery struct {
	PageSize     uint
	Offset       uint
	AfterAddress string
	Filters      BalancesQueryFilters
}

type BalancesQueryFilters struct {
	AddressRegexp string
}

func NewBalancesQuery() *BalancesQuery {
	return &BalancesQuery{
		PageSize: QueryDefaultPageSize,
	}
}

func (b *BalancesQuery) WithAfterAddress(after string) *BalancesQuery {
	b.AfterAddress = after

	return b
}

func (b *BalancesQuery) WithOffset(offset uint) *BalancesQuery {
	b.Offset = offset

	return b
}

func (b *BalancesQuery) WithAddressFilter(address string) *BalancesQuery {
	b.Filters.AddressRegexp = address

	return b
}

func (b *BalancesQuery) WithPageSize(pageSize uint) *BalancesQuery {
	b.PageSize = pageSize
	return b
}

type LogsQuery struct {
	AfterID  uint64
	PageSize uint

	Filters LogsQueryFilters
}

type LogsQueryFilters struct {
	EndTime   time.Time
	StartTime time.Time
}

func NewLogsQuery() *LogsQuery {
	return &LogsQuery{
		PageSize: QueryDefaultPageSize,
	}
}

func (l *LogsQuery) WithAfterID(after uint64) *LogsQuery {
	l.AfterID = after

	return l
}

func (l *LogsQuery) WithPageSize(pageSize uint) *LogsQuery {
	if pageSize != 0 {
		l.PageSize = pageSize
	}

	return l
}

func (l *LogsQuery) WithStartTimeFilter(start time.Time) *LogsQuery {
	if !start.IsZero() {
		l.Filters.StartTime = start
	}

	return l
}

func (l *LogsQuery) WithEndTimeFilter(end time.Time) *LogsQuery {
	if !end.IsZero() {
		l.Filters.EndTime = end
	}

	return l
}
