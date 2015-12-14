CREATE TABLE dialog_users(
  dialog_id SERIAL not null,
  user_id SERIAL not null,
  created_at timestamp(6) DEFAULT now(),
  updated_at timestamp(6) DEFAULT now(),
  last_seen_message_id SERIAL
);

CREATE INDEX  index_dialog_users_on_dialog_id ON dialog_users USING btree(dialog_id ASC NULLS LAST);
CREATE INDEX  index_dialog_users_on_user_id ON dialog_users USING btree(user_id ASC NULLS LAST);
CREATE INDEX  index_dialog_users_on_updated_at ON dialog_users USING btree(updated_at ASC NULLS LAST);
CREATE INDEX  index_dialog_users_on_last_seen_message_id ON dialog_users USING btree(last_seen_message_id);
