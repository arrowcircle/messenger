CREATE TABLE dialogs(
  id int8 NOT NULL PRIMARY KEY,
  name char(255),
  created_at timestamp(6),
  updated_at timestamp(6)
);

CREATE INDEX index_dialogs_on_updated_at ON dialogs USING btree(updated_at ASC NULLS LAST);
