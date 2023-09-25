package parser

const (
	KeywordAdd          = "ADD"
	KeywordAfter        = "AFTER"
	KeywordAlias        = "ALIAS"
	KeywordAll          = "ALL"
	KeywordAlter        = "ALTER"
	KeywordAnd          = "AND"
	KeywordAnti         = "ANTI"
	KeywordAny          = "ANY"
	KeywordArray        = "ARRAY"
	KeywordAs           = "AS"
	KeywordAsc          = "ASC"
	KeywordAscending    = "ASCENDING"
	KeywordAsof         = "ASOF"
	KeywordAst          = "AST"
	KeywordAsync        = "ASYNC"
	KeywordAttach       = "ATTACH"
	KeywordBetween      = "BETWEEN"
	KeywordBoth         = "BOTH"
	KeywordBy           = "BY"
	KeywordCache        = "CACHE"
	KeywordCase         = "CASE"
	KeywordCast         = "CAST"
	KeywordCheck        = "CHECK"
	KeywordClear        = "CLEAR"
	KeywordCluster      = "CLUSTER"
	KeywordCodec        = "CODEC"
	KeywordCollate      = "COLLATE"
	KeywordColumn       = "COLUMN"
	KeywordComment      = "COMMENT"
	KeywordCompiled     = "COMPILED"
	KeywordConstraint   = "CONSTRAINT"
	KeywordCreate       = "CREATE"
	KeywordCross        = "CROSS"
	KeywordCube         = "CUBE"
	KeywordCurrent      = "CURRENT"
	KeywordDatabase     = "DATABASE"
	KeywordDatabases    = "DATABASES"
	KeywordDate         = "DATE"
	KeywordDay          = "DAY"
	KeywordDeduplicate  = "DEDUPLICATE"
	KeywordDefault      = "DEFAULT"
	KeywordDelay        = "DELAY"
	KeywordDelete       = "DELETE"
	KeywordDesc         = "DESC"
	KeywordDescending   = "DESCENDING"
	KeywordDescribe     = "DESCRIBE"
	KeywordDetach       = "DETACH"
	KeywordDetached     = "DETACHED"
	KeywordDictionaries = "DICTIONARIES"
	KeywordDictionary   = "DICTIONARY"
	KeywordDisk         = "DISK"
	KeywordDistinct     = "DISTINCT"
	KeywordDistributed  = "DISTRIBUTED"
	KeywordDrop         = "DROP"
	KeywordDNS          = "DNS"
	KeywordElse         = "ELSE"
	KeywordEnd          = "END"
	KeywordEngine       = "ENGINE"
	KeywordEvents       = "EVENTS"
	KeywordExcept       = "EXCEPT"
	KeywordExists       = "EXISTS"
	KeywordExplain      = "EXPLAIN"
	KeywordExpression   = "EXPRESSION"
	KeywordExtract      = "EXTRACT"
	KeywordFetches      = "FETCHES"
	KeywordFileSystem   = "FILESYSTEM"
	KeywordFinal        = "FINAL"
	KeywordFirst        = "FIRST"
	KeywordFlush        = "FLUSH"
	KeywordFollowing    = "FOLLOWING"
	KeywordFor          = "FOR"
	KeywordFormat       = "FORMAT"
	KeywordFreeze       = "FREEZE"
	KeywordFrom         = "FROM"
	KeywordFull         = "FULL"
	KeywordFunction     = "FUNCTION"
	KeywordGlobal       = "GLOBAL"
	KeywordGranularity  = "GRANULARITY"
	KeywordGroup        = "GROUP"
	KeywordHaving       = "HAVING"
	KeywordHierarchical = "HIERARCHICAL"
	KeywordHour         = "HOUR"
	KeywordId           = "ID"
	KeywordIf           = "IF"
	KeywordIlike        = "ILIKE"
	KeywordIn           = "IN"
	KeywordIndex        = "INDEX"
	KeywordInf          = "INF"
	KeywordInjective    = "INJECTIVE"
	KeywordInner        = "INNER"
	KeywordInsert       = "INSERT"
	KeywordInterval     = "INTERVAL"
	KeywordInto         = "INTO"
	KeywordIs           = "IS"
	KeywordIs_object_id = "IS_OBJECT_ID"
	KeywordJoin         = "JOIN"
	KeywordKey          = "KEY"
	KeywordKill         = "KILL"
	KeywordLast         = "LAST"
	KeywordLayout       = "LAYOUT"
	KeywordLeading      = "LEADING"
	KeywordLeft         = "LEFT"
	KeywordLifetime     = "LIFETIME"
	KeywordLike         = "LIKE"
	KeywordLimit        = "LIMIT"
	KeywordLive         = "LIVE"
	KeywordLocal        = "LOCAL"
	KeywordLogs         = "LOGS"
	KeywordMark         = "MARK"
	KeywordMaterialize  = "MATERIALIZE"
	KeywordMaterialized = "MATERIALIZED"
	KeywordMax          = "MAX"
	KeywordMerges       = "MERGES"
	KeywordMin          = "MIN"
	KeywordMinute       = "MINUTE"
	KeywordModify       = "MODIFY"
	KeywordMonth        = "MONTH"
	KeywordMove         = "MOVE"
	KeywordMutation     = "MUTATION"
	KeywordNan_sql      = "NAN_SQL"
	KeywordNo           = "NO"
	KeywordNot          = "NOT"
	KeywordNull         = "NULL"
	KeywordNulls        = "NULLS"
	KeywordOffset       = "OFFSET"
	KeywordOn           = "ON"
	KeywordOptimize     = "OPTIMIZE"
	KeywordOr           = "OR"
	KeywordOrder        = "ORDER"
	KeywordOuter        = "OUTER"
	KeywordOutfile      = "OUTFILE"
	KeywordOver         = "OVER"
	KeywordPartition    = "PARTITION"
	KeywordPopulate     = "POPULATE"
	KeywordPreceding    = "PRECEDING"
	KeywordPrewhere     = "PREWHERE"
	KeywordPrimary      = "PRIMARY"
	KeywordProjection   = "PROJECTION"
	KeywordQuarter      = "QUARTER"
	KeywordQuery        = "QUERY"
	KeywordRange        = "RANGE"
	KeywordReload       = "RELOAD"
	KeywordRemove       = "REMOVE"
	KeywordRename       = "RENAME"
	KeywordReplace      = "REPLACE"
	KeywordReplica      = "REPLICA"
	KeywordReplicated   = "REPLICATED"
	KeywordRight        = "RIGHT"
	KeywordRollup       = "ROLLUP"
	KeywordRow          = "ROW"
	KeywordRows         = "ROWS"
	KeywordSample       = "SAMPLE"
	KeywordSecond       = "SECOND"
	KeywordSelect       = "SELECT"
	KeywordSemi         = "SEMI"
	KeywordSends        = "SENDS"
	KeywordSet          = "SET"
	KeywordSettings     = "SETTINGS"
	KeywordShow         = "SHOW"
	KeywordSource       = "SOURCE"
	KeywordStart        = "START"
	KeywordStop         = "STOP"
	KeywordSubstring    = "SUBSTRING"
	KeywordSync         = "SYNC"
	KeywordSyntax       = "SYNTAX"
	KeywordSystem       = "SYSTEM"
	KeywordTable        = "TABLE"
	KeywordTables       = "TABLES"
	KeywordTemporary    = "TEMPORARY"
	KeywordTest         = "TEST"
	KeywordThen         = "THEN"
	KeywordTies         = "TIES"
	KeywordTimeout      = "TIMEOUT"
	KeywordTimestamp    = "TIMESTAMP"
	KeywordTo           = "TO"
	KeywordTop          = "TOP"
	KeywordTotals       = "TOTALS"
	KeywordTrailing     = "TRAILING"
	KeywordTrim         = "TRIM"
	KeywordTruncate     = "TRUNCATE"
	KeywordTtl          = "TTL"
	KeywordType         = "TYPE"
	KeywordUnbounded    = "UNBOUNDED"
	KeywordUncompressed = "UNCOMPRESSED"
	KeywordUnion        = "UNION"
	KeywordUpdate       = "UPDATE"
	KeywordUse          = "USE"
	KeywordUsing        = "USING"
	KeywordUuid         = "UUID"
	KeywordValues       = "VALUES"
	KeywordView         = "VIEW"
	KeywordVolume       = "VOLUME"
	KeywordWatch        = "WATCH"
	KeywordWeek         = "WEEK"
	KeywordWhen         = "WHEN"
	KeywordWhere        = "WHERE"
	KeywordWindow       = "WINDOW"
	KeywordWith         = "WITH"
	KeywordYear         = "YEAR"
)

var keywords = NewSet(
	KeywordAdd,
	KeywordAfter,
	KeywordAlias,
	KeywordAll,
	KeywordAlter,
	KeywordAnd,
	KeywordAnti,
	KeywordAny,
	KeywordArray,
	KeywordAs,
	KeywordAsc,
	KeywordAscending,
	KeywordAsof,
	KeywordAst,
	KeywordAsync,
	KeywordAttach,
	KeywordBetween,
	KeywordBoth,
	KeywordBy,
	KeywordCache,
	KeywordCase,
	KeywordCast,
	KeywordCheck,
	KeywordClear,
	KeywordCluster,
	KeywordCodec,
	KeywordCollate,
	KeywordColumn,
	KeywordComment,
	KeywordCompiled,
	KeywordConstraint,
	KeywordCreate,
	KeywordCross,
	KeywordCube,
	KeywordCurrent,
	KeywordDatabase,
	KeywordDatabases,
	KeywordDate,
	KeywordDay,
	KeywordDeduplicate,
	KeywordDefault,
	KeywordDelay,
	KeywordDelete,
	KeywordDesc,
	KeywordDescending,
	KeywordDescribe,
	KeywordDetach,
	KeywordDetached,
	KeywordDictionaries,
	KeywordDictionary,
	KeywordDisk,
	KeywordDistinct,
	KeywordDistributed,
	KeywordDrop,
	KeywordDNS,
	KeywordElse,
	KeywordEnd,
	KeywordEngine,
	KeywordEvents,
	KeywordExcept,
	KeywordExists,
	KeywordExplain,
	KeywordExpression,
	KeywordExtract,
	KeywordFetches,
	KeywordFileSystem,
	KeywordFinal,
	KeywordFirst,
	KeywordFlush,
	KeywordFollowing,
	KeywordFor,
	KeywordFormat,
	KeywordFreeze,
	KeywordFrom,
	KeywordFull,
	KeywordFunction,
	KeywordGlobal,
	KeywordGranularity,
	KeywordGroup,
	KeywordHaving,
	KeywordHierarchical,
	KeywordHour,
	KeywordId,
	KeywordIf,
	KeywordIlike,
	KeywordIn,
	KeywordIndex,
	KeywordInf,
	KeywordInjective,
	KeywordInner,
	KeywordInsert,
	KeywordInterval,
	KeywordInto,
	KeywordIs,
	KeywordIs_object_id,
	KeywordJoin,
	KeywordKey,
	KeywordKill,
	KeywordLast,
	KeywordLayout,
	KeywordLeading,
	KeywordLeft,
	KeywordLifetime,
	KeywordLike,
	KeywordLimit,
	KeywordLive,
	KeywordLocal,
	KeywordLogs,
	KeywordMark,
	KeywordMaterialize,
	KeywordMaterialized,
	KeywordMax,
	KeywordMerges,
	KeywordMin,
	KeywordMinute,
	KeywordModify,
	KeywordMonth,
	KeywordMove,
	KeywordMutation,
	KeywordNan_sql,
	KeywordNo,
	KeywordNot,
	KeywordNull,
	KeywordNulls,
	KeywordOffset,
	KeywordOn,
	KeywordOptimize,
	KeywordOr,
	KeywordOrder,
	KeywordOuter,
	KeywordOutfile,
	KeywordOver,
	KeywordPartition,
	KeywordPopulate,
	KeywordPreceding,
	KeywordPrewhere,
	KeywordPrimary,
	KeywordProjection,
	KeywordQuarter,
	KeywordQuery,
	KeywordRange,
	KeywordReload,
	KeywordRemove,
	KeywordRename,
	KeywordReplace,
	KeywordReplica,
	KeywordReplicated,
	KeywordRight,
	KeywordRollup,
	KeywordRow,
	KeywordRows,
	KeywordSample,
	KeywordSecond,
	KeywordSelect,
	KeywordSemi,
	KeywordSends,
	KeywordSet,
	KeywordSettings,
	KeywordShow,
	KeywordSource,
	KeywordStart,
	KeywordStop,
	KeywordSubstring,
	KeywordSync,
	KeywordSyntax,
	KeywordSystem,
	KeywordTable,
	KeywordTables,
	KeywordTemporary,
	KeywordTest,
	KeywordThen,
	KeywordTies,
	KeywordTimeout,
	KeywordTimestamp,
	KeywordTo,
	KeywordTop,
	KeywordTotals,
	KeywordTrailing,
	KeywordTrim,
	KeywordTruncate,
	KeywordTtl,
	KeywordType,
	KeywordUnbounded,
	KeywordUncompressed,
	KeywordUnion,
	KeywordUpdate,
	KeywordUse,
	KeywordUsing,
	KeywordUuid,
	KeywordValues,
	KeywordView,
	KeywordVolume,
	KeywordWatch,
	KeywordWeek,
	KeywordWhen,
	KeywordWhere,
	KeywordWindow,
	KeywordWith,
	KeywordYear,
)
