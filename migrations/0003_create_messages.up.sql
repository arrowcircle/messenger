CREATE TABLE messages(
  id SERIAL PRIMARY KEY,
  dialog_id SERIAL NOT NULL,
  text TEXT NOT NULL,
  user_id SERIAL NOT NULL,
  created_at timestamp(6) DEFAULT now()
);

CREATE INDEX index_messages_on_dialog_id ON messages USING btree(dialog_id);
CREATE INDEX index_messages_on_user_id ON messages USING btree(user_id);
