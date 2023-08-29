-- DROP TABLE stat_count;

CREATE TABLE stat_count (
                            id bigserial NOT NULL,
                            chat_id int8 NULL,
                            stat_type int8 NULL,
                            user_id int8 NULL,
                            ts int8 NULL,
                            count int8 NULL,
                            CONSTRAINT stat_count_pkey PRIMARY KEY (id)
);
CREATE INDEX stat_count_chat_id_idx ON stat_count USING btree (chat_id, ts, stat_type);
CREATE UNIQUE INDEX stat_count_chat_id_ts_idx ON stat_count USING btree (chat_id, stat_type, user_id, ts);
