# Datadog

~~~~sql
SELECT
    thread_a.thread_id,
    thread_a.processlist_id,
    thread_a.processlist_user,
    thread_a.processlist_host,
    thread_a.processlist_db,
    thread_a.processlist_command,
    thread_a.processlist_state,
    COALESCE(statement.sql_text, thread_a.PROCESSLIST_info) AS sql_text,
    statement.timer_start AS event_timer_start,
    statement.timer_end AS event_timer_end,
    statement.lock_time,
    statement.current_schema,
    waits_a.event_id AS event_id,
    waits_a.end_event_id,
    IF(waits_a.thread_id IS NULL,
        'other',
        COALESCE(
            IF(thread_a.processlist_state = 'User sleep', 'User sleep',
            IF(waits_a.event_id = waits_a.end_event_id, 'CPU', waits_a.event_name)), 'CPU'
        )
    ) AS wait_event,
    waits_a.operation,
    waits_a.timer_start AS wait_timer_start,
    waits_a.timer_end AS wait_timer_end,
    waits_a.object_schema,
    waits_a.object_name,
    waits_a.index_name,
    waits_a.object_type,
    waits_a.source
FROM
    performance_schema.threads AS thread_a
    LEFT JOIN performance_schema.events_waits_current AS waits_a ON waits_a.thread_id = thread_a.thread_id
    LEFT JOIN performance_schema.events_statements_current AS statement ON statement.thread_id = thread_a.thread_id
WHERE
    thread_a.processlist_state IS NOT NULL
    AND thread_a.processlist_command != 'Sleep'
    AND thread_a.processlist_id != CONNECTION_ID()
    AND thread_a.PROCESSLIST_COMMAND != 'Daemon'
    AND (waits_a.EVENT_NAME != 'idle' OR waits_a.EVENT_NAME IS NULL)
    AND (waits_a.operation != 'idle' OR waits_a.operation IS NULL)
    -- events_waits_current can have multiple rows per thread, thus we use EVENT_ID to identify the row we want to use.
    -- Additionally, we want the row with the highest EVENT_ID which reflects the most recent and current wait.
    AND (
        waits_a.event_id = (
           SELECT
              MAX(waits_b.EVENT_ID)
          FROM  performance_schema.events_waits_current AS waits_b
          Where waits_b.thread_id = thread_a.thread_id
    ) OR waits_a.event_id is NULL)
    -- We ignore rows without SQL text because there will be rows for background operations that do not have
    -- SQL text associated with it.
    AND COALESCE(statement.sql_text, thread_a.PROCESSLIST_info) != "";
~~~~

~~~~sql
SELECT `schema_name`,
                   `digest`,
                   `digest_text`,
                   `count_star`,
                   `sum_timer_wait`,
                   `sum_lock_time`,
                   `sum_errors`,
                   `sum_rows_affected`,
                   `sum_rows_sent`,
                   `sum_rows_examined`,
                   `sum_select_scan`,
                   `sum_select_full_join`,
                   `sum_no_index_used`,
                   `sum_no_good_index_used`
            FROM performance_schema.events_statements_summary_by_digest
            WHERE `digest_text` NOT LIKE 'EXPLAIN %' OR `digest_text` IS NULL
            ORDER BY `count_star` DESC
            LIMIT 10000
~~~~

~~~~sql
SELECT current_schema, sql_text, digest, digest_text, timer_start, @startup_time_s+timer_end*1e-12 as timer_end_time_s, timer_wait / 1000 AS timer_wait_ns, lock_time / 1000 AS lock_time_ns, rows_affected, rows_sent, rows_examined, select_full_join, select_full_range_join, select_range, select_range_check, select_scan, sort_merge_passes, sort_range, sort_rows, sort_scan, no_index_used, no_good_index_used, processlist_user, processlist_host, processlist_db FROM performance_schema.events_statements_current E LEFT JOIN performance_schema.threads as T ON E.thread_id = T.thread_id WHERE sql_text IS NOT NULL AND event_name like 'statement/%%' AND digest_text is NOT NULL AND digest_text NOT LIKE 'EXPLAIN %%' ORDER BY timer_wait DESC
~~~~

~~~~sql
set @startup_time_s = (SELECT UNIX_TIMESTAMP()-VARIABLE_VALUE FROM performance_schema.global_status WHERE VARIABLE_NAME='UPTIME')
~~~~

~~~~sql
CALL datadog.explain_statement('SELECT current_schema, sql_text, digest, digest_text, timer_start, @startup_time_s+timer_end*1e-12 as timer_end_time_s, timer_wait / 1000 AS timer_wait_ns, lock_time / 1000 AS lock_time_ns, rows_affected, rows_sent, rows_examined, select_full_join, select_full_range_join, select_range, select_range_check, select_scan, sort_merge_passes, sort_range, sort_rows, sort_scan, no_index_used, no_good_index_used, processlist_user, processlist_host, processlist_db FROM performance_schema.events_statements_current E LEFT JOIN performance_schema.threads as T ON E.thread_id = T.thread_id WHERE sql_text IS NOT NULL AND event_name like \'statement/%%\' AND digest_text is NOT NULL AND digest_text NOT LIKE \'EXPLAIN %%\' ORDER BY timer_wait DESC')
~~~~

~~~~sql
SELECT
            COUNT(processlist_user) AS connections,
            processlist_user,
            processlist_host,
            processlist_db,
            processlist_state
        FROM
            performance_schema.threads
        WHERE
            processlist_user IS NOT NULL AND
            processlist_state IS NOT NULL
        GROUP BY processlist_user, processlist_host, processlist_db, processlist_state
~~~~

~~~~sql
SHOW VARIABLES LIKE 'pid_file'
~~~~

~~~~sql
SELECT * FROM INFORMATION_SCHEMA.PROCESSLIST WHERE COMMAND LIKE '%Binlog dump%'
~~~~

~~~~sql
SHOW MASTER STATUS;
~~~~

~~~~sql
SHOW REPLICA STATUS;
~~~~

~~~~sql
SELECT plugin_status
FROM information_schema.plugins WHERE plugin_name='group_replication'
~~~~

~~~~sql
SHOW BINARY LOGS;
~~~~

~~~~sql
SHOW /*!50000 ENGINE*/ INNODB STATUS
~~~~

~~~~sql
SELECT engine
FROM information_schema.ENGINES
WHERE engine='InnoDB' and support != 'no' and support != 'disabled'
~~~~

~~~~sql
SHOW GLOBAL VARIABLES;
~~~~

~~~~sql
SHOW /*!50002 GLOBAL */ STATUS;
~~~~

~~~~sql
SELECT VERSION()
~~~~

~~~~sql

~~~~

~~~~sql

~~~~

~~~~sql

~~~~

~~~~sql

~~~~

~~~~sql

~~~~
