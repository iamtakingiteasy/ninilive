package postgres

import (
	"database/sql"
	"fmt"

	"github.com/iamtakingiteasy/ninilive/internal/config"
	"github.com/iamtakingiteasy/ninilive/internal/db/model"
	"github.com/jmoiron/sqlx"
)

type persister struct {
	values             *config.ValuesDB
	db                 *sqlx.DB
	deleteMessage      *sqlx.NamedStmt
	saveMessage        *sqlx.NamedStmt
	editMessage        *sqlx.NamedStmt
	loadMessagesLast   *sqlx.Stmt
	loadMessagesBefore *sqlx.Stmt
	loadMessagesPage   *sqlx.Stmt
	deleteUser         *sqlx.Stmt
	loadUser           *sqlx.Stmt
	loadUserByID       *sqlx.Stmt
	saveUser           *sqlx.NamedStmt
	ensureUser         *sqlx.NamedStmt
	loadUsersPage      *sqlx.Stmt
	deleteChannel      *sqlx.Stmt
	saveChannel        *sqlx.NamedStmt
	loadChannels       *sqlx.Stmt
}

const deleteMessageQuery = `
delete from messages where
  message_id = :message_id and 
  (
    message_origin = :message_origin or 
    (
      message_user = :user_id and
      exists(select * from users where user_login = :user_login and user_password = :user_password)
    ) or
    :user_mod
  )
`

func (persister *persister) DeleteMessage(message *model.Message) (err error) {
	var res sql.Result

	res, err = persister.deleteMessage.Exec(message)

	if err != nil {
		return
	}

	var rows int64

	rows, err = res.RowsAffected()
	if err != nil {
		return
	}

	if rows == 0 {
		return fmt.Errorf("not removed")
	}

	return
}

const saveMessageQuery = `
with u as (insert into users(
    user_id,
    user_name,
    user_login,
    user_password,
    user_mod
  ) values (
    coalesce(nullif(:user_id, 0), case when 
      :user_login = '' then 
        0
      else
        nextval('users_user_id_seq')
      end),
    :user_name,
    :user_login,
    :user_password,
    :user_mod
  ) on conflict (user_id) do update set user_id = EXCLUDED.user_id
    returning user_id, user_name, user_login, user_password, user_mod),
res as (insert into messages(
  message_id,
  message_channel_id,
  message_body,
  message_time,
  message_edit,
  message_trip,
  message_origin,
  message_remote,
  message_file_name,
  message_file_path,
  message_user
) select
  coalesce(nullif(:message_id, 0), nextval('messages_message_id_seq')),
  :message_channel_id,
  :message_body,
  now(),
  now(),
  :message_trip,
  :message_origin,
  :message_remote,
  :message_file_name,
  :message_file_path,
  user_id from u
on conflict (message_id) do update set
  message_body = EXCLUDED.message_body,
  message_edit = EXCLUDED.message_edit,
  message_file_name = EXCLUDED.message_file_name,
  message_file_path = EXCLUDED.message_file_path
returning *
)
select 
  coalesce(m.message_id, 0) message_id,
  coalesce(m.message_channel_id, 0) message_channel_id,
  coalesce(m.message_body, '') message_body,
  coalesce(m.message_time, 'epoch') message_time,
  coalesce(m.message_edit, 'epoch') message_edit,
  coalesce(m.message_trip, '') message_trip,
  coalesce(m.message_origin, '') message_origin,
  coalesce(m.message_remote, '') message_remote,
  coalesce(m.message_file_name, '') message_file_name,
  coalesce(m.message_file_path, '') message_file_path,
  coalesce(us.user_id, 0) user_id,
  coalesce(us.user_name, '') user_name,
  coalesce(us.user_login, '') user_login,
  coalesce(us.user_password, '') user_password,
  coalesce(us.user_mod, false) user_mod
from res m left join users us on
  m.message_user = us.user_id
`

func (persister *persister) SaveMessage(msg *model.Message) (err error) {
	row := persister.saveMessage.QueryRowx(msg)

	if row.Err() != nil {
		return row.Err()
	}

	return row.StructScan(msg)
}

const editMessageQuery = `
with res as (
update messages set
  message_body = :message_body,
  message_edit = now(),
  message_file_name = :message_file_name,
  message_file_path = :message_file_path
where
  message_id = :message_id and 
  (
    message_origin = :message_origin or 
    (
      message_user = :user_id and
      exists(select * from users where user_login = :user_login and user_password = :user_password)
    ) or
    :user_mod
  )
returning *)
select
  coalesce(m.message_id, 0) message_id,
  coalesce(m.message_channel_id, 0) message_channel_id,
  coalesce(m.message_body, '') message_body,
  coalesce(m.message_time, 'epoch') message_time,
  coalesce(m.message_edit, 'epoch') message_edit,
  coalesce(m.message_trip, '') message_trip,
  coalesce(m.message_origin, '') message_origin,
  coalesce(m.message_remote, '') message_remote,
  coalesce(m.message_file_name, '') message_file_name,
  coalesce(m.message_file_path, '') message_file_path,
  coalesce(us.user_id, 0) user_id,
  coalesce(us.user_name, '') user_name,
  coalesce(us.user_login, '') user_login,
  coalesce(us.user_password, '') user_password,
  coalesce(us.user_mod, false) user_mod
from res m left join users us on
  m.message_user = us.user_id
`

func (persister *persister) EditMessage(msg *model.Message) (err error) {
	row := persister.editMessage.QueryRowx(msg)

	if row.Err() != nil {
		return row.Err()
	}

	return row.StructScan(msg)
}

const loadMessagesLastQuery = `
with cte as (
  select *, true as present from messages m 
    left join users u on
      m.message_user = u.user_id
  where message_channel_id = $1
),
res as (
  select * from (
    table cte order by message_time desc limit $1
  ) sub right join (select count(*) > $2 from cte) c(more) on true
)
select
  coalesce(message_id, 0) message_id,
  coalesce(message_channel_id, 0) message_channel_id,
  coalesce(message_body, '') message_body,
  coalesce(message_time, 'epoch') message_time,
  coalesce(message_edit, 'epoch') message_edit,
  coalesce(message_trip, '') message_trip,
  coalesce(message_origin, '') message_origin,
  coalesce(message_remote, '') message_remote,
  coalesce(message_file_name, '') message_file_name,
  coalesce(message_file_path, '') message_file_path,
  coalesce(user_id, 0) user_id,
  coalesce(user_name, '') user_name,
  coalesce(user_login, '') user_login,
  coalesce(user_password, '') user_password,
  coalesce(user_mod, false) user_mod,
  coalesce(present, false) present,
  more
from res
`

func (persister *persister) LoadMessagesLast(channelID, reqSize uint64) (msgs []*model.Message, more bool, err error) {
	if reqSize > persister.values.HardLimit {
		reqSize = persister.values.HardLimit
	}

	rows, err := persister.loadMessagesLast.Queryx(channelID, reqSize)
	if err != nil {
		return nil, false, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		v := &struct {
			model.Message
			User    uint64 `db:"message_user"`
			More    bool   `db:"more"`
			Present bool   `db:"present"`
		}{
			More: more,
		}

		err = rows.StructScan(v)
		if err != nil {
			return nil, false, err
		}

		if v.Present {
			msgs = append(msgs, &v.Message)
		}

		more = v.More
	}

	return
}

const loadMessagesBeforeQuery = `
with cte as (
  select *, true as present from messages m 
    left join users u on
      m.message_user = u.user_id
  where
    message_channel_id = $1 and 
    message_id < $2
),
res as (
  select * from (
    table cte order by message_time desc limit $3
  ) sub right join (select count(*) > $3 from cte) c(more) on true
)
select
  coalesce(message_id, 0) message_id,
  coalesce(message_channel_id, 0) message_channel_id,
  coalesce(message_body, '') message_body,
  coalesce(message_time, 'epoch') message_time,
  coalesce(message_edit, 'epoch') message_edit,
  coalesce(message_trip, '') message_trip,
  coalesce(message_origin, '') message_origin,
  coalesce(message_remote, '') message_remote,
  coalesce(message_file_name, '') message_file_name,
  coalesce(message_file_path, '') message_file_path,
  coalesce(user_id, 0) user_id,
  coalesce(user_name, '') user_name,
  coalesce(user_login, '') user_login,
  coalesce(user_password, '') user_password,
  coalesce(user_mod, false) user_mod,
  coalesce(present, false) present,
  more
from res
`

func (persister *persister) LoadMessagesBefore(
	channelID, messageID, reqSize uint64,
) (msgs []*model.Message, more bool, err error) {
	if reqSize > persister.values.HardLimit {
		reqSize = persister.values.HardLimit
	}

	rows, err := persister.loadMessagesBefore.Queryx(channelID, messageID, reqSize)
	if err != nil {
		return nil, false, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		v := &struct {
			model.Message
			User    uint64 `db:"message_user"`
			More    bool   `db:"more"`
			Present bool   `db:"present"`
		}{
			More: more,
		}

		err = rows.StructScan(v)
		if err != nil {
			return nil, false, err
		}

		if v.Present {
			msgs = append(msgs, &v.Message)
		}

		more = v.More
	}

	return
}

const loadMessagesPageQuery = `
with cte as (
  select *, true as present from messages m 
    left join users u on
      m.message_user = u.user_id
  where message_channel_id = $1
),
res as (
  select * from (
    table cte order by message_time desc offset $2 limit $3
  ) sub right join (select (count(*) / $3 + case when (count(*) % $3) != 0 then 1 else 0 end) from cte) c(pages) on true
)
select
  coalesce(message_id, 0) message_id,
  coalesce(message_channel_id, 0) message_channel_id,
  coalesce(message_body, '') message_body,
  coalesce(message_time, 'epoch') message_time,
  coalesce(message_edit, 'epoch') message_edit,
  coalesce(message_trip, '') message_trip,
  coalesce(message_origin, '') message_origin,
  coalesce(message_remote, '') message_remote,
  coalesce(message_file_name, '') message_file_name,
  coalesce(message_file_path, '') message_file_path,
  coalesce(user_id, 0) user_id,
  coalesce(user_name, '') user_name,
  coalesce(user_login, '') user_login,
  coalesce(user_password, '') user_password,
  coalesce(user_mod, false) user_mod,
  coalesce(present, false) present,
  pages
from res
`

func (persister *persister) LoadMessagesPage(
	channelID, page, reqSize uint64,
) (msgs []*model.Message, actSize, pages uint64, err error) {
	if reqSize > persister.values.HardLimit {
		reqSize = persister.values.HardLimit
	}

	actSize = reqSize

	rows, err := persister.loadMessagesPage.Queryx(channelID, page*actSize, actSize)
	if err != nil {
		return nil, 0, 0, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		v := &struct {
			model.Message
			User    uint64 `db:"message_user"`
			Pages   uint64 `db:"pages"`
			Present bool   `db:"present"`
		}{
			Pages: pages,
		}

		err = rows.StructScan(v)
		if err != nil {
			return nil, 0, 0, err
		}

		if v.Present {
			msgs = append(msgs, &v.Message)
		}

		pages = v.Pages
	}

	return
}

const deleteUserQuery = `
delete from users where
  user_id = $1
`

func (persister *persister) DeleteUser(userID uint64) (err error) {
	_, err = persister.deleteUser.Exec(userID)
	return
}

const loadUserQuery = `
select
  coalesce(user_id, 0) user_id,
  coalesce(user_name, '') user_name,
  coalesce(user_name, '') user_name,
  coalesce(user_login, '') user_login,
  coalesce(user_password, '') user_password,
  coalesce(user_mod, false) user_mod
from users where
  user_login = $1 and
  user_password = $2
`

func (persister *persister) LoadUser(login, password string) (user *model.User, err error) {
	row := persister.loadUser.QueryRowx(login, password)

	if row.Err() != nil {
		return nil, row.Err()
	}

	user = &model.User{}
	err = row.StructScan(user)

	return
}

const loadUserByIDQuery = `
select
  coalesce(user_id, 0) user_id,
  coalesce(user_name, '') user_name,
  coalesce(user_name, '') user_name,
  coalesce(user_login, '') user_login,
  coalesce(user_password, '') user_password,
  coalesce(user_mod, false) user_mod
from users where
  user_id = $1
`

func (persister *persister) LoadUserByID(id uint64) (user *model.User, err error) {
	row := persister.loadUserByID.QueryRowx(id)

	if row.Err() != nil {
		return nil, row.Err()
	}

	user = &model.User{}
	err = row.StructScan(user)

	return
}

const saveUserQuery = `
insert into users(
  user_id,
  user_name,
  user_login,
  user_password,
  user_mod
) values (
  coalesce(nullif(:user_id, 0), nextval('users_user_id_seq')),
  :user_name,
  :user_login,
  :user_password,
  :user_mod
) on conflict (user_id) do update set
  user_name     = EXCLUDED.user_name,
  user_login    = EXCLUDED.user_login,
  user_password = EXCLUDED.user_password,
  user_mod      = EXCLUDED.user_mod
returning
  coalesce(user_id, 0) user_id,
  coalesce(user_name, '') user_name,
  coalesce(user_login, '') user_login,
  coalesce(user_password, '') user_password,
  coalesce(user_mod, false) user_mod
`

func (persister *persister) SaveUser(user *model.User) (err error) {
	row := persister.saveUser.QueryRowx(user)

	if row.Err() != nil {
		return row.Err()
	}

	return row.StructScan(user)
}

const ensureUserQuery = `
insert into users(
  user_id,
  user_name,
  user_login,
  user_password,
  user_mod
) values (
  default,
  :user_name,
  :user_login,
  :user_password,
  :user_mod
) on conflict (user_login) do nothing
`

func (persister *persister) EnsureUser(user *model.User) (err error) {
	_, err = persister.ensureUser.Exec(user)
	return
}

const loadUsersPageQuery = `
with cte as (
  select *, true as present from users
),
res as (
  select * from (
    table cte order by user_id asc offset $1 limit $2
  ) sub right join (select (count(*) / $2 + case when (count(*) % $2) != 0 then 1 else 0 end) from cte) c(pages) on true
)
select
  coalesce(user_id, 0) user_id,
  coalesce(user_name, '') user_name,
  coalesce(user_login, '') user_login,
  coalesce(user_password, '') user_password,
  coalesce(user_mod, false) user_mod,
  coalesce(present, false) present,
  pages
from res
`

func (persister *persister) LoadUsersPage(
	page, reqSize uint64,
) (users []*model.User, actSize, pages uint64, err error) {
	if reqSize > persister.values.HardLimit {
		reqSize = persister.values.HardLimit
	}

	actSize = reqSize

	rows, err := persister.loadUsersPage.Queryx(page*actSize, actSize)
	if err != nil {
		return nil, 0, 0, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		v := &struct {
			model.User
			Pages   uint64 `db:"pages"`
			Present bool   `db:"present"`
		}{
			Pages: pages,
		}

		err = rows.StructScan(v)
		if err != nil {
			return nil, 0, 0, err
		}

		if v.Present {
			users = append(users, &v.User)
		}

		pages = v.Pages
	}

	return
}

const deleteChannelQuery = `
delete from channels where
  channel_id = $1
`

func (persister *persister) DeleteChannel(channelID uint64) (err error) {
	_, err = persister.deleteChannel.Exec(channelID)
	return
}

const saveChannelQuery = `
insert into channels(
  channel_id,
  channel_name,
  channel_order
) values (
  coalesce(nullif(:channel_id, 0), nextval('channels_channel_id_seq')),
  :channel_name,
  :channel_order
) on conflict (channel_id) do update set
  channel_name    = EXCLUDED.channel_name,
  channel_order   = EXCLUDED.channel_order
returning
  coalesce(channel_id, 0) channel_id,
  coalesce(channel_name, '') channel_name,
  coalesce(channel_order, 0) channel_order
`

func (persister *persister) SaveChannel(channel *model.Channel) (err error) {
	row := persister.saveChannel.QueryRowx(channel)

	if row.Err() != nil {
		return row.Err()
	}

	return row.StructScan(channel)
}

const loadChannelsQuery = `
select
  coalesce(channel_id, 0) channel_id,
  coalesce(channel_name, '') channel_name,
  coalesce(channel_order, 0) channel_order
from channels
`

func (persister *persister) LoadChannels() (channels []*model.Channel, err error) {
	rows, err := persister.loadChannels.Queryx()
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = rows.Close()
	}()

	for rows.Next() {
		v := &model.Channel{}

		err = rows.StructScan(v)
		if err != nil {
			return nil, err
		}

		channels = append(channels, v)
	}

	return
}
