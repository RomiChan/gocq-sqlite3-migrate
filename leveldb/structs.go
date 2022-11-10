package leveldb

import "github.com/Mrs4s/go-cqhttp/db"

func (r *reader) readStoredGroupMessage() *db.StoredGroupMessage {
	coder := r.coder()
	if coder == coderNil {
		return nil
	}
	x := &db.StoredGroupMessage{}
	x.ID = r.string()
	x.GlobalID = r.int32()
	x.Attribute = r.readStoredMessageAttribute()
	x.SubType = r.string()
	x.QuotedInfo = r.readQuotedInfo()
	x.GroupCode = r.int64()
	x.AnonymousID = r.string()
	x.Content = r.arrayMsg()
	return x
}

func (r *reader) readStoredPrivateMessage() *db.StoredPrivateMessage {
	coder := r.coder()
	if coder == coderNil {
		return nil
	}
	x := &db.StoredPrivateMessage{}
	x.ID = r.string()
	x.GlobalID = r.int32()
	x.Attribute = r.readStoredMessageAttribute()
	x.SubType = r.string()
	x.QuotedInfo = r.readQuotedInfo()
	x.SessionUin = r.int64()
	x.TargetUin = r.int64()
	x.Content = r.arrayMsg()
	return x
}

func (r *reader) readStoredGuildChannelMessage() *db.StoredGuildChannelMessage {
	coder := r.coder()
	if coder == coderNil {
		return nil
	}
	x := &db.StoredGuildChannelMessage{}
	x.ID = r.string()
	x.Attribute = r.readStoredGuildMessageAttribute()
	x.GuildID = r.uint64()
	x.ChannelID = r.uint64()
	x.QuotedInfo = r.readQuotedInfo()
	x.Content = r.arrayMsg()
	return x
}

func (r *reader) readStoredMessageAttribute() *db.StoredMessageAttribute {
	coder := r.coder()
	if coder == coderNil {
		return nil
	}
	x := &db.StoredMessageAttribute{}
	x.MessageSeq = r.int32()
	x.InternalID = r.int32()
	x.SenderUin = r.int64()
	x.SenderName = r.string()
	x.Timestamp = r.int64()
	return x
}

func (r *reader) readStoredGuildMessageAttribute() *db.StoredGuildMessageAttribute {
	coder := r.coder()
	if coder == coderNil {
		return nil
	}
	x := &db.StoredGuildMessageAttribute{}
	x.MessageSeq = r.uint64()
	x.InternalID = r.uint64()
	x.SenderTinyID = r.uint64()
	x.SenderName = r.string()
	x.Timestamp = r.int64()
	return x
}

func (r *reader) readQuotedInfo() *db.QuotedInfo {
	coder := r.coder()
	if coder == coderNil {
		return nil
	}
	x := &db.QuotedInfo{}
	x.PrevID = r.string()
	x.PrevGlobalID = r.int32()
	x.QuotedContent = r.arrayMsg()
	return x
}
