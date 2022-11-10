package leveldb

import (
	"strconv"

	"github.com/Mrs4s/MiraiGo/utils"
	"github.com/pkg/errors"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

type Database struct {
	db *leveldb.DB
}

func Open(dbpath string) (ldb *Database, err error) {
	d, err := leveldb.OpenFile(dbpath, &opt.Options{
		WriteBuffer: 32 * opt.KiB,
	})
	if err != nil {
		return nil, errors.Wrap(err, "open leveldb error")
	}
	return &Database{db: d}, nil
}

func (ldb *Database) Close() error {
	return ldb.db.Close()
}

func (ldb *Database) ForEach(f func(x any) error) (errs []error) {
	iter := ldb.db.NewIterator(nil, nil)
	for iter.Next() {
		value := iter.Value()
		if len(value) == 0 {
			continue
		}
		r, err := newReader(utils.B2S(value))
		if err != nil {
			errs = append(errs, errors.Wrap(err, "new reader failed"))
			continue
		}
		flg := r.uvarint()
		switch flg {
		case group:
			err = f(r.readStoredGroupMessage())
			if err != nil {
				errs = append(errs, errors.Wrap(err, "decode group message callback failed"))
			}
		case private:
			err = f(r.readStoredPrivateMessage())
			if err != nil {
				errs = append(errs, errors.Wrap(err, "decode private message callback failed"))
			}
		case guildChannel:
			err = f(r.readStoredGuildChannelMessage())
			if err != nil {
				errs = append(errs, errors.Wrap(err, "decode guild channel message callback failed"))
			}
		default:
			errs = append(errs, errors.New("unknown message flag "+strconv.Itoa(int(flg))))
		}

	}
	iter.Release()
	return
}
