-- Origin SQL:
SET param_a = 13;
SET param_b = 'str';
SET param_c = '2022-08-04 18:30:53';
SET param_d = {'10': [11, 12], '13': [14, 15]};

SELECT
    {a: UInt32},
    {b: String},
    {c: DateTime},
    {d: Map(String, Array(UInt8))};

SELECT * FROM clickhouse WHERE tenant_id = {tenant_id: String};


-- Format SQL:
SET param_a=13;
SET param_b='str';
SET param_c='2022-08-04 18:30:53';
SET param_d={'10': [11, 12], '13': [14, 15]};
SELECT {a: UInt32}, {b: String}, {c: DateTime}, {d: Map(String, Array(UInt8))};
SELECT * FROM clickhouse WHERE tenant_id = {tenant_id: String};
