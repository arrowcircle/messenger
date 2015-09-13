CREATE TABLE dialog_users(
  dialog_id int8 not null,
  user_id int8 not null,
  created_at timestamp(6),
  updated_at timestamp(6)
);

CREATE INDEX  index_dialog_users_on_dialog_id ON dialog_users USING btree(dialog_id ASC NULLS LAST);
CREATE INDEX  index_dialog_users_on_user_id ON dialog_users USING btree(user_id ASC NULLS LAST);
CREATE INDEX  index_dialog_users_on_updated_at ON dialog_users USING btree(updated_at ASC NULLS LAST);
