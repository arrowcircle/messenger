package main

import "strconv"

// GetDialogs fetches dialogs from db for user with offset
func (i *Impl) GetDialogs(userID string, offset int) []DialogJSON {
	dialogs := []DialogJSON{}
	i.DB.Raw(`
    SELECT c.*, array_agg(du.user_id) AS user_ids
    FROM
      (SELECT
        dialogs.id AS id,
        dialogs.name AS name,
        dialogs.created_at AS created_at,
        dialogs.updated_at AS updated_at,
        dialogs.last_message_id AS last_message_id,
        messages.text AS last_message,
        messages.user_id AS last_message_user_id,
    	  dialog_users.last_seen_message_id AS last_seen_message_id
      FROM dialogs
      JOIN messages ON messages.id = dialogs.last_message_id
      JOIN dialog_users ON dialog_users.dialog_id = dialogs.id
      WHERE dialog_users.user_id = ?
      ORDER BY dialogs.last_message_id DESC
      ) c
    JOIN dialog_users du ON c.id = du.dialog_id
    GROUP BY
      c.id,
      c.name,
      c.created_at,
      c.updated_at,
      c.last_message_id,
      c.last_message,
      c.last_message_user_id,
      c.last_seen_message_id
		ORDER BY c.last_message_id DESC
    LIMIT 10
    OFFSET ?
  `, userID, offset).Find(&dialogs)
	return dialogs
}

// ShowDialog gets one dialog for user
func (i *Impl) ShowDialog(userID string, dialogID int) DialogJSON {
	dialog := DialogJSON{}
	i.DB.Raw(`
    SELECT c.*, array_agg(du.user_id) AS user_ids
    FROM
      (SELECT
        dialogs.id AS id,
        dialogs.name AS name,
        dialogs.created_at AS created_at,
        dialogs.updated_at AS updated_at,
        dialogs.last_message_id AS last_message_id,
        messages.text AS last_message,
        messages.user_id AS last_message_user_id,
    	  dialog_users.last_seen_message_id AS last_seen_message_id
      FROM dialogs
      JOIN messages ON messages.id = dialogs.last_message_id
      JOIN dialog_users ON dialog_users.dialog_id = dialogs.id
      WHERE dialog_users.user_id = ?
      ORDER BY dialogs.last_message_id DESC
      ) c
    JOIN dialog_users du ON c.id = du.dialog_id
    WHERE c.id = ?
    GROUP BY
      c.id,
      c.name,
      c.created_at,
      c.updated_at,
      c.last_message_id,
      c.last_message,
      c.last_message_user_id,
      c.last_seen_message_id
  `, userID, dialogID).Find(&dialog)

	i.UpdateLastMessage(userID, dialogID)

	return dialog
}

// IndexMessages fetches messages for the dialog
func (i *Impl) IndexMessages(userID string, dialogID int, offset int) []MessageJSON {
	messages := []MessageJSON{}
	i.DB.Raw(`
    SELECT * FROM messages
    WHERE messages.dialog_id = ?
    ORDER BY messages.id DESC
    LIMIT 10
    OFFSET ?
  `, dialogID, offset).Find(&messages)

	i.UpdateLastMessage(userID, dialogID)
	return messages
}

// ShowUser fetches user and number of unread dialogs
func (i *Impl) ShowUser(userID string) UserJSON {
	user := UserJSON{}
	user.ID, _ = strconv.Atoi(userID)
	dialogsCount := 0
	i.DB.Raw(`
		SELECT COUNT(dialogs.id)
		FROM dialogs, dialog_users
		WHERE
	  dialogs.last_message_id > dialog_users.last_seen_message_id AND
	  dialog_users.user_id = ? AND
	  dialog_users.dialog_id = dialogs.id
	`, user.ID).Row().Scan(&dialogsCount)
	user.DialogsCount = dialogsCount
	return user
}

// CreateMessage inserts message in db
func (i *Impl) CreateMessage(userID string, dialogID int, message Message) (Message, error) {
	message.DialogID = dialogID
	message.UserID, _ = strconv.Atoi(userID)

	if err := i.DB.Save(&message).Error; err != nil {
		return message, err
	}

	i.DB.Exec("UPDATE dialogs SET last_message_id = ?, updated_at = now() WHERE dialogs.id = ?", message.ID, message.DialogID)
	i.DB.Exec("UPDATE dialog_users SET last_seen_message_id = ? WHERE dialog_id = ? AND user_id = ?", message.ID, message.DialogID, message.UserID)

	return message, nil
}

// CreateDialog creates dialog. If ony two people are present, it uses existing dialog
func (i *Impl) CreateDialog(userID string, params DialogCreateJSON) (Dialog, error) {
	dialog, err := i.FindDialogByUserIds(params)
	if err != nil {
		dialog := Dialog{}
		dialog.Name = params.Name
		if err := i.DB.Save(&dialog).Error; err != nil {
			return dialog, err
		}
	}

	message := Message{}
	message.DialogID = dialog.ID
	message.Text = params.Message
	message.UserID, _ = strconv.Atoi(userID)

	if err := i.DB.Save(&message).Error; err != nil {
		return dialog, err
	}

	for _, element := range params.UserIds {
		i.DB.Exec("INSERT INTO dialog_users (dialog_id, user_id, last_seen_message_id) VALUES (?, ?, 0)", dialog.ID, element)
	}
	i.DB.Exec("UPDATE dialogs SET last_message_id = ? WHERE id = ?", message.ID, dialog.ID)
	dialog.LastMessageID = message.ID

	return dialog, nil
}

// FindDialogByUserIds return dialog for two users if it exists
func (i *Impl) FindDialogByUserIds(params DialogCreateJSON) (Dialog, error) {
	dialogID := 0
	if len(params.UserIds) == 2 {
		i.DB.Raw(`SELECT dialog_users.dialog_id
    FROM dialog_users
    WHERE dialog_users.user_id = ?
    UNION
    SELECT dialog_users.dialog_id
    FROM dialog_users
    WHERE dialog_users.user_id = ?
    ORDER BY dialog_id DESC
    LIMIT 1`, params.UserIds[0], params.UserIds[1]).Row().Scan(&dialogID)
	}

	dialog := Dialog{}
	if dialogID == 0 {
		dialog.Name = params.Name
		if err := i.DB.Save(&dialog).Error; err != nil {
			return dialog, err
		}
	} else {
		i.DB.Find(&dialog, dialogID)
	}

	return dialog, nil
}

// UpdateLastMessage sets last_message_id for dialog
func (i *Impl) UpdateLastMessage(userID string, dialogID int) {
	lastMessageID := 0
	i.DB.Raw("SELECT last_message_id FROM dialogs WHERE id = ?", dialogID).Row().Scan(&lastMessageID)
	i.DB.Exec("UPDATE dialog_users SET last_seen_message_id = ? WHERE dialog_id = ? AND user_id = ?", lastMessageID, dialogID, userID)
}
