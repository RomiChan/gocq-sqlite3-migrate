//go:build !(windows && (arm || arm64))
// +build !windows !arm,!arm64

package sqlite3

import (
	"encoding/json"
	"hash/crc64"
	"sync"
	"time"

	sql "github.com/FloatTech/sqlite"
	"github.com/pkg/errors"

	"github.com/Mrs4s/MiraiGo/binary"
	"github.com/Mrs4s/MiraiGo/utils"
	"github.com/Mrs4s/go-cqhttp/db"
)

type Database struct {
	sync.RWMutex
	db  *sql.Sqlite
	ttl time.Duration
}

func Open(dbpath string) (s *Database, err error) {
	s = &Database{db: new(sql.Sqlite)}
	s.db.DBPath = dbpath
	err = s.db.Open(s.ttl)
	if err != nil {
		return nil, errors.Wrap(err, "open sqlite3 error")
	}
	err = s.db.Create(Sqlite3GroupMessageTableName, &StoredGroupMessage{})
	if err != nil {
		return nil, errors.Wrap(err, "create sqlite3 table error")
	}
	err = s.db.Create(Sqlite3MessageAttributeTableName, &StoredMessageAttribute{})
	if err != nil {
		return nil, errors.Wrap(err, "create sqlite3 table error")
	}
	err = s.db.Create(Sqlite3GuildMessageAttributeTableName, &StoredGuildMessageAttribute{})
	if err != nil {
		return nil, errors.Wrap(err, "create sqlite3 table error")
	}
	err = s.db.Create(Sqlite3QuotedInfoTableName, &QuotedInfo{})
	if err != nil {
		return nil, errors.Wrap(err, "create sqlite3 table error")
	}
	err = s.db.Create(Sqlite3PrivateMessageTableName, &StoredPrivateMessage{})
	if err != nil {
		return nil, errors.Wrap(err, "create sqlite3 table error")
	}
	err = s.db.Create(Sqlite3GuildChannelMessageTableName, &StoredGuildChannelMessage{})
	if err != nil {
		return nil, errors.Wrap(err, "create sqlite3 table error")
	}
	err = s.db.Create(Sqlite3UinInfoTableName, &UinInfo{})
	if err != nil {
		return nil, errors.Wrap(err, "create sqlite3 table error")
	}
	err = s.db.Create(Sqlite3TinyInfoTableName, &TinyInfo{})
	if err != nil {
		return nil, errors.Wrap(err, "create sqlite3 table error")
	}
	return s, nil
}

func (s *Database) Close() error {
	return s.db.Close()
}

func (s *Database) InsertGroupMessage(msg *db.StoredGroupMessage) error {
	grpmsg := &StoredGroupMessage{
		GlobalID:    msg.GlobalID,
		ID:          msg.ID,
		SubType:     msg.SubType,
		GroupCode:   msg.GroupCode,
		AnonymousID: msg.AnonymousID,
	}
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	if msg.Attribute != nil {
		h.Write(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt32(uint32(msg.Attribute.MessageSeq))
			w.WriteUInt32(uint32(msg.Attribute.InternalID))
			w.WriteUInt64(uint64(msg.Attribute.SenderUin))
			w.WriteUInt64(uint64(msg.Attribute.Timestamp))
		}))
		h.Write(utils.S2B(msg.Attribute.SenderName))
		id := int64(h.Sum64())
		if id == 0 {
			id++
		}
		s.Lock()
		err := s.db.Insert(Sqlite3UinInfoTableName, &UinInfo{
			Uin:  msg.Attribute.SenderUin,
			Name: msg.Attribute.SenderName,
		})
		if err == nil {
			err = s.db.Insert(Sqlite3MessageAttributeTableName, &StoredMessageAttribute{
				ID:         id,
				MessageSeq: msg.Attribute.MessageSeq,
				InternalID: msg.Attribute.InternalID,
				SenderUin:  msg.Attribute.SenderUin,
				Timestamp:  msg.Attribute.Timestamp,
			})
		}
		s.Unlock()
		if err == nil {
			grpmsg.AttributeID = id
		}
		h.Reset()
	}
	if msg.QuotedInfo != nil {
		h.Write(utils.S2B(msg.QuotedInfo.PrevID))
		h.Write(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt32(uint32(msg.QuotedInfo.PrevGlobalID))
		}))
		content, err := json.Marshal(&msg.QuotedInfo.QuotedContent)
		if err != nil {
			return errors.Wrap(err, "insert marshal QuotedContent error")
		}
		h.Write(content)
		id := int64(h.Sum64())
		if id == 0 {
			id++
		}
		s.Lock()
		err = s.db.Insert(Sqlite3QuotedInfoTableName, &QuotedInfo{
			ID:            id,
			PrevID:        msg.QuotedInfo.PrevID,
			PrevGlobalID:  msg.QuotedInfo.PrevGlobalID,
			QuotedContent: utils.B2S(content),
		})
		s.Unlock()
		if err == nil {
			grpmsg.QuotedInfoID = id
		}
	}
	content, err := json.Marshal(&msg.Content)
	if err != nil {
		return errors.Wrap(err, "insert marshal Content error")
	}
	grpmsg.Content = utils.B2S(content)
	s.Lock()
	err = s.db.Insert(Sqlite3GroupMessageTableName, grpmsg)
	s.Unlock()
	if err != nil {
		return errors.Wrap(err, "insert error")
	}
	return nil
}

func (s *Database) InsertPrivateMessage(msg *db.StoredPrivateMessage) error {
	privmsg := &StoredPrivateMessage{
		GlobalID:   msg.GlobalID,
		ID:         msg.ID,
		SubType:    msg.SubType,
		SessionUin: msg.SessionUin,
		TargetUin:  msg.TargetUin,
	}
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	if msg.Attribute != nil {
		h.Write(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt32(uint32(msg.Attribute.MessageSeq))
			w.WriteUInt32(uint32(msg.Attribute.InternalID))
			w.WriteUInt64(uint64(msg.Attribute.SenderUin))
			w.WriteUInt64(uint64(msg.Attribute.Timestamp))
		}))
		h.Write(utils.S2B(msg.Attribute.SenderName))
		id := int64(h.Sum64())
		if id == 0 {
			id++
		}
		s.Lock()
		err := s.db.Insert(Sqlite3UinInfoTableName, &UinInfo{
			Uin:  msg.Attribute.SenderUin,
			Name: msg.Attribute.SenderName,
		})
		if err == nil {
			err = s.db.Insert(Sqlite3MessageAttributeTableName, &StoredMessageAttribute{
				ID:         id,
				MessageSeq: msg.Attribute.MessageSeq,
				InternalID: msg.Attribute.InternalID,
				SenderUin:  msg.Attribute.SenderUin,
				Timestamp:  msg.Attribute.Timestamp,
			})
		}
		s.Unlock()
		if err == nil {
			privmsg.AttributeID = id
		}
		h.Reset()
	}
	if msg.QuotedInfo != nil {
		h.Write(utils.S2B(msg.QuotedInfo.PrevID))
		h.Write(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt32(uint32(msg.QuotedInfo.PrevGlobalID))
		}))
		content, err := json.Marshal(&msg.QuotedInfo.QuotedContent)
		if err != nil {
			return errors.Wrap(err, "insert marshal QuotedContent error")
		}
		h.Write(content)
		id := int64(h.Sum64())
		if id == 0 {
			id++
		}
		s.Lock()
		err = s.db.Insert(Sqlite3QuotedInfoTableName, &QuotedInfo{
			ID:            id,
			PrevID:        msg.QuotedInfo.PrevID,
			PrevGlobalID:  msg.QuotedInfo.PrevGlobalID,
			QuotedContent: utils.B2S(content),
		})
		s.Unlock()
		if err == nil {
			privmsg.QuotedInfoID = id
		}
	}
	content, err := json.Marshal(&msg.Content)
	if err != nil {
		return errors.Wrap(err, "insert marshal Content error")
	}
	privmsg.Content = utils.B2S(content)
	s.Lock()
	err = s.db.Insert(Sqlite3PrivateMessageTableName, privmsg)
	s.Unlock()
	if err != nil {
		return errors.Wrap(err, "insert error")
	}
	return nil
}

func (s *Database) InsertGuildChannelMessage(msg *db.StoredGuildChannelMessage) error {
	guildmsg := &StoredGuildChannelMessage{
		ID:        msg.ID,
		GuildID:   int64(msg.GuildID),
		ChannelID: int64(msg.ChannelID),
	}
	h := crc64.New(crc64.MakeTable(crc64.ISO))
	if msg.Attribute != nil {
		h.Write(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt32(uint32(msg.Attribute.MessageSeq))
			w.WriteUInt32(uint32(msg.Attribute.InternalID))
			w.WriteUInt64(uint64(msg.Attribute.SenderTinyID))
			w.WriteUInt64(uint64(msg.Attribute.Timestamp))
		}))
		h.Write(utils.S2B(msg.Attribute.SenderName))
		id := int64(h.Sum64())
		if id == 0 {
			id++
		}
		s.Lock()
		err := s.db.Insert(Sqlite3TinyInfoTableName, &TinyInfo{
			ID:   int64(msg.Attribute.SenderTinyID),
			Name: msg.Attribute.SenderName,
		})
		if err == nil {
			err = s.db.Insert(Sqlite3MessageAttributeTableName, &StoredGuildMessageAttribute{
				ID:           id,
				MessageSeq:   int64(msg.Attribute.MessageSeq),
				InternalID:   int64(msg.Attribute.InternalID),
				SenderTinyID: int64(msg.Attribute.SenderTinyID),
				Timestamp:    msg.Attribute.Timestamp,
			})
		}
		s.Unlock()
		if err == nil {
			guildmsg.AttributeID = id
		}
		h.Reset()
	}
	if msg.QuotedInfo != nil {
		h.Write(utils.S2B(msg.QuotedInfo.PrevID))
		h.Write(binary.NewWriterF(func(w *binary.Writer) {
			w.WriteUInt32(uint32(msg.QuotedInfo.PrevGlobalID))
		}))
		content, err := json.Marshal(&msg.QuotedInfo.QuotedContent)
		if err != nil {
			return errors.Wrap(err, "insert marshal QuotedContent error")
		}
		h.Write(content)
		id := int64(h.Sum64())
		if id == 0 {
			id++
		}
		s.Lock()
		err = s.db.Insert(Sqlite3QuotedInfoTableName, &QuotedInfo{
			ID:            id,
			PrevID:        msg.QuotedInfo.PrevID,
			PrevGlobalID:  msg.QuotedInfo.PrevGlobalID,
			QuotedContent: utils.B2S(content),
		})
		s.Unlock()
		if err == nil {
			guildmsg.QuotedInfoID = id
		}
	}
	content, err := json.Marshal(&msg.Content)
	if err != nil {
		return errors.Wrap(err, "insert marshal Content error")
	}
	guildmsg.Content = utils.B2S(content)
	s.Lock()
	err = s.db.Insert(Sqlite3GuildChannelMessageTableName, guildmsg)
	s.Unlock()
	if err != nil {
		return errors.Wrap(err, "insert error")
	}
	return nil
}
