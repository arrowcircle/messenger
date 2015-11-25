CREATE TABLE dialogs(
  id int NOT NULL PRIMARY KEY,
  name char(255),
  created_at timestamp(6),
  updated_at timestamp(6),
  last_message_id int
);

CREATE INDEX index_dialogs_on_last_message_id ON dialogs USING btree(last_message_id DESC);
