CREATE TABLE messages(
  id int NOT NULL PRIMARY KEY,
  dialog_id int NOT NULL,
  text TEXT NOT NULL,
  user_id int NOT NULL,
  created_at timestamp(6)
);

CREATE INDEX index_messages_on_dialog_id ON messages USING btree(dialog_id);
CREATE INDEX index_messages_on_user_id ON messages USING btree(user_id);
