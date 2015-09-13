CREATE TABLE messages(
  id int8 NOT NULL PRIMARY KEY,
  dialog_id int8 NOT NULL,
  text TEXT NOT NULL,
  user_id int4 NOT NULL,
  created_at timestamp(6)
);

CREATE INDEX index_messages_on_dialog_id ON messages USING btree(dialog_id ASC NULLS LAST);
CREATE INDEX index_messages_on_user_id ON messages USING btree(user_id ASC NULLS LAST);
