INSERT INTO messages (dialog_id, text, user_id, created_at)
VALUES
  (1, 'id: 1, dialog_id: 1, user_id: 1', 1, now()),
  (1, 'id: 2, dialog_id: 1, user_id: 2', 2, CURRENT_TIMESTAMP + INTERVAL '1 minute'),
  (1, 'id: 3, dialog_id: 1, user_id: 1', 1, CURRENT_TIMESTAMP + INTERVAL '2 minutes'),
  (2, 'id: 4, dialog_id: 2, user_id: 2', 2, CURRENT_TIMESTAMP + INTERVAL '3 minutes'),
  (2, 'id: 5, dialog_id: 2, user_id: 1', 1, CURRENT_TIMESTAMP + INTERVAL '4 minutes'),
  (2, 'id: 6, dialog_id: 2, user_id: 2', 2, CURRENT_TIMESTAMP + INTERVAL '5 minutes'),
  (3, 'id: 7, dialog_id: 3, user_id: 2', 2, CURRENT_TIMESTAMP + INTERVAL '6 minutes');

INSERT INTO dialogs (name, created_at, updated_at, last_message_id)
VALUES
  ('test dialog 1', now(), now(), 3),
  ('test dialog 2', now(), now(), 6),
  ('test dialog 3', now(), now(), 7);

INSERT INTO dialog_users (dialog_id, user_id, created_at, updated_at, last_seen_message_id)
VALUES
  (1, 1, now(), now(), 3),
  (1, 2, now(), now(), 3),
  (2, 1, now(), now(), 5),
  (2, 2, now(), now(), 6),
  (3, 3, now(), now(), 0),
  (3, 2, now(), now(), 7),
  (1, 289622, now(), now(), 3);
