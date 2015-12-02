CREATE TABLE dialogs(
  id SERIAL PRIMARY KEY,
  name varchar(255),
  created_at timestamp(6) DEFAULT now(),
  updated_at timestamp(6) DEFAULT now(),
  last_message_id SERIAL
);

CREATE INDEX index_dialogs_on_last_message_id ON dialogs USING btree(last_message_id DESC);
