package parser

const (
	KeywordAdd          = "ADD"
	KeywordAdmin        = "ADMIN"
	KeywordAfter        = "AFTER"
	KeywordAlias        = "ALIAS"
	KeywordAll          = "ALL"
	KeywordAlter        = "ALTER"
	KeywordAnd          = "AND"
	KeywordAnti         = "ANTI"
	KeywordAny          = "ANY"
	KeywordAppend       = "APPEND"
	KeywordApply        = "APPLY"
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
	KeywordColumns      = "COLUMNS"
	KeywordComment      = "COMMENT"
	KeywordCompiled     = "COMPILED"
	KeywordConfig       = "CONFIG"
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
	KeywordDepends      = "DEPENDS"
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
	KeywordEmbedded     = "EMBEDDED"
	KeywordEmpty        = "EMPTY"
	KeywordEnd          = "END"
	KeywordEngine       = "ENGINE"
	KeywordEstimate     = "ESTIMATE"
	KeywordEvents       = "EVENTS"
	KeywordEvery        = "EVERY"
	KeywordExcept       = "EXCEPT"
	KeywordExists       = "EXISTS"
	KeywordExplain      = "EXPLAIN"
	KeywordExpression   = "EXPRESSION"
	KeywordExtract      = "EXTRACT"
	KeywordFalse        = "FALSE"
	KeywordFetches      = "FETCHES"
	KeywordFileSystem   = "FILESYSTEM"
	KeywordFill         = "FILL"
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
	KeywordFunctions    = "FUNCTIONS"
	KeywordGlobal       = "GLOBAL"
	KeywordGrant        = "GRANT"
	KeywordGrantees     = "GRANTEES"
	KeywordGranularity  = "GRANULARITY"
	KeywordGroup        = "GROUP"
	KeywordGrouping     = "GROUPING"
	KeywordHaving       = "HAVING"
	KeywordHierarchical = "HIERARCHICAL"
	KeywordHost         = "HOST"
	KeywordHour         = "HOUR"
	KeywordId           = "ID"
	KeywordIdentified   = "IDENTIFIED"
	KeywordIf           = "IF"
	KeywordIlike        = "ILIKE"
	KeywordIn           = "IN"
	KeywordIndex        = "INDEX"
	KeywordInf          = "INF"
	KeywordInjective    = "INJECTIVE"
	KeywordInner        = "INNER"
	KeywordInsert       = "INSERT"
	KeywordInterval     = "INTERVAL"
	KeywordInterpolate  = "INTERPOLATE"
	KeywordInto         = "INTO"
	KeywordIp           = "IP"
	KeywordIs           = "IS"
	KeywordIs_object_id = "IS_OBJECT_ID"
	KeywordJoin         = "JOIN"
	KeywordJSON         = "JSON"
	KeywordKey          = "KEY"
	KeywordKill         = "KILL"
	KeywordKerberos     = "KERBEROS"
	KeywordLast         = "LAST"
	KeywordLayout       = "LAYOUT"
	KeywordLdap         = "LDAP"
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
	KeywordMoves        = "MOVES"
	KeywordMutation     = "MUTATION"
	KeywordName         = "NAME"
	KeywordNan_sql      = "NAN_SQL"
	KeywordNo           = "NO"
	KeywordNone         = "NONE"
	KeywordNot          = "NOT"
	KeywordNull         = "NULL"
	KeywordNulls        = "NULLS"
	KeywordOffset       = "OFFSET"
	KeywordOn           = "ON"
	KeywordOptimize     = "OPTIMIZE"
	KeywordOption       = "OPTION"
	KeywordOr           = "OR"
	KeywordOrder        = "ORDER"
	KeywordOuter        = "OUTER"
	KeywordOutfile      = "OUTFILE"
	KeywordOver         = "OVER"
	KeywordPartition    = "PARTITION"
	KeywordPipeline     = "PIPELINE"
	KeywordPolicy       = "POLICY"
	KeywordPopulate     = "POPULATE"
	KeywordPreceding    = "PRECEDING"
	KeywordPrewhere     = "PREWHERE"
	KeywordPrimary      = "PRIMARY"
	KeywordProjection   = "PROJECTION"
	KeywordQuarter      = "QUARTER"
	KeywordQuery        = "QUERY"
	KeywordQueues       = "QUEUES"
	KeywordQuota        = "QUOTA"
	KeywordRandomize    = "RANDOMIZE"
	KeywordRange        = "RANGE"
	KeywordRealm        = "REALM"
	KeywordRecompress   = "RECOMPRESS"
	KeywordRefresh      = "REFRESH"
	KeywordRegexp       = "REGEXP"
	KeywordReload       = "RELOAD"
	KeywordRemove       = "REMOVE"
	KeywordRename       = "RENAME"
	KeywordReplace      = "REPLACE"
	KeywordReset        = "RESET"
	KeywordReplica      = "REPLICA"
	KeywordReplicated   = "REPLICATED"
	KeywordReplication  = "REPLICATION"
	KeywordRestart      = "RESTART"
	KeywordRight        = "RIGHT"
	KeywordRole         = "ROLE"
	KeywordRollup       = "ROLLUP"
	KeywordRow          = "ROW"
	KeywordRows         = "ROWS"
	KeywordSample       = "SAMPLE"
	KeywordSecond       = "SECOND"
	KeywordSelect       = "SELECT"
	KeywordSemi         = "SEMI"
	KeywordSends        = "SENDS"
	KeywordServer       = "SERVER"
	KeywordSet          = "SET"
	KeywordSets         = "SETS"
	KeywordSetting      = "SETTING"
	KeywordSettings     = "SETTINGS"
	KeywordShow         = "SHOW"
	KeywordShutdown     = "SHUTDOWN"
	KeywordSkip         = "SKIP"
	KeywordSource       = "SOURCE"
	KeywordStart        = "START"
	KeywordStaleness    = "STALENESS"
	KeywordStep         = "STEP"
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
	KeywordTrue         = "TRUE"
	KeywordTruncate     = "TRUNCATE"
	KeywordTtl          = "TTL"
	KeywordType         = "TYPE"
	KeywordUnbounded    = "UNBOUNDED"
	KeywordUncompressed = "UNCOMPRESSED"
	KeywordUnion        = "UNION"
	KeywordUpdate       = "UPDATE"
	KeywordUse          = "USE"
	KeywordUser         = "USER"
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
	KeywordDefiner      = "DEFINER"
	KeywordSQL          = "SQL"
	KeywordSecurity     = "SECURITY"
)

var keywords = NewSet(
	KeywordAdd,
	KeywordAdmin,
	KeywordAfter,
	KeywordAlias,
	KeywordAll,
	KeywordAlter,
	KeywordAnd,
	KeywordAnti,
	KeywordAny,
	KeywordAppend,
	KeywordApply,
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
	KeywordColumns,
	KeywordComment,
	KeywordCompiled,
	KeywordConfig,
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
	KeywordDepends,
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
	KeywordEstimate,
	KeywordEmbedded,
	KeywordEmpty,
	KeywordEvents,
	KeywordEvery,
	KeywordExcept,
	KeywordExists,
	KeywordExplain,
	KeywordExpression,
	KeywordExtract,
	KeywordFalse,
	KeywordFetches,
	KeywordFileSystem,
	KeywordFill,
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
	KeywordFunctions,
	KeywordGlobal,
	KeywordGrant,
	KeywordGrantees,
	KeywordGranularity,
	KeywordGroup,
	KeywordGrouping,
	KeywordHaving,
	KeywordHierarchical,
	KeywordHost,
	KeywordHour,
	KeywordId,
	KeywordIdentified,
	KeywordIf,
	KeywordIlike,
	KeywordIn,
	KeywordIndex,
	KeywordInf,
	KeywordInjective,
	KeywordInner,
	KeywordInsert,
	KeywordInterval,
	KeywordInterpolate,
	KeywordInto,
	KeywordIp,
	KeywordIs,
	KeywordIs_object_id,
	KeywordJoin,
	KeywordJSON,
	KeywordKey,
	KeywordKill,
	KeywordKerberos,
	KeywordLast,
	KeywordLayout,
	KeywordLdap,
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
	KeywordMoves,
	KeywordMutation,
	KeywordName,
	KeywordNan_sql,
	KeywordNo,
	KeywordNone,
	KeywordNot,
	KeywordNull,
	KeywordNulls,
	KeywordOffset,
	KeywordOn,
	KeywordOptimize,
	KeywordOption,
	KeywordOr,
	KeywordOrder,
	KeywordOuter,
	KeywordOutfile,
	KeywordOver,
	KeywordPartition,
	KeywordPipeline,
	KeywordPolicy,
	KeywordPopulate,
	KeywordPreceding,
	KeywordPrewhere,
	KeywordPrimary,
	KeywordProjection,
	KeywordQuarter,
	KeywordQuery,
	KeywordQueues,
	KeywordQuota,
	KeywordRandomize,
	KeywordRange,
	KeywordRealm,
	KeywordRecompress,
	KeywordRefresh,
	KeywordRegexp,
	KeywordReload,
	KeywordRemove,
	KeywordRename,
	KeywordReplace,
	KeywordReset,
	KeywordReplica,
	KeywordReplicated,
	KeywordReplication,
	KeywordRestart,
	KeywordRight,
	KeywordRole,
	KeywordRollup,
	KeywordRow,
	KeywordRows,
	KeywordSample,
	KeywordSecond,
	KeywordSelect,
	KeywordSemi,
	KeywordSends,
	KeywordServer,
	KeywordSet,
	KeywordSets,
	KeywordSetting,
	KeywordSettings,
	KeywordShow,
	KeywordShutdown,
	KeywordSkip,
	KeywordSource,
	KeywordStart,
	KeywordStaleness,
	KeywordStep,
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
	KeywordTrue,
	KeywordTruncate,
	KeywordTtl,
	KeywordType,
	KeywordUnbounded,
	KeywordUncompressed,
	KeywordUnion,
	KeywordUpdate,
	KeywordUse,
	KeywordUser,
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
	KeywordDefiner,
	KeywordSQL,
	KeywordSecurity,
)
