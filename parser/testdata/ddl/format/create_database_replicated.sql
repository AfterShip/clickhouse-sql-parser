-- Origin SQL:
CREATE DATABASE IF NOT EXISTS `test` ENGINE=Replicated('/root/test_local', 'shard', 'replica');


-- Format SQL:
CREATE DATABASE IF NOT EXISTS `test`  ENGINE = Replicated('/root/test_local', 'shard', 'replica');
